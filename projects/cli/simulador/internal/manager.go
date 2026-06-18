package internal

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

var simuladorBaseURL = fmt.Sprintf("https://localhost:%s", SimuladorPort)

func simuladorURL(path string) string {
	baseURL := simuladorBaseURL
	if configuredURL := os.Getenv("SIMULADOR_BASE_URL"); configuredURL != "" {
		baseURL = strings.TrimRight(configuredURL, "/")
	}
	return baseURL + path
}

func IsPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return false
	}

	ln.Close()

	return true
}

// localiza o executável do Java 21. Se não encontrado no sistema ou na pasta
// gerenciada, inicia o download automático do JDK.
func GetJavaPath(binName string) (string, error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(binName, ".exe") {
		binName += ".exe"
	}

	managedBin := filepath.Join(GetJDKDir(), "bin", binName)

	// tenta encontrar no PATH do sistema
	if path, err := exec.LookPath(binName); err == nil {
		if isJava21(path) {
			return path, nil
		}
	}

	// tenta encontrar na pasta .hubsaude
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

// verifica se o binário apontado é de fato a versão 21 do Java.
func isJava21(javaPath string) bool {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false
	}

	return strings.Contains(string(output), "version \"21")
}

// baixa o JDK 21 para o diretório alvo.
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

	LogFeedback("SIMULADOR CONFIG", "Extraindo arquivos...")

	if strings.HasSuffix(downloadUrl, ".zip") {
		return extractZip(tmpFile, targetDir)
	}

	return extractTarGz(tmpFile, targetDir)
}

// baixa o JAR do simulador diretamente do GitHub do professor.
func DownloadSimuladorJar(targetPath string) error {
	LogFeedback("SIMULADOR CONFIG", "JAR do simulador não encontrado. Baixando versão %s...", CompatibleSimuladorVersion)

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

// garante que o simulador esteja ativo. Se não estiver, inicia-o em background
// após validar a porta e provisionar os arquivos necessários.
func EnsureSimuladorRunning() error {
	// verifica se já está respondendo via HTTPS
	url := simuladorURL("/api/info")
	resp, err := httpClient.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		resp.Body.Close()
		LogFeedback("SIMULADOR CONFIG", "Simulador já está em execução.")
		return nil
	}

	// valida se a porta 8443 está livre conforme requisito US-03
	if !IsPortAvailable(SimuladorPort) {
		return fmt.Errorf("a porta %s já está em uso por outro processo. Não é possível iniciar o simulador", SimuladorPort)
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

	// inicia o processo Java de forma independente (background)
	cmd := exec.Command(javaPath, "-jar", localJarPath)
	SetDetachedProcess(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("falha ao iniciar simulador: %w", err)
	}

	if cmd.Process != nil {
		savePID(cmd.Process.Pid, SimuladorPort)
		cmd.Process.Release()
	}

	// aguarda o servidor subir e responder ao endpoint de saúde
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
	url := simuladorURL("/shutdown")
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
	url := simuladorURL("/api/info")
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

func DownloadDriver() error {
	ext := ".so"
	if runtime.GOOS == "windows" {
		ext = ".dll"
	}

	driverName := "libscuba" + ext
	targetPath := filepath.Join(GetHubSaudeDir(), "bin", driverName)

	LogFeedback("SIMULADOR CONFIG", "Baixando driver PKCS#11 do simulador...")
	url := simuladorURL("/api/driver")

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("falha ao conectar no simulador: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("simulador retornou erro %d ao solicitar driver", resp.StatusCode)
	}

	os.MkdirAll(filepath.Dir(targetPath), 0755)
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	LogFeedback("SIMULADOR CONFIG", "Driver %s salvo em %s", driverName, targetPath)

	return nil
}
