package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {

	// Marsruudid
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", landingPage)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/client/return", requestIdentityToken)

	// Käivita HTTPS server
	log.Println("** Klientrakenduse näidis käivitatud pordil 8081 **")
	err := http.ListenAndServeTLS(
		":8081",
		"keys/https-server.crt",
		"keys/https-server.key",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// healthCheck pakub elutukset (/health).
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"name":"TARA-Mock klientrakendus", "status":"ok"}`)
}

// LandingPage on klientrakenduse avaleht; kasutaja saab seal sisse logida (/).
func landingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p := filepath.Join("ui", "index.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, nil)
}

// LoginUser suunab kasutaja TARA-Mock-i autentima.
func loginUser(w http.ResponseWriter, r *http.Request) {
	// Ümbersuunamis-URL
	ru := "https://localhost:8080/oidc/authorize?=" +
		"redirect_uri=https%3A%2F%2Flocalhost%3A8081%2Freturn&" +
		"scope=openid&" +
		"state=1111&" +
		"response_type=code&" +
		"client_id=1"

	// Suuna kasutaja TARA-Mock-i
	http.Redirect(w, r, ru, 301)
}

// RequestIdentityToken võtab TARA-Moc-st tagasi suunatud kasutaja vastu ja pärib TARA-Mock-lt identsustõendi (otspunkt /client/return).
func requestIdentityToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Kuva saadud volituskood
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)
	fmt.Fprint(w, "Klientrakendus - küsi identsustõend")
}
