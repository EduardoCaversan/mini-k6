package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mini-k6/handlers"

	"github.com/gorilla/mux"
)

func main() {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err := os.Mkdir("uploads", 0755)
		if err != nil {
			log.Fatalf("Erro ao criar diretório de uploads: %v", err)
		}
	}

	r := mux.NewRouter()

	r.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	r.HandleFunc("/run-test", handlers.RunTestHandler).Methods("POST")
	r.HandleFunc("/progress", handlers.ProgressStream)
	r.HandleFunc("/summary", handlers.SummaryHandler).Methods("POST")
	r.HandleFunc("/report", handlers.ReportHandler).Methods("POST")

	fmt.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
