package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLifecycleCommands(t *testing.T) {
	shutdownCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/info":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"name":"HubSaude","version":"0.1.7"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/shutdown":
			shutdownCalled = true
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	start := runCLI(t, server.URL, "start")
	if start.exitCode != 0 || !strings.Contains(start.stdout, "já está em execução") {
		t.Fatalf("start falhou (código %d): stdout=%s stderr=%s", start.exitCode, start.stdout, start.stderr)
	}

	status := runCLI(t, server.URL, "status")
	if status.exitCode != 0 || !strings.Contains(status.stdout, "Simulador Online") || !strings.Contains(status.stdout, `"version":"0.1.7"`) {
		t.Fatalf("status falhou (código %d): stdout=%s stderr=%s", status.exitCode, status.stdout, status.stderr)
	}

	stop := runCLI(t, server.URL, "stop")
	if stop.exitCode != 0 || !strings.Contains(stop.stdout, "parada enviado com sucesso") {
		t.Fatalf("stop falhou (código %d): stdout=%s stderr=%s", stop.exitCode, stop.stdout, stop.stderr)
	}
	if !shutdownCalled {
		t.Fatal("comando stop não chamou POST /shutdown")
	}
}

func TestVersionCommand(t *testing.T) {
	result := runCLI(t, "http://127.0.0.1", "version")
	if result.exitCode != 0 {
		t.Fatalf("version retornou código %d: %s", result.exitCode, result.stderr)
	}
	if !strings.Contains(result.stdout, "Simulador CLI v1.1.0") || !strings.Contains(result.stdout, "Simulador JAR v0.1.7") {
		t.Fatalf("versões esperadas não encontradas: %s", result.stdout)
	}
}

func TestStatusReturnsFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "indisponível", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	result := runCLI(t, server.URL, "status")
	if result.exitCode == 0 {
		t.Fatalf("status offline deveria falhar: stdout=%s", result.stdout)
	}
	if !strings.Contains(result.stdout, "Offline ou erro") || !strings.Contains(result.stderr, "status 503") {
		t.Fatalf("diagnóstico inesperado: stdout=%s stderr=%s", result.stdout, result.stderr)
	}
}

func TestStopReturnsFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	result := runCLI(t, server.URL, "stop")
	if result.exitCode == 0 || !strings.Contains(result.stderr, "status 503") {
		t.Fatalf("stop deveria propagar erro HTTP: stdout=%s stderr=%s", result.stdout, result.stderr)
	}
}

func TestCommandsRejectUnexpectedArguments(t *testing.T) {
	result := runCLI(t, "http://127.0.0.1", "version", "extra")
	if result.exitCode == 0 || !strings.Contains(result.stderr, "unknown command") {
		t.Fatalf("argumento inesperado deveria falhar: stdout=%s stderr=%s", result.stdout, result.stderr)
	}
}
