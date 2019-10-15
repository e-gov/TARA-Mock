package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

const taraMockHost = "localhost"
const returnURL = "https://localhost:8081/return"

func main() {

	// Marsruudid
	// In Go the pattern "/" matches all request paths, rather than just the empty path.
	http.HandleFunc("/", LandingPage)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/oidc/authorize", AuthenticateUser)
	http.HandleFunc("/back", SendUserBack)
	http.HandleFunc("/oidc/token", SendIdentityToken)
	http.HandleFunc("/oidc/jwks", SendKey)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs)) // fileServer serveerib kasutajaliidese muutumatuid faile.

	// Käivita HTTPS server
	log.Println("** TARA-Mock käivitatud pordil 8080 **")
	err := http.ListenAndServeTLS(
		":8080",
		"keys/https-server.crt",
		"keys/https-server.key",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// LandingPage annab teavet TARA-Mock-i kohta (/).
func LandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "TARA-Mock")
}

// AuthenticateUser etendab kasutaja autentimise dialoogi (/oidc/authorize).
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// OidcParams hoiab Klientrakendusest saadetud päringu OIDC parameetreid.
	type OidcParams struct {
		RedirectURI  string // redirect_uri
		Scope        string
		State        string
		ResponseType string // response_type
		ClientID     string // client_id
		UILocales    string // ui_locales
		Nonce        string
		AcrValues    string // acr_values
	}

	r.ParseForm() // Parsi päringuparameetrid.

	// pr hoiab päringuparameetreid; kasutatakse malli täitmiseks.
	var pr OidcParams

	// Valmista päringuparameetrid ette malli täitmiseks.
	pr.RedirectURI = getP("redirect_uri", r)
	pr.Scope = getP("scope", r)
	pr.State = getP("state", r)
	pr.ResponseType = getP("response_type", r)
	pr.ClientID = getP("client_id", r)
	pr.UILocales = getP("ui_locales", r)
	pr.Nonce = getP("nonce", r)
	pr.AcrValues = getP("acr_values", r)

	// Loe mall, täida ja saada leht sirvikusse
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p := filepath.Join("templates", "index.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, pr)
}

// SendUserBack : 1) võtab sirvikust vastu kasutaja tehtud valiku;
// 2) genereerib OIDC volituskoodi; 3) genereerib identsustõendi ja
// paneb selle ootele ning 4) saadab kasutaja klientrakendusse tagasi
// (/return).
func SendUserBack(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Parsi päringuparameetrid.

	// TO DO: Lahenda valitud isikuandmete edastamine. Peidetud väljadega?
	fmt.Println(getP("nr", r))
	fmt.Println(getP("idcode", r))
	fmt.Println(getP("firstname", r))
	fmt.Println(getP("lastname", r))

	// TO DO: Lahenda pordi edastamise probleem
	// redirectURI := getP("redirect_uri", r)
	state := getP("state", r)
	nonce := getP("nonce", r)

	// Genereeri volituskood
	c := randSeq(5)
	// TO DO: Moodusta identsustõend ja pane ootele

	// Moodusta tagasisuunamis-URL
	ru := returnURL +
		"?code=" + c +
		"&state=" + state +
		"&nonce=" + nonce

	// Suuna kasutaja tagasi
	http.Redirect(w, r, ru, 301)
}

// SendIdentityToken väljastab klientrakendusele identsustõendi (OIDC identsustõendi otspunkt /oidc/token).
func SendIdentityToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "TARA-Mock - token")
}

// SendKey väljastab klientrakendusele identsustõendi allkirjastamisel kasutatavale privaatvõtmele vastava avaliku võtme (sellega teostab OIDC avaliku võtme otspunkti oidc/jwks).
func SendKey(w http.ResponseWriter, r *http.Request) {

}

// healthCheck pakub elutukset (/health).
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"name":"TARA-Mock", "status":"ok"}`)
}
