package webhook

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Climatik-Project/Climatik-Project/internal/webhook/factory"
)

// Make sure this interface is defined in the factory package and imported here

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	runnerType := r.URL.Query().Get("runner")
	path := r.URL.Query().Get("path")
	if source == "" || runnerType == "" || path == "" {
		http.Error(w, "Source, runner, and path query parameters are required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// runnerFactory := &runners.RunnerFactory{}
	// runner, err := runnerFactory.GetRunner(runnerType, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// runner.Run()

	handlerFactory := &factory.AlertHandlerFactory{}
	handler, err := handlerFactory.GetHandler(source)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.HandleAlert(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alert received and handled"))
}

func CreateWebhook(port int) {
	http.HandleFunc("/alert", AlertHandler)
	portStr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
