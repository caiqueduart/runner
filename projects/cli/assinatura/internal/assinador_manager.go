package internal

import (
	"bytes"
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

// localiza o executável do Java 21, priorizando o provisionamento automático
// da pasta .hubsaude/jdk caso não esteja no PATH.
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

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório do JDK: %w", err)
	}
	tmpFile := filepath.Join(os.TempDir(), "jdk21_download"+filepath.Ext(downloadUrl))
	out, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo temporário do JDK: %w", err)
	}

	defer os.Remove(tmpFile)

	resp, err = http.Get(downloadUrl)
	if err != nil {
		out.Close()
		return fmt.Errorf("falha ao baixar JDK: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		out.Close()
		return fmt.Errorf("falha ao baixar JDK: status HTTP %d", resp.StatusCode)
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return fmt.Errorf("falha ao salvar JDK: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("falha ao finalizar arquivo do JDK: %w", err)
	}

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
	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo local: %w", err)
	}
	defer out.Close()

	resp, err = http.Get(jarURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao baixar JAR da release: %v (Status: %v)", err, resp.StatusCode)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("falha ao salvar conteúdo do JAR: %w", err)
	}

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

type SignatureResponse struct {
	FileName       string `json:"fileName"`
	Code           string `json:"code"`
	SignOutputPath string `json:"signOutputPath"`
	Message        string `json:"message"`
	Status         int    `json:"status"`
}

type ValidationResponse struct {
	FileName string `json:"fileName"`
	Code     string `json:"code"`
	Valid    bool   `json:"valid"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Type   string `json:"type"`
}

// custom erro para transportar o tipo de erro (user ou system)
type JavaError struct {
	Msg  string
	Type string
}

func (e *JavaError) Error() string {
	return e.Msg
}

type ExecError struct {
	Output string
	Code   int
}

func (e *ExecError) Error() string {
	return e.Output
}

func CallJavaServer(endpoint string, cmdKey string, fileName string) (string, error) {
	url := fmt.Sprintf("http://localhost:%s/%s", ServerPort, endpoint)

	requestData := map[string]string{
		"command": cmdKey,
		"file":    fileName,
		"flag":    "--file",
	}
	jsonData, _ := json.Marshal(requestData)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			return "", &JavaError{Msg: errResp.Error, Type: errResp.Type}
		}
		return "", fmt.Errorf("erro %d: %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// verifica se o servidor HTTP do assinador está ativo.
// caso contrário, inicia-o em background e aguarda o status online.
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

	var cmd *exec.Cmd
	if os.Getenv("DEV_MODE") == "true" {
		assinadorDir := getAssinadorDir()

		if assinadorDir == "" {
			return fmt.Errorf("não foi possível localizar projects/assinador")
		}
		srcDir := filepath.Join(assinadorDir, "src", "main", "java")
		appJava := filepath.Join(srcDir, "App.java")
		cmd = exec.Command(javaPath, "-cp", srcDir, appJava, "server", "--port", ServerPort, "--timeout", ServerTimeoutMinutes)

	} else {
		localJarPath := GetJarPath()

		if _, err := os.Stat(localJarPath); os.IsNotExist(err) {
			if err := DownloadAssinadorJar(localJarPath); err != nil {
				return err
			}
		}

		cmd = exec.Command(javaPath, "-jar", localJarPath, "server", "--port", ServerPort, "--timeout", ServerTimeoutMinutes)
	}

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

func getAssinadorDir() string {
	currDir, _ := os.Getwd()
	tempDir := currDir

	for i := 0; i < 6; i++ {
		target := filepath.Join(tempDir, "projects", "assinador")
		if _, err := os.Stat(target); err == nil {
			return target
		}
		parent := filepath.Dir(tempDir)
		if parent == tempDir {
			break
		}
		tempDir = parent
	}
	return ""
}

// orquestra a execução da assinatura/validação.
// tenta usar o modo servidor (HTTP) e regride para modo local (JAR direto) se necessário.
func ExecJavaSigner(cmdKey string, cmdArgs []string) (string, error) {
	// MODO DESENVOLVEDOR: Executa direto do .java se a variável DEV_MODE=true no arquivo .env
	if os.Getenv("DEV_MODE") == "true" {
		LogFeedback("ASSINATURA CONFIG", "Modo Desenvolvedor Ativo (Lendo App.java).")

		assinadorDir := getAssinadorDir()
		if assinadorDir == "" {
			return "", fmt.Errorf("não foi possível localizar projects/assinador")
		}

		srcDir := filepath.Join(assinadorDir, "src", "main", "java")
		appJava := filepath.Join(srcDir, "App.java")

		fileName := ""
		for i, arg := range cmdArgs {
			if arg == "--file" && i+1 < len(cmdArgs) {
				fileName = cmdArgs[i+1]
			} else if !strings.HasPrefix(arg, "-") && cmdKey == "validate" {
				fileName = arg
			}
		}

		requestData := map[string]string{
			"command": cmdKey,
			"file":    fileName,
			"flag":    "--file",
		}
		jsonData, _ := json.Marshal(requestData)

		args := []string{"-cp", srcDir, appJava, string(jsonData)}
		javaCmd := exec.Command("java", args...)

		output, err := javaCmd.CombinedOutput()

		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				var errResp ErrorResponse
				if parseErr := json.Unmarshal(output, &errResp); parseErr == nil && errResp.Error != "" {
					return "", &JavaError{Msg: errResp.Error, Type: errResp.Type}
				}
				return "", &ExecError{Output: string(output), Code: exitError.ExitCode()}
			}
			return "", err
		}

		return string(output), nil
	}

	if err := EnsureServerRunning(); err == nil {
		fileName := ""

		for i, arg := range cmdArgs {
			if arg == "--file" && i+1 < len(cmdArgs) {
				fileName = cmdArgs[i+1]
			} else if !strings.HasPrefix(arg, "-") && cmdKey == "validate" {
				fileName = arg
			}
		}

		return CallJavaServer(cmdKey, cmdKey, fileName)
	}

	LogFeedback("ASSINATURA CONFIG", "Servidor indisponível. Usando modo local...")

	javaPath, _ := GetJavaPath("java")
	localJarPath := GetJarPath()

	fileName := ""
	for i, arg := range cmdArgs {
		if arg == "--file" && i+1 < len(cmdArgs) {
			fileName = cmdArgs[i+1]
		} else if !strings.HasPrefix(arg, "-") && cmdKey == "validate" {
			fileName = arg
		}
	}

	requestData := map[string]string{
		"command": cmdKey,
		"file":    fileName,
		"flag":    "--file",
	}

	jsonData, _ := json.Marshal(requestData)

	args := []string{"-jar", localJarPath, string(jsonData)}
	javaCmd := exec.Command(javaPath, args...)

	output, err := javaCmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			var errResp ErrorResponse
			if parseErr := json.Unmarshal(output, &errResp); parseErr == nil && errResp.Error != "" {
				return "", &JavaError{Msg: errResp.Error, Type: errResp.Type}
			}
			return "", &ExecError{Output: string(output), Code: exitError.ExitCode()}
		}
		return "", err
	}

	return string(output), nil
}
