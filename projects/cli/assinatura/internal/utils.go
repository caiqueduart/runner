package internal

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CompatibleAssinadorVersion define a versão do JAR que esta CLI sabe operar.
const CompatibleAssinadorVersion = "0.1.4"
const RepoPath = "caiqueduart/runner"

func PrintError(format string, a ...any) {
	fmt.Printf(ColorRed+format+ColorReset, a...)
}

// GetJavaPath localiza ou baixa o JDK 21 e retorna o caminho para o binário solicitado.
func GetJavaPath(binName string) (string, error) {
	if runtime.GOOS == "windows" && !strings.HasSuffix(binName, ".exe") {
		binName += ".exe"
	}

	// Tenta encontrar no PATH ou no local gerenciado
	home, _ := os.UserHomeDir()
	managedDir := filepath.Join(home, ".hubsaude", "jdk")
	managedBin := filepath.Join(managedDir, "bin", binName)

	// Verifica se o java no PATH é o 21
	if binName == "java" || binName == "java.exe" {
		if path, err := exec.LookPath(binName); err == nil {
			if isJava21(path) {
				return path, nil
			}
		}
	}

	// Verifica no diretório gerenciado
	if _, err := os.Stat(managedBin); err == nil {
		// Se for o java, valida a versão
		if binName == "java" || binName == "java.exe" {
			if isJava21(managedBin) {
				return managedBin, nil
			}
		} else {
			return managedBin, nil
		}
	}

	// Se não encontrou (ou versão errada), inicia download
	if binName == "java" || binName == "java.exe" {
		fmt.Println("JDK 21 não encontrado ou versão incompatível.")
		fmt.Println("Iniciando download do JDK 21...")

		if err := DownloadJava21(managedDir); err != nil {
			return "", fmt.Errorf("Falha ao baixar JDK 21: %w", err)
		}

		fmt.Println("Download e instalação do JDK 21 concluídos.")

		if _, err := os.Stat(managedBin); err == nil {
			return managedBin, nil
		}
	}

	return "", fmt.Errorf("Binário %s não encontrado", binName)
}

// isJava21 verifica se o binário java informado é da versão 21.
func isJava21(javaPath string) bool {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "version \"21")
}

// DownloadAssinadorJar baixa o JAR do assinador de uma release do GitHub.
func DownloadAssinadorJar(targetPath string) error {
	tag := "assinador-v" + CompatibleAssinadorVersion
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", RepoPath, tag)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Falha ao consultar release no GitHub: status %d", resp.StatusCode)
	}

	var release struct {
		Assets []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return err
	}

	var jarDownloadURL string
	expectedName := fmt.Sprintf("assinador-v%s.jar", CompatibleAssinadorVersion)
	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			jarDownloadURL = asset.DownloadURL
			break
		}
	}

	if jarDownloadURL == "" {
		return fmt.Errorf("Arquivo %s não encontrado na release %s", expectedName, tag)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err = http.Get(jarDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// DownloadJava21 baixa e extrai o JDK 21 da Adoptium.
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
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return err
	}
	if len(releases) == 0 {
		return fmt.Errorf("nenhum release encontrado na API da Adoptium")
	}

	downloadUrl := releases[0]["binaries"].([]interface{})[0].(map[string]interface{})["package"].(map[string]interface{})["link"].(string)
	fmt.Printf("Baixando JDK de: %s\n", downloadUrl)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	tmpFile := filepath.Join(os.TempDir(), "jdk21_download"+filepath.Ext(downloadUrl))
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	resp, err = http.Get(downloadUrl)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return err
	}

	fmt.Println("Extraindo arquivos...")
	if strings.HasSuffix(downloadUrl, ".zip") {
		return extractZip(tmpFile, targetDir)
	}
	return extractTarGz(tmpFile, targetDir)
}

// Funções de extração simplificadas
func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	defer r.Close()

	var rootFolder string

	if len(r.File) > 0 {
		rootFolder = strings.Split(r.File[0].Name, "/")[0]
	}

	for _, f := range r.File {
		fpath := filepath.Join(dest, strings.TrimPrefix(f.Name, rootFolder))
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		outFile, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		rc, _ := f.Open()
		io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}

	return nil
}

func extractTarGz(src, dest string) error {
	f, _ := os.Open(src)

	defer f.Close()

	gzr, _ := gzip.NewReader(f)

	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var rootFolder string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if rootFolder == "" {
			rootFolder = strings.Split(header.Name, "/")[0]
		}
		fpath := filepath.Join(dest, strings.TrimPrefix(header.Name, rootFolder))
		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(fpath, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(fpath), 0755)
			outFile, _ := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			io.Copy(outFile, tr)
			outFile.Close()
		}
	}

	return nil
}

func ExecJavaSigner(fileName string, cmdKey string) (string, error) {
	// Prepara o ambiente Java
	javaPath, err := GetJavaPath("java")
	if err != nil {
		return "", err
	}

	// Local do JAR gerenciado
	home, _ := os.UserHomeDir()
	jarDir := filepath.Join(home, ".hubsaude", "bin")
	jarName := fmt.Sprintf("assinador-v%s.jar", CompatibleAssinadorVersion)
	localJarPath := filepath.Join(jarDir, jarName)

	// Verifica se o JAR existe, se não, baixa
	if _, err := os.Stat(localJarPath); os.IsNotExist(err) {
		fmt.Printf("Assinador v%s não encontrado localmente.\n", CompatibleAssinadorVersion)
		fmt.Println("Iniciando download do Assinador...")
		if err := DownloadAssinadorJar(localJarPath); err != nil {
			return "", fmt.Errorf("falha ao baixar Assinador: %w", err)
		}
		fmt.Println("Download do Assinador concluído.")
	}

	// Execução
	fmt.Println("Executando Assinador...")

	javaCmd := exec.Command(javaPath, "-jar", localJarPath, cmdKey, fileName)
	output, err := javaCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro na execução: %w\n%s", err, string(output))
	}

	return string(output), nil
}
