package main

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

// authenticateUser võtab vastu klientrakendusest autentimisele
// saadetud kasutaja ja kuvab talle autentimisdialoogi avalehe
// (dialoogis ongi üks leht). Kasutaja saabub päringuga otspunkti
// /oidc/authorize.
// Kui päringus on parameeter autologin=<isikukood>, siis tehakse
// automaatautentimine. Kui isikukood on ettantud identiteetide
// hulgas, siis tagastatakse vastav ees- ja perekonnanimi; vastasel
// korral eesnimi: Auto, perekonnanimi Maat.
func authenticateUser(w http.ResponseWriter, r *http.Request) {
	// OidcParams hoiab klientrakendusest saadetud päringu
	// OpenID Connect kohaseid parameetreid.
	type OidcParams struct {
		RedirectURI  string // redirect_uri
		Scope        string // scope
		State        string // state
		ResponseType string // response_type
		ClientID     string // client_id
		UILocales    string // ui_locales
		Nonce        string // nonce
		AcrValues    string // acr_values
	}

	r.ParseForm() // Parsi päringuparameetrid.
	if log.IsLevelEnabled(log.DebugLevel) {
		// Kuva kontrolliks mäpi Form kõik elemendid
		log.Debug("Authentication request parameters/Autentimispäringu parameetrid:")
		for k, v := range r.Form {
			log.Debugf("    %s: %s", k, v)
		}
	}

	// Automaatautentimine?
	if v, ok := r.Form["autologin"]; ok {
		// Parameeter võib korduda. Võtame esimese.
		ik := v[0]
		// Järgnevas on ühisosa sendUserBack-ga. TO DO: Refaktoori

		// Genereeri volituskood
		var c volituskood
		c = volituskood(randSeq(6))

		// Kogu identsustõendi koostamiseks ja väljastamiseks vajalikud
		// andmed.
		var forToken forTokenType

		forToken.sub = ik
		forToken.givenName = "Auto"
		forToken.familyName = "Maat"

		// Kas autologin parameetris näidatud isik on etteantute hulgas?
		for _, identity := range identities {
			if identity.Isikukood == ik {
				forToken.givenName = identity.Eesnimi
				forToken.familyName = identity.Perekonnanimi
				break
			}
		}

		log.WithFields(log.Fields{
			"sub": forToken.sub,
			"givenName": forToken.givenName,
			"familyName": forToken.familyName,
		}).Debug("Automaatautentimine")

		forToken.clientID = getPtr("client_id", r)
		forToken.state = getPtr("state", r)
		forToken.nonce = getPtr("nonce", r)

		// ..ja pane tallele
		mutex.Lock()
		idToendid[c] = forToken
		mutex.Unlock()

		log.WithFields(log.Fields{
			"token": forToken,
		}).Debug("Id-tõendi andmed talletatud")

		// Moodusta tagasisuunamis-URL
		ru := getPtr("redirect_uri", r) +
			"?code=" + string(c) +
			"&state=" + getPtr("state", r) +
			"&nonce=" + getPtr("nonce", r)

		// Suuna kasutaja tagasi
		http.Redirect(w, r, ru, 301)
	}

	// pr hoiab päringuparameetreid; kasutatakse malli täitmiseks.
	var pr OidcParams

	// Valmista päringuparameetrid ette malli täitmiseks.
	pr.RedirectURI = getPtr("redirect_uri", r)
	pr.Scope = getPtr("scope", r)
	pr.State = getPtr("state", r)
	pr.ResponseType = getPtr("response_type", r)
	pr.ClientID = getPtr("client_id", r)
	pr.UILocales = getPtr("ui_locales", r)
	pr.Nonce = getPtr("nonce", r)
	pr.AcrValues = getPtr("acr_values", r)

	// Valmista ette malli parameetrid. Mallile saadetakse päringu-
	// parameetrid (taustateabeks) ja etteantud identiteedid (isikud,
	// kelle hulgast kasutaja saab valida sobiva.
	type templateParams struct {
		Request    OidcParams
		Identities []Identity
	}
	mp := templateParams{
		Request:    pr,
		Identities: identities,
	}

	// Loe mall, täida parameetritega ja saada leht sirvikusse.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(conf.AuthenticateUserTmpl)
	if err != nil {
		log.Debug("Unable to load template")
		return
	}
	t.Execute(w, mp)
}
