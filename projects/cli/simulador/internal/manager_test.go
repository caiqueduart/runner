package internal

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsPortAvailable(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("falha ao reservar porta de teste: %v", err)
	}

	port := fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
	if IsPortAvailable(port) {
		t.Fatal("porta ocupada foi reportada como disponível")
	}

	if err := listener.Close(); err != nil {
		t.Fatalf("falha ao liberar porta de teste: %v", err)
	}
	if !IsPortAvailable(port) {
		t.Fatal("porta liberada foi reportada como ocupada")
	}
}

func TestGetSimuladorStatus(t *testing.T) {
	server := newSimuladorTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/info" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"HubSaude","status":"ready"}`)
	})
	defer server.Close()

	status, err := GetSimuladorStatus()
	if err != nil {
		t.Fatalf("status deveria ser obtido: %v", err)
	}
	if status != `{"name":"HubSaude","status":"ready"}` {
		t.Fatalf("status inesperado: %s", status)
	}
}

func TestGetSimuladorStatusRejectsHTTPError(t *testing.T) {
	server := newSimuladorTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "indisponível", http.StatusServiceUnavailable)
	})
	defer server.Close()

	if _, err := GetSimuladorStatus(); err == nil {
		t.Fatal("status HTTP de erro deveria ser propagado")
	}
}

func TestStopSimulador(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	server := newSimuladorTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/shutdown" {
			http.Error(w, "requisição inválida", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	if err := StopSimulador(); err != nil {
		t.Fatalf("parada deveria ser aceita: %v", err)
	}
}

func newSimuladorTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	server := httptest.NewTLSServer(handler)
	previousClient := httpClient
	previousBaseURL := simuladorBaseURL
	httpClient = server.Client()
	simuladorBaseURL = server.URL

	t.Cleanup(func() {
		httpClient = previousClient
		simuladorBaseURL = previousBaseURL
	})

	return server
}
