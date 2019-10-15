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
	http.HandleFunc("/return", requestIdentityToken)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // fileServer serveerib kasutajaliidese muutumatuid faile.

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
	// Loe avalehe mall, täida ja saada sirvikusse.
	p := filepath.Join("templates", "index.html")
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

	// Päringuparameetrid, lehele "Autenditud" edastamiseks
	type ReturnParams struct {
		Code  string
		State string
		Nonce string
	}
	var ps ReturnParams

	r.ParseForm() // Parsi päringuparameetrid.
	// Võta päringust volituskood, state ja nonce
	ps.Code = getP("code", r)
	ps.State = getP("state", r)
	ps.Nonce = getP("nonce", r)

	// TO DO: Päri identsustõend

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Loe lehe "Autenditud" vmall, täida ja saada sirvikusse.
	p := filepath.Join("templates", "autenditud.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, ps)

}
