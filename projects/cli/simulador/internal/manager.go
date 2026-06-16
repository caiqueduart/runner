package internal

import (
	"crypto/tls"
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

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

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

	LogFeedback("SIMULADOR CONFIG", "JDK 21 não encontrado. Iniciando download...")
	if err := DownloadJava21(GetJDKDir()); err != nil {
		return "", fmt.Errorf("falha ao baixar JDK 21: %w", err)
	}

	LogFeedback("SIMULADOR CONFIG", "JDK 21 instalado com sucesso.")
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
	LogFeedback("SIMULADOR CONFIG", "Baixando JDK...")

	os.MkdirAll(targetDir, 0755)
	tmpFile := filepath.Join(os.TempDir(), "jdk21_download"+filepath.Ext(downloadUrl))
	out, _ := os.Create(tmpFile)
	defer os.Remove(tmpFile)

	resp, _ = http.Get(downloadUrl)
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	out.Close()

	LogFeedback("SIMULADOR CONFIG", "Extraindo arquivos...")
	if strings.HasSuffix(downloadUrl, ".zip") {
		return extractZip(tmpFile, targetDir)
	}
	return extractTarGz(tmpFile, targetDir)
}

func DownloadSimuladorJar(targetPath string) error {
	LogFeedback("SIMULADOR CONFIG", "JAR do simulador não encontrado. Baixando versão %s...", CompatibleSimuladorVersion)
	
	// Usando link direto para a versão solicitada
	downloadUrl := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", RepoPath, SimuladorTag, SimuladorJarName)

	os.MkdirAll(filepath.Dir(targetPath), 0755)
	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo JAR: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(downloadUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao baixar JAR da URL: %s (Status: %v)", downloadUrl, resp.StatusCode)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("falha ao salvar conteúdo do JAR: %w", err)
	}

	LogFeedback("SIMULADOR CONFIG", "Simulador baixado com sucesso.")
	return nil
}

func EnsureSimuladorRunning() error {
	// O simulador é um servidor HTTPS
	url := fmt.Sprintf("https://localhost:%s/api/info", SimuladorPort)
	resp, err := httpClient.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return nil
	}

	LogFeedback("SIMULADOR CONFIG", "Simulador não detectado na porta %s. Iniciando...", SimuladorPort)

	javaPath, err := GetJavaPath("java")
	if err != nil {
		return err
	}

	localJarPath := GetSimuladorJarPath()
	if _, err := os.Stat(localJarPath); os.IsNotExist(err) {
		if err := DownloadSimuladorJar(localJarPath); err != nil {
			return err
		}
	}

	cmd := exec.Command(javaPath, "-jar", localJarPath)
	SetDetachedProcess(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("falha ao iniciar simulador: %w", err)
	}

	if cmd.Process != nil {
		savePID(cmd.Process.Pid, SimuladorPort)
		cmd.Process.Release()
	}

	// Aguarda o simulador subir
	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		resp, err := httpClient.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			LogFeedback("SIMULADOR SERVIDOR", "Simulador online.")
			return nil
		}
	}
	return fmt.Errorf("timeout ao aguardar o simulador subir")
}

func savePID(pid int, port string) {
	path := GetPIDFilePath()
	os.MkdirAll(filepath.Dir(path), 0755)
	content := fmt.Sprintf("PID=%d\nPORT=%s\n", pid, port)
	os.WriteFile(path, []byte(content), 0644)
	LogFeedback("SIMULADOR CONFIG", "Rastreabilidade registrada (PID: %d, Porta: %s).", pid, port)
}

func ClearPIDFile() {
	os.Remove(GetPIDFilePath())
}

func StopSimulador() error {
	url := fmt.Sprintf("https://localhost:%s/shutdown", SimuladorPort)
	resp, err := httpClient.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("falha ao enviar comando de parada: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		LogFeedback("SIMULADOR SERVIDOR", "Comando de parada enviado com sucesso.")
		ClearPIDFile()
		return nil
	}
	
	return fmt.Errorf("erro ao parar simulador: status %d", resp.StatusCode)
}

func GetSimuladorStatus() (string, error) {
	url := fmt.Sprintf("https://localhost:%s/api/info", SimuladorPort)
	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}
