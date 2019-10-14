package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

const taraMockHost = "localhost"
const returnURL = "https://localhost:8080/client/return"

func main() {

	// Marsruudid
	http.HandleFunc("/", LandingPage)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/oidc/authorize", AuthenticateUser)
	http.HandleFunc("/back", SendUserBack)
	http.HandleFunc("/oidc/token", SendIdentityToken)
	http.HandleFunc("/oidc/jwks", SendKey)
	http.HandleFunc("/client/return", RequestIdentityToken)

	// Käivita HTTPS server
	log.Println("** TARA-Makett käivitatud pordil 8080 **")
	err := http.ListenAndServeTLS(
		":8080",
		"keys/https-server.crt",
		"keys/https-server.key",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// healthCheck pakub elutukset.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

// LandingPage annab teavet TARA-Mock-i kohta.
func LandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "TARA-Mock")
}

// AuthenticateUser etendab kasutaja autentimise dialoogi.
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p := filepath.Join("ui", "index.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}

	user := user{ID: 1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
		Phone: "000099999"}
	t.Execute(w, user)
}

// SendUserBack saadab kasutaja klientrakendusse tagasi.
func SendUserBack(w http.ResponseWriter, r *http.Request) {
	// Genereeri volituskood
	c := PseudoUUID()
	// Moodusta identsustõend ja pane ootele

	// Moodusta tagasisuunamis-URL
	ru := returnURL + "?code=" + c
	// Suuna kasutaja tagasi
	http.Redirect(w, r, ru, 301)
}

// SendIdentityToken teostab OIDC identsustõendi otspunkti.
func SendIdentityToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "TARA-Mock - token")
}

// RequestIdentityToken teostab klientrakenduses identsustõendi küsimist
func RequestIdentityToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Kuva saadud volituskood
	fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)
	fmt.Fprint(w, "Klientrakendus - küsi identsustõend")
}

// SendKey teostab OIDC avaliku võtme otspunkti.
func SendKey(w http.ResponseWriter, r *http.Request) {

}
