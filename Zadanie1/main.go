package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Osadzanie plików HTML wewnątrz binarki
//go:embed templates/*
var content embed.FS

const (
	Author = "Adam Sidor"
	Port   = "8080"
)

func main() {
	// Informacje w logach przy uruchomieniu
	startTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Data uruchomienia: %s\n", startTime)
	fmt.Printf("Autor programu:    %s\n", Author)
	fmt.Printf("Nasłuchiwanie na:  TCP %s\n", Port)

	// Parsowanie szablonów z osadzonego systemu plików
	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	// Obsługa głównej strony
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// Uruchomienie serwera
	serverErr := http.ListenAndServe(":"+Port, nil)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}