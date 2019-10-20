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

	// TO DO: Lahenda tagasisuunamis-URL-i pordi edastamise probleem

	r.ParseForm() // Parsi päringuparameetrid.

	returnURI := getPtr("redirect_uri", r)
	state := getPtr("state", r)
	nonce := getPtr("nonce", r)

	// Genereeri volituskood
	var c volituskood
	c = volituskood(randSeq(6))

	// Kogu identsustõendi koostamiseks ja väljastamiseks vajalikud
	// andmed..
	var dataForToken dataForTokenType

	// Selgita, millise identiteedi kasutaja valis. Kui valis etteantute
	// hulgast, siis Form submit saatis elemendi isik=nr (0-based).
	isikunr := getPtr("isik", r)

	if isikunr != "" {
		// Teisenda int-ks
		i, err := strconv.Atoi(isikunr)
		if err != nil {
			fmt.Println("sendUserBack: Viga vormist saadud andmete kasutamisel")
			i = 0 // Kasuta esimest etteantud identiteeti
		}
		dataForToken.sub = identities[i].Isikukood
		dataForToken.givenName = identities[i].Eesnimi
		dataForToken.familyName = identities[i].Perekonnanimi
	} else {
		// Kasutaja ei valinud etteantud identiteetide seast, vaid
		// sisestas identiteedi ise.
		dataForToken.sub = getPtr("idcode", r)
		dataForToken.givenName = getPtr("firstname", r)
		dataForToken.familyName = getPtr("lastname", r)
	}

	dataForToken.state = state
	dataForToken.nonce = nonce

	// ..ja pane tallele
	mutex.Lock()
	idToendid[c] = dataForToken
	mutex.Unlock()

	fmt.Println("sendUserBack: Id-tõendi andmed talletatud: ", dataForToken)

	// Moodusta tagasisuunamis-URL
	ru := returnURI +
		"?code=" + string(c) +
		"&state=" + state +
		"&nonce=" + nonce

	// Suuna kasutaja tagasi
	http.Redirect(w, r, ru, 301)
}
