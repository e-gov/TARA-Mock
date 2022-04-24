package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func main() {

	cFilePtr := flag.String("conf", "config.json", "Seadistusfaili asukoht")
	flag.Parse()

	// Loe seadistus sisse.
	conf = loadConf(*cFilePtr)
	log.Infoln("* Klientrakenduse näidis: Seadistus loetud")

	// Marsruudid
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/", landingPage)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/autologin", autologinUser)
	http.HandleFunc("/return", finalize)

	// fileServer serveerib kasutajaliidese muutumatuid faile.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Käivita HTTPS server.
	log.Infoln("** Klientrakenduse näidis: Käivitatud pordil 8081")
	err := http.ListenAndServeTLS(
		conf.AppPort,
		conf.AppCert,
		conf.AppKey,
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// LandingPage on klientrakenduse avaleht; kasutaja saab seal sisse logida (/).
func landingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Valmista ette malli parameetrid.
	type MalliParameetrid struct {
		appHost     string
		appPort     string
		RedirectURI string
	}
	mp := MalliParameetrid{conf.AppHost, conf.AppPort, conf.RedirectURI}

	// Loe avalehe mall, täida ja saada sirvikusse.
	p := filepath.Join("templates", "index.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, mp)
}

// loginUser suunab kasutaja TARA-Mock-i autentima.
func loginUser(w http.ResponseWriter, r *http.Request) {
	// Ümbersuunamis-URL
	ru := conf.TaraMockAuthorizeEndpoint + "?" +
		"redirect_uri=" +
		url.PathEscape(conf.RedirectURI) + "&" +
		"scope=openid&" +
		"state=1111&" +
		"nonce=2222&" +
		"response_type=code&" +
		"client_id=1"

	log.Infof("\nloginUser:\n    Saadan autentimispäringu:\n    %v\n", ru)

	// Suuna kasutaja TARA-Mock-i.
	http.Redirect(w, r, ru, http.StatusMovedPermanently) // 301
}

// autologinUser suunab kasutaja TARA-Mock-i automaatautentimisele.
// F-n erib loginUser-st ainult parameetri autologin=<isikukood>
// poolest. TO DO: Kaalu refaktoorimist.
func autologinUser(w http.ResponseWriter, r *http.Request) {
	// Ümbersuunamis-URL
	ru := conf.TaraMockAuthorizeEndpoint + "?" +
		"redirect_uri=" +
		url.PathEscape(conf.RedirectURI) + "&" +
		"scope=openid&" +
		"state=1111&" +
		"nonce=2222&" +
		"response_type=code&" +
		"client_id=1&" +
		"autologin=36107120334"

	log.Infof("\nautologinUser:\n    Saadan autentimispäringu:\n    %v\n", ru)

	// Suuna kasutaja TARA-Mock-i.
	http.Redirect(w, r, ru, http.StatusMovedPermanently) // 301
}

// finalize : 1) võtab TARA-Moc-st tagasi suunatud kasutaja
// vastu; 2) kutsub välja identsustõendi pärimise; 3) viib sisselogimise
// lõpule - saadab sirvikusse lehe "autenditud". (Otspunkt /client/return).
func finalize(w http.ResponseWriter, r *http.Request) {

	// PassParams koondab lehele "Autenditud" edastatavaid väärtusi.
	type PassParams struct {
		Code        string
		State       string
		Nonce       string
		Isikuandmed string
		Success     bool
	}
	var ps PassParams

	r.ParseForm() // Parsi päringuparameetrid.
	// Võta päringust volituskood, state ja nonce
	ps.Code = getP("code", r)
	ps.State = getP("state", r)
	ps.Nonce = getP("nonce", r)

	// Päri identsustõend
	// t []byte - Identsustõend
	t, ok := getIdentityToken(getP("code", r))
	if !ok {
		log.Infoln("\nfinalize: Identsustõendi pärimine ebaõnnestus")
		ps.Success = false
	} else {
		log.Printf("\nfinalize:\n    Saadud identsustõend:\n    %v\n", string(t))
		ps.Success = true
	}

	ps.Isikuandmed = t

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Loe lehe "Autenditud" mall, täida ja saada sirvikusse.
	p := filepath.Join("templates", "autenditud.html")
	tpl, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	tpl.Execute(w, ps)

}

// healthCheck pakub elutukset (/health).
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"name":"TARA-Mock klientrakendus", "status":"ok"}`)
}
