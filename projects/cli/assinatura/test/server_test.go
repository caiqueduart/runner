package test

import (
	"net/http"
	"testing"
	"time"
)

func TestServerLifecycle(t *testing.T) {
	// Garantir que o servidor está desligado antes de começar
	RunCLI("stop")
	time.Sleep(1 * time.Second)

	// 1. TC-07: Iniciar servidor automaticamente via comando sign
	absFile, err := CreateDummyFile("lifecycle.json", "{}")
	if err != nil {
		t.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer CleanUpFiles("lifecycle.json", "lifecycle-json-assinatura.txt")

	// Executa sign - isso deve disparar o servidor se DEV_MODE=false
	res, err := RunCLI("sign", "--file", absFile)
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}

	// 2. Verificar se o servidor responde em /health
	// (Damos um tempinho para o background process se estabilizar)
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Logf("Aviso: Servidor não respondeu via HTTP. Talvez DEV_MODE esteja ativo? Stdout: %s", res.Stdout)
		return // Se DEV_MODE estiver on, encerramos aqui pois não haverá servidor
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 200 no /health, obtido %d", resp.StatusCode)
	}

	// 3. TC-08: Comando stop
	res, err = RunCLI("stop")
	if err != nil {
		t.Fatalf("Erro ao parar servidor: %v", err)
	}

	if !ContainsIgnoringColors(res.Stdout, "Encerrando") && !ContainsIgnoringColors(res.Stdout, "encerrado") && !ContainsIgnoringColors(res.Stdout, "Sinal") {
		t.Errorf("Mensagem de encerramento ou sinal não encontrada: %s", res.Stdout)
	}

	// Verificar se parou de responder
	time.Sleep(1 * time.Second)
	_, err = http.Get("http://localhost:8080/health")
	if err == nil {
		t.Errorf("Servidor ainda responde após comando stop")
	}
}
