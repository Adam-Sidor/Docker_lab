package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func getIP() string {
	ip := "127.0.0.1"
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if !v.IP.IsLoopback() && v.IP.To4() != nil {
						ip = v.IP.String()
					}
				}
			}
		}
	}
	return ip
}

func handler(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	ip := getIP()
	
	version := os.Getenv("VERSION")
	if version == "" {
		version = "Brak podanej wersji"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="pl">
<head>
    <meta charset="utf-8">
    <title>Laboratorium 5</title>
</head>
<body>
    <h1>Informacje o serwerze</h1>
    <p><strong>Adres IP kontenera:</strong> %s</p>
    <p><strong>Nazwa serwera (hostname):</strong> %s</p>
    <p><strong>Wersja aplikacji:</strong> %s</p>
</body>
</html>`, ip, hostname, version)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Serwer aplikacji Go wystartował na porcie 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Błąd uruchomienia serwera: %v\n", err)
	}
}
