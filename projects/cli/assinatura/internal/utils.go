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

// Caminho para o arquivo Main.java (temporário)
const jarPath = "../../assinador/src/Main.java"

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
	if binName == "java" || binName == "java.exe" || binName == "javac" || binName == "javac.exe" {
		fmt.Println("JDK 21 não encontrado ou versão incompatível.")
		fmt.Println("Iniciando download automático do JDK 21...")

		if err := DownloadJava21(managedDir); err != nil {
			return "", fmt.Errorf("falha ao baixar JDK 21: %w", err)
		}

		fmt.Println("Download e instalação do JDK 21 concluídos com sucesso.")

		if _, err := os.Stat(managedBin); err == nil {
			return managedBin, nil
		}
	}

	return "", fmt.Errorf("binário %s não encontrado", binName)
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

	javacPath, err := GetJavaPath("javac")
	if err != nil {
		return "", err
	}

	// Verifica a existência do código fonte
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		return "", fmt.Errorf("arquivo do Assinador não encontrado em: %s", jarPath)
	}

	// Preparação do Assinador
	fmt.Println("Preparando o Assinador...")

	absJarPath, _ := filepath.Abs(jarPath)
	javaProjectRoot := filepath.Dir(filepath.Dir(absJarPath))

	// Criamos um diretório 'bin' temporário para as classes compiladas
	binDir := filepath.Join(javaProjectRoot, "bin")
	os.MkdirAll(binDir, 0755)

	// Listamos os fontes necessários recursivamente a partir de 'src'
	var sources []string
	err = filepath.Walk(filepath.Join(javaProjectRoot, "src"), func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".java") {
			sources = append(sources, path)
			return nil
		}
		return err
	})

	if len(sources) == 0 {
		return "", fmt.Errorf("nenhum arquivo .java encontrado em %s", filepath.Join(javaProjectRoot, "src"))
	}

	// Compilação
	compileCmd := exec.Command(javacPath, append([]string{"-d", binDir}, sources...)...)
	if output, err := compileCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("erro na preparação do Assinador: %w\n%s", err, string(output))
	}

	// Execução
	fmt.Println("Executando Assinador...")

	javaCmd := exec.Command(javaPath, "-cp", binDir, "src.Main", cmdKey, fileName)
	javaCmd.Dir = javaProjectRoot

	output, err := javaCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro na execução: %w\n%s", err, string(output))
	}

	return string(output), nil
}
