package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

// authenticateUser võtab vastu klientrakendusest autentimisele
// saadetud kasutaja ja kuvab talle autentimisdialoogi avalehe
// (dialoogis ongi üks leht). Kasutaja saabub päringuga otspunkti
// /oidc/authorize.
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
	// Kuva kontrolliks mäpi Form kõik elemendid
	fmt.Println("authenticateUser: Autentimispäringu parameetrid:")
	for k, v := range r.Form {
		fmt.Printf("  %s: %s\n", k, v)
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
	p := filepath.Join("templates", "authenticateUser.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, mp)
}
