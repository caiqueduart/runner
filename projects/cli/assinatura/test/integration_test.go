package test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSignValidateFlow(t *testing.T) {
	// Limpeza inicial
	CleanUpFiles("test_receita.json", "test_receita-json-assinatura.txt")
	defer CleanUpFiles("test_receita.json", "test_receita-json-assinatura.txt")

	// 1. Criar arquivo para assinar
	absFile, err := CreateDummyFile("test_receita.json", `{"item": "vacina"}`)
	if err != nil {
		t.Fatalf("Erro ao criar arquivo de teste: %v", err)
	}

	// 2. Testar Sign Sucesso
	res, err := RunCLI("sign", "--file", absFile)
	if err != nil {
		t.Fatalf("Erro ao executar sign: %v", err)
	}

	if res.ExitCode != 0 {
		t.Errorf("Esperado exit code 0, obtido %d. Stdout: %s", res.ExitCode, res.Stdout)
	}

	if !ContainsIgnoringColors(res.Stdout, "gerou o código de assinatura") {
		t.Errorf("Mensagem de sucesso não encontrada no stdout: %s", res.Stdout)
	}

	// O arquivo é gerado no CWD do comando
	generatedFile := "test_receita-json-assinatura.txt"
	absGenerated, _ := filepath.Abs(generatedFile)

	if _, err := os.Stat(absGenerated); os.IsNotExist(err) {
		t.Errorf("Arquivo '%s' não foi gerado", absGenerated)
	}

	// 3. Testar Validate Sucesso
	res, err = RunCLI("validate", absGenerated)
	if err != nil {
		t.Fatalf("Erro ao executar validate: %v", err)
	}

	if res.ExitCode != 0 {
		t.Errorf("Esperado exit code 0 no validate, obtido %d. Stdout: %s", res.ExitCode, res.Stdout)
	}

	if !ContainsIgnoringColors(res.Stdout, "está assinado sob o código") {
		t.Errorf("Mensagem de validação não encontrada no stdout: %s", res.Stdout)
	}
}

func TestSignValidationErrors(t *testing.T) {
	// TC-01: Sign sem flag --file
	res, err := RunCLI("sign")
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}
	// O exit code vem do JAR (ou do wrapper que detectou falha no JAR)
	// Como a validação é no JAR, esperamos que ele retorne erro
	if !ContainsIgnoringColors(res.Stdout, "Erro do usuário: O parâmetro '--file' é obrigatório") {
		t.Errorf("Esperado erro de parâmetro obrigatório, obtido: %s", res.Stdout)
	}

	// TC-02: Flag inválida
	res, err = RunCLI("sign", "--f", "test.txt")
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}
	if !ContainsIgnoringColors(res.Stdout, "Você quis dizer '--file'?") {
		t.Errorf("Esperada sugestão de flag correta, obtido: %s", res.Stdout)
	}

	// TC-06: Arquivo inexistente
	res, err = RunCLI("sign", "--file", "non_existent.json")
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}
	if !ContainsIgnoringColors(res.Stdout, "não encontrado") {
		t.Errorf("Esperado erro de arquivo não encontrado, obtido: %s", res.Stdout)
	}
}

func TestValidateErrors(t *testing.T) {
	// TC-03: Validate sem argumento
	res, err := RunCLI("validate")
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}
	if !ContainsIgnoringColors(res.Stdout, "Forneça o caminho do arquivo para validação") {
		t.Errorf("Esperado erro de falta de argumento, obtido: %s", res.Stdout)
	}
}
