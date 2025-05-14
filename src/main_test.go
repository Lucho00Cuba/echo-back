package main

import (
	"errors"
	"net/http"
	"testing"
)

func TestStartServer_InvalidPort(t *testing.T) {
	// Guardar valores originales
	originalPort := PORT
	defer func() { PORT = originalPort }()
	PORT = "not_a_number"

	// Capturar os.Exit
	calledExit := false
	originalExit := osExit
	defer func() { osExit = originalExit }()
	osExit = func(code int) {
		calledExit = true
		panic("os.Exit called")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from os.Exit, but none occurred")
		}
		if !calledExit {
			t.Error("expected osExit to be called, but it wasn't")
		}
	}()

	startServer(dummyListen)
}

func TestStartServer_Success(t *testing.T) {
	// Guardar valores originales
	originalPort := PORT
	defer func() { PORT = originalPort }()
	PORT = "8080"

	started := false

	mockListen := func(addr string, handler http.Handler) error {
		started = true
		if addr != ":8080" {
			t.Errorf("expected addr :8080, got %s", addr)
		}
		// simulate success
		return nil
	}

	startServer(mockListen)

	if !started {
		t.Error("expected server to be started")
	}
}

func TestStartServer_Failure(t *testing.T) {
	originalPort := PORT
	defer func() { PORT = originalPort }()
	PORT = "8081"

	mockListen := func(addr string, handler http.Handler) error {
		return errors.New("mock failure")
	}

	startServer(mockListen)
}

// Dummy listener used to prevent actual server start in invalid port test
func dummyListen(addr string, handler http.Handler) error {
	return nil
}
