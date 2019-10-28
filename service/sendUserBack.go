package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// sendUserBack : 1) võtab sirvikust vastu kasutaja tehtud valiku
// (Form submit);
// 2) genereerib OIDC volituskoodi;
// 3) kogub identsustõendi koostamiseks vajalikud andmed ja talletab
// need mälus peetavas tõendihoidlas;
// 4) saadab kasutaja klientrakendusse tagasi.
// sendUserBack teostab otspunkti /return.
func sendUserBack(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // Parsi päringuparameetrid.

	returnURI := getPtr("redirect_uri", r)
	state := getPtr("state", r)
	nonce := getPtr("nonce", r)
	clientID := getPtr("client_id", r)

	// Genereeri volituskood
	var c volituskood
	c = volituskood(randSeq(6))
	fmt.Printf("sendUserBack:\n    Genereeritud volituskood: %v\n", c)

	// Kogu identsustõendi koostamiseks ja väljastamiseks vajalikud
	// andmed.
	var forToken forTokenType

	// Selgita, millise identiteedi kasutaja valis. Kui valis etteantute
	// hulgast, siis Form submit saatis elemendi isik=<nr> (0-based).
	isikunr := getPtr("isik", r)
	fmt.Printf("    Kasutaja valis isiku: %v\n", isikunr)

	if isikunr != "" {
		// Teisenda int-ks
		i, err := strconv.Atoi(isikunr)
		if err != nil {
			fmt.Println("    Viga vormist saadud andmete kasutamisel")
			i = 0 // Kasuta esimest etteantud identiteeti
		}
		forToken.sub = identities[i].Isikukood
		forToken.givenName = identities[i].Eesnimi
		forToken.familyName = identities[i].Perekonnanimi
	} else {
		// Kasutaja ei valinud etteantud identiteetide seast, vaid
		// sisestas identiteedi ise.
		forToken.sub = getPtr("idcode", r)
		forToken.givenName = getPtr("firstname", r)
		forToken.familyName = getPtr("lastname", r)
	}

	forToken.clientID = clientID
	forToken.state = state
	forToken.nonce = nonce

	// ..ja pane tallele
	mutex.Lock()
	idToendid[c] = forToken
	mutex.Unlock()

	fmt.Printf("--- Id-tõendi andmed talletatud:\n    %+v\n", forToken)

	// Moodusta tagasisuunamis-URL
	ru := returnURI +
		"?code=" + string(c) +
		"&state=" + state +
		"&nonce=" + nonce

	fmt.Printf("--- Suunan kasutaja tagasi:\n    %v\n", ru)

	// Suuna kasutaja tagasi
	http.Redirect(w, r, ru, 301)
}
