package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"mini-k6/internal"
	"mini-k6/models"
	"mini-k6/results"
)

var (
	clients   = make(map[chan string]bool)
	clientsMu sync.Mutex
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := saveUploadedFile(file, header)
	if err != nil {
		http.Error(w, "Erro ao salvar o arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Arquivo salvo em: %s\n", dst)
}

func saveUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	dstPath := filepath.Join("uploads", header.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	return dstPath, err
}

func RunTestHandler(w http.ResponseWriter, r *http.Request) {
	var scenario models.TestScenario

	if err := json.NewDecoder(r.Body).Decode(&scenario); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("▶️  Recebido cenário: %+v\n", scenario)

	results := internal.ExecuteScenario(scenario)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Erro ao codificar resposta JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProgressStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string)

	clientsMu.Lock()
	clients[messageChan] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, messageChan)
		clientsMu.Unlock()
		close(messageChan)
	}()

	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		clientsMu.Lock()
		delete(clients, messageChan)
		clientsMu.Unlock()
		close(messageChan)
	}()

	for msg := range messageChan {
		fmt.Fprintf(w, "data: %s\n\n", msg)
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

func BroadcastProgress(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for ch := range clients {
		select {
		case ch <- message:
		default:
			delete(clients, ch)
		}
	}
}

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	var scenario models.TestScenario
	if err := json.NewDecoder(r.Body).Decode(&scenario); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	results := internal.ExecuteScenario(scenario)
	summary := summarizeResults(results)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, "Erro ao codificar resumo JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

type Summary struct {
	Total             int         `json:"total"`
	Success           int         `json:"success"`
	Failures          int         `json:"failures"`
	AverageSeconds    float64     `json:"average_seconds"`
	TotalSeconds      float64     `json:"total_seconds"`
	RequestsPerSecond float64     `json:"requests_per_second"`
	ByStatusCode      map[int]int `json:"by_status_code"`
}

func summarizeResults(results []results.TestResult) Summary {
	var sumTime float64
	success := 0
	fail := 0
	statusMap := map[int]int{}

	for _, r := range results {
		sumTime += r.Duration.Seconds()
		statusMap[r.StatusCode]++
		if r.Error == "" {
			success++
		} else {
			fail++
		}
	}

	avg := 0.0
	total := len(results)
	if total > 0 {
		avg = sumTime / float64(total)
	}

	reqPerSec := 0.0
	if sumTime > 0 {
		reqPerSec = float64(total) / sumTime
	}

	return Summary{
		Total:             total,
		Success:           success,
		Failures:          fail,
		AverageSeconds:    avg,
		TotalSeconds:      sumTime,
		RequestsPerSecond: reqPerSec,
		ByStatusCode:      statusMap,
	}
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	var scenario models.TestScenario
	if err := json.NewDecoder(r.Body).Decode(&scenario); err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	results := internal.ExecuteScenario(scenario)
	summary := summarizeResults(results)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, "Erro ao codificar JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
