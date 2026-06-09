package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetJavaPath(binName string) (string, error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(binName, ".exe") {
		binName += ".exe"
	}

	managedBin := filepath.Join(GetJDKDir(), "bin", binName)

	if path, err := exec.LookPath(binName); err == nil {
		if isJava21(path) {
			return path, nil
		}
	}

	if _, err := os.Stat(managedBin); err == nil {
		if isJava21(managedBin) {
			return managedBin, nil
		}
	}

	LogFeedback("ASSINATURA CONFIG", "JDK 21 não encontrado. Iniciando download...")
	if err := DownloadJava21(GetJDKDir()); err != nil {
		return "", fmt.Errorf("falha ao baixar JDK 21: %w", err)
	}

	LogFeedback("ASSINATURA CONFIG", "JDK 21 instalado com sucesso.")
	return managedBin, nil
}

func isJava21(javaPath string) bool {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "version \"21")
}

func DownloadJava21(targetDir string) error {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "mac"
	}

	apiUrl := fmt.Sprintf("https://api.adoptium.net/v3/assets/feature_releases/21/ga?architecture=%s&image_type=jdk&os=%s&project=jdk&vendor=eclipse", arch, osName)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var releases []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&releases)
	if len(releases) == 0 {
		return fmt.Errorf("nenhum release encontrado na Adoptium")
	}

	downloadUrl := releases[0]["binaries"].([]interface{})[0].(map[string]interface{})["package"].(map[string]interface{})["link"].(string)
	LogFeedback("ASSINATURA CONFIG", "Baixando JDK...")

	os.MkdirAll(targetDir, 0755)
	tmpFile := filepath.Join(os.TempDir(), "jdk21_download"+filepath.Ext(downloadUrl))
	out, _ := os.Create(tmpFile)
	defer os.Remove(tmpFile)

	resp, _ = http.Get(downloadUrl)
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	out.Close()

	LogFeedback("ASSINATURA CONFIG", "Extraindo arquivos...")
	if strings.HasSuffix(downloadUrl, ".zip") {
		return extractZip(tmpFile, targetDir)
	}
	return extractTarGz(tmpFile, targetDir)
}

func DownloadAssinadorJar(targetPath string) error {
	LogFeedback("ASSINATURA CONFIG", "JAR não encontrado. Baixando...")
	tag := "assinador-v" + CompatibleAssinadorVersion
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", RepoPath, tag)

	resp, err := http.Get(apiUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao consultar release no GitHub")
	}
	defer resp.Body.Close()

	var release struct {
		Assets []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
			Digest      string `json:"digest"`
		} `json:"assets"`
	}
	json.NewDecoder(resp.Body).Decode(&release)

	var jarURL, digest string
	expectedName := fmt.Sprintf("assinador-v%s.jar", CompatibleAssinadorVersion)
	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			jarURL = asset.DownloadURL
			digest = asset.Digest
			break
		}
	}

	if jarURL == "" {
		return fmt.Errorf("arquivo %s não encontrado", expectedName)
	}

	os.MkdirAll(filepath.Dir(targetPath), 0755)
	out, _ := os.Create(targetPath)
	resp, _ = http.Get(jarURL)
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	out.Close()

	if digest != "" {
		LogFeedback("ASSINATURA CONFIG", "Validando integridade...")
		if ok, _ := checkFileSHA256(targetPath, digest); !ok {
			os.Remove(targetPath)
			return fmt.Errorf("ERRO DE SEGURANÇA: SHA256 não coincide")
		}
		LogFeedback("ASSINATURA CONFIG", "Integridade OK.")
	}
	return nil
}

func savePID(pid int, port string) {
	path := GetPIDFilePath()
	os.MkdirAll(filepath.Dir(path), 0755)
	content := fmt.Sprintf("PID=%d\nPORT=%s\n", pid, port)
	os.WriteFile(path, []byte(content), 0644)
	LogFeedback("ASSINATURA CONFIG", "Rastreabilidade registrada (PID: %d, Porta: %s).", pid, port)
}

func ClearPIDFile() {
	os.Remove(GetPIDFilePath())
}

func CallJavaServer(endpoint string, data string) (string, error) {
	url := fmt.Sprintf("http://localhost:%s/%s", ServerPort, endpoint)
	resp, err := http.Post(url, "text/plain", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("erro %d: %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

func EnsureServerRunning() error {
	resp, err := http.Get("http://localhost:" + ServerPort + "/health")
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return nil
	}

	LogFeedback("ASSINATURA CONFIG", "Servidor não detectado. Iniciando...")

	javaPath, err := GetJavaPath("java")
	if err != nil {
		return err
	}

	localJarPath := GetJarPath()
	if _, err := os.Stat(localJarPath); os.IsNotExist(err) {
		if err := DownloadAssinadorJar(localJarPath); err != nil {
			return err
		}
	}

	cmd := exec.Command(javaPath, "-jar", localJarPath, "server", "--port", ServerPort, "--timeout", ServerTimeoutMinutes)
	SetDetachedProcess(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("falha ao iniciar servidor: %w", err)
	}

	if cmd.Process != nil {
		savePID(cmd.Process.Pid, ServerPort)
		cmd.Process.Release()
	}

	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		resp, err := http.Get("http://localhost:" + ServerPort + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			LogFeedback("ASSINATURA SERVIDOR", "Servidor online.")
			return nil
		}
	}
	return fmt.Errorf("timeout ao aguardar o servidor subir")
}

func ExecJavaSigner(fileName string, cmdKey string) (string, error) {
	// MODO DESENVOLVEDOR: Executa direto do .java se a variável DEV_MODE=true
	if os.Getenv("DEV_MODE") == "true" {
		LogFeedback("ASSINATURA CONFIG", "Modo Desenvolvedor Ativo (Lendo Main.java).")

		// Descobre o caminho absoluto para a raiz do projeto e depois para o Main.java
		wd, _ := os.Getwd()
		root := filepath.Join(wd, "..", "..", "..")
		javaSource := filepath.Join(root, "projects", "assinador", "src", "main", "Main.java")

		var javaCmd *exec.Cmd
		if cmdKey == "sign" {
			javaCmd = exec.Command("java", javaSource, cmdKey, "--file", fileName)
		} else {
			javaCmd = exec.Command("java", javaSource, cmdKey, fileName)
		}

		output, err := javaCmd.CombinedOutput()
		if err != nil {
			return string(output), nil
		}
		return string(output), nil
	}

	if err := EnsureServerRunning(); err == nil {
		return CallJavaServer(cmdKey, fileName)
	}

	LogFeedback("ASSINATURA CONFIG", "Servidor indisponível. Usando modo local...")
	javaPath, _ := GetJavaPath("java")
	localJarPath := GetJarPath()

	var javaCmd *exec.Cmd
	if cmdKey == "sign" {
		// No sign, obrigamos a flag --file no JAR
		javaCmd = exec.Command(javaPath, "-jar", localJarPath, cmdKey, "--file", fileName)
	} else {
		// No validate, passamos como argumento posicional
		javaCmd = exec.Command(javaPath, "-jar", localJarPath, cmdKey, fileName)
	}

	output, err := javaCmd.CombinedOutput()

	if err != nil {
		return string(output), nil
	}
	return string(output), nil
}
