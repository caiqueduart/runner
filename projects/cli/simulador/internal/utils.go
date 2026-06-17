package internal

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// retorna o caminho para a pasta oculta .hubsaude no home do usuário,
// usada para centralizar binários, JDK e arquivos de estado.
func GetHubSaudeDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".hubsaude")
}

// retorna o caminho completo onde o simulador.jar deve ser armazenado.
func GetSimuladorJarPath() string {
	return filepath.Join(GetHubSaudeDir(), "bin", SimuladorJarName)
}

// retorna o caminho para a instalação gerenciada do JDK 21.
func GetJDKDir() string {
	return filepath.Join(GetHubSaudeDir(), "jdk")
}

// retorna o caminho do arquivo .pid usado para rastrear o processo em background.
func GetPIDFilePath() string {
	return filepath.Join(GetHubSaudeDir(), "simulador.pid")
}

func PrintError(format string, a ...any) {
	fmt.Printf(ColorRed+format+ColorReset, a...)
}

func checkFileSHA256(filePath string, expectedDigest string) (bool, error) {
	expectedHash := strings.TrimPrefix(expectedDigest, "sha256:")
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}

	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}

	calculatedHash := hex.EncodeToString(hash.Sum(nil))

	return strings.EqualFold(calculatedHash, expectedHash), nil
}

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
