package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed templates/*
var content embed.FS

const (
	Author = "Adam Sidor"
	Port   = "8080"
)

func main() {
	// 1. Definicja flagi healthcheck
	isHealthCheck := flag.Bool("health", false, "Sprawdź stan serwera")
	flag.Parse()

	// 2. Logika Healthcheck: jeśli aplikacja jest uruchomiona z flagą -health
	if *isHealthCheck {
		// Próba połączenia z lokalnym serwerem
		client := http.Client{Timeout: 2 * time.Second}
		_, err := client.Get("http://localhost:" + Port)
		if err != nil {
			// Jeśli serwer nie odpowiada, kończymy z błędem (exit code 1)
			os.Exit(1)
		}
		// Jeśli odpowiedział, kończymy sukcesem (exit code 0)
		os.Exit(0)
	}

	// 3. Normalny start aplikacji (wyświetlanie logów przy starcie)
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