package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

const taraMockHost = "localhost"
const returnURL = "https://localhost:8081/return"

type volituskood string

// Andmed identsustõendi moodustamiseks ja väljastamiseks. Identsustõend koostatakse
// vahetult enne väljastamist.
type idtAndmed struct {
	sub        string // subject, isikutõendi väli "sub"
	familyName string // family_name
	givenName  string // given_name
}

// Identsustõendite hoidla
var idToendid = make(map[volituskood]idtAndmed)

var mutex = &sync.Mutex{}

// IdentityToken on väljastatav identsustõendi struktuur (hetkel
// mittetäielik).
type IdentityToken struct {
	Sub       string `json:"sub"`
	FirstName string `json:"given_name"`
	LastName  string `json:"first_name"`
}

// IDTokenReqBody on vastuvõetava identsustõendi päringu keha
// struktuur.
type IDTokenReqBody struct {
	GrantType  string      `json:"grant_type"`
	Code       volituskood `json:"code"`
	RequestURI string      `json:"request_uri"`
}

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
		"vault/https-server.crt",
		"vault/https-server.key",
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// LandingPage annab teavet TARA-Mock-i kohta (/).
// Siinne vastus saadetakse ka päringutele, mida ei suudeta
// marsruutida.
func LandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, "TARA-Mock")
}

// AuthenticateUser võtab vastu klientrakendusest autentimisele
// saadetud kasutaja ja kuvab talle autentimisdialoogi avalehe
// (dialoogis ongi üks leht). (Kasutaja saabub päringuga otspunkti
//	/oidc/authorize).
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

	// TO DO: Lahenda pordi edastamise probleem
	// redirectURI := getP("redirect_uri", r)
	state := getP("state", r)
	nonce := getP("nonce", r)

	// Genereeri volituskood
	var c volituskood
	c = volituskood(randSeq(6))

	// Kogu identsustõendi koostamiseks ja väljastamiseks vajalikud
	// andmed..
	var idta idtAndmed
	switch n := getP("nr", r); n {
	case "1":
		i := idtAndmed{"Isikukood1", "Eesnimi1", "Perekonnanimi1"}
		idta = i
	case "2":
		i := idtAndmed{"Isikukood2", "Eesnimi2", "Perekonnanimi2"}
		idta = i
	case "3":
		i := idtAndmed{"Isikukood3", "Eesnimi3", "Perekonnanimi3"}
		idta = i
	default:
		idta.sub = getP("idcode", r)
		idta.givenName = getP("firstname", r)
		idta.familyName = getP("lastname", r)
	}

	// ..ja pane hoidlasse
	mutex.Lock()
	idToendid[c] = idta
	mutex.Unlock()

	fmt.Println("Id-tõendi andmed hoiustatud: ", idta)

	// Moodusta tagasisuunamis-URL
	ru := returnURL +
		"?code=" + string(c) +
		"&state=" + state +
		"&nonce=" + nonce

	// Suuna kasutaja tagasi
	http.Redirect(w, r, ru, 301)
}

// SendIdentityToken väljastab klientrakendusele identsustõendi (OIDC identsustõendi otspunkt /oidc/token).
func SendIdentityToken(w http.ResponseWriter, r *http.Request) {
	// Võta päringust volituskood, seejärel võta volituskoodile
	// vastavad identsustõendi andmed, koosta identsustõend ja saada
	// päringu vastuses

	// Loe päringu keha
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Paki lahti
	var msg IDTokenReqBody
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// b.Code
	// idToendid[b.Code]
	// idtAndmed
	v, ok := idToendid[msg.Code]
	if !ok {
		http.Error(w, "Identsustõend ei eksisteeri", 404)
		return
	}

	// Koosta identsustõend
	var mt IdentityToken
	mt.Sub, mt.FirstName, mt.LastName = v.sub, v.givenName, v.familyName
	output, err := json.Marshal(mt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Väljasta
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// SendKey väljastab klientrakendusele identsustõendi allkirjastamisel kasutatavale privaatvõtmele vastava avaliku võtme (eostab OIDC avaliku võtme otspunkti oidc/jwks).
func SendKey(w http.ResponseWriter, r *http.Request) {

}

// healthCheck pakub elutukset (/health).
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"name":"TARA-Mock", "status":"ok"}`)
}
