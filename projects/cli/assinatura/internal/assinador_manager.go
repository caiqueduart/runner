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

	versionOutput := strings.TrimSpace(string(output))
	return strings.Contains(versionOutput, "version \"21") || strings.HasPrefix(versionOutput, "javac 21")
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

type ExecutionOptions struct {
	Local          bool
	Port           int
	TimeoutMinutes int
}

var signerHTTPClient = &http.Client{Timeout: 10 * time.Second}

func (options ExecutionOptions) validate() error {
	if options.Port < 1 || options.Port > 65535 {
		return fmt.Errorf("porta inválida: %d; use um valor entre 1 e 65535", options.Port)
	}
	if options.TimeoutMinutes < 1 {
		return fmt.Errorf("timeout inválido: %d; use ao menos 1 minuto", options.TimeoutMinutes)
	}
	return nil
}

func serverURL(port int, endpoint string) string {
	return fmt.Sprintf("http://localhost:%d/%s", port, endpoint)
}

func CallJavaServer(endpoint string, cmdKey string, fileName string, options ExecutionOptions) (string, error) {
	if err := options.validate(); err != nil {
		return "", err
	}

	url := serverURL(options.Port, endpoint)

	requestData := map[string]string{
		"command": cmdKey,
		"file":    fileName,
		"flag":    "--file",
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("falha ao criar requisição JSON: %w", err)
	}

	resp, err := signerHTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("falha ao chamar assinador em %s: %w", url, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("falha ao ler resposta do assinador: %w", err)
	}

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
	return EnsureServerRunningWithOptions(ExecutionOptions{
		Port:           DefaultServerPort,
		TimeoutMinutes: DefaultServerTimeoutMinutes,
	})
}

func EnsureServerRunningWithOptions(options ExecutionOptions) error {
	if err := options.validate(); err != nil {
		return err
	}

	healthURL := serverURL(options.Port, "health")
	resp, err := signerHTTPClient.Get(healthURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		return nil
	}
	if resp != nil {
		resp.Body.Close()
	}

	LogFeedback("ASSINATURA CONFIG", "Servidor não detectado. Iniciando...")

	javaPath, err := GetJavaPath("java")
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	port := fmt.Sprintf("%d", options.Port)
	timeout := fmt.Sprintf("%d", options.TimeoutMinutes)
	if os.Getenv("DEV_MODE") == "true" {
		classesDir, err := compileDevAssinador()
		if err != nil {
			return err
		}
		cmd = exec.Command(javaPath, "-cp", classesDir, "App", "server", "--port", port, "--timeout", timeout)

	} else {
		localJarPath := GetJarPath()

		if _, err := os.Stat(localJarPath); os.IsNotExist(err) {
			if err := DownloadAssinadorJar(localJarPath); err != nil {
				return err
			}
		}

		cmd = exec.Command(javaPath, "-jar", localJarPath, "server", "--port", port, "--timeout", timeout)
	}

	logPath := filepath.Join(GetHubSaudeDir(), "assinador-server.log")
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("falha ao preparar logs do servidor: %w", err)
	}
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("falha ao abrir log do servidor: %w", err)
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	SetDetachedProcess(cmd)

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("falha ao iniciar servidor: %w", err)
	}
	logFile.Close()

	if cmd.Process != nil {
		savePID(cmd.Process.Pid, port)
		cmd.Process.Release()
	}

	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		resp, err := signerHTTPClient.Get(healthURL)

		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			LogFeedback("ASSINATURA SERVIDOR", "Servidor online.")

			return nil
		}
	}

	diagnostic, _ := os.ReadFile(logPath)
	if message := strings.TrimSpace(string(diagnostic)); message != "" {
		return fmt.Errorf("timeout ao aguardar o servidor subir; diagnóstico: %s", message)
	}
	return fmt.Errorf("timeout ao aguardar o servidor subir; consulte %s", logPath)
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

func compileDevAssinador() (string, error) {
	assinadorDir := getAssinadorDir()
	if assinadorDir == "" {
		return "", fmt.Errorf("não foi possível localizar projects/assinador")
	}

	srcDir := filepath.Join(assinadorDir, "src", "main", "java")
	serviceSources, err := filepath.Glob(filepath.Join(srcDir, "services", "*.java"))
	if err != nil {
		return "", fmt.Errorf("falha ao localizar fontes Java: %w", err)
	}
	sources := append([]string{
		filepath.Join(srcDir, "App.java"),
		filepath.Join(srcDir, "Version.java"),
	}, serviceSources...)

	classesDir := filepath.Join(assinadorDir, "target", "dev-classes")
	if err := os.MkdirAll(classesDir, 0755); err != nil {
		return "", fmt.Errorf("falha ao criar diretório de classes: %w", err)
	}

	javacPath, err := GetJavaPath("javac")
	if err != nil {
		return "", err
	}
	args := append([]string{"-d", classesDir}, sources...)
	output, err := exec.Command(javacPath, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao compilar Assinador em modo desenvolvedor: %s", strings.TrimSpace(string(output)))
	}
	return classesDir, nil
}

// orquestra a execução da assinatura/validação.
// tenta usar o modo servidor (HTTP) e regride para modo local (JAR direto) se necessário.
func ExecJavaSigner(cmdKey string, fileName string, options ExecutionOptions) (string, error) {
	if err := options.validate(); err != nil {
		return "", err
	}
	if options.Local {
		return execLocalSigner(cmdKey, fileName)
	}
	if err := EnsureServerRunningWithOptions(options); err != nil {
		return "", fmt.Errorf("modo servidor indisponível: %w; use --local para execução direta", err)
	}
	return CallJavaServer(cmdKey, cmdKey, fileName, options)
}

func execLocalSigner(cmdKey string, fileName string) (string, error) {
	requestData := map[string]string{
		"command": cmdKey,
		"file":    fileName,
		"flag":    "--file",
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("falha ao criar requisição JSON: %w", err)
	}

	javaPath, err := GetJavaPath("java")
	if err != nil {
		return "", err
	}

	var javaCmd *exec.Cmd
	if os.Getenv("DEV_MODE") == "true" {
		LogFeedback("ASSINATURA CONFIG", "Modo Desenvolvedor Ativo (Lendo App.java).")
		classesDir, err := compileDevAssinador()
		if err != nil {
			return "", err
		}
		javaCmd = exec.Command(javaPath, "-cp", classesDir, "App", string(jsonData))
	} else {
		localJarPath := GetJarPath()
		if _, err := os.Stat(localJarPath); os.IsNotExist(err) {
			if err := DownloadAssinadorJar(localJarPath); err != nil {
				return "", err
			}
		}
		javaCmd = exec.Command(javaPath, "-jar", localJarPath, string(jsonData))
	}

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
