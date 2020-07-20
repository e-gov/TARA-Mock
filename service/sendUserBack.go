package main

import (
	log "github.com/sirupsen/logrus"
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
	log.Debugf("Generated authorization code/Genereeritud volituskood: %v", c)

	// Kogu identsustõendi koostamiseks ja väljastamiseks vajalikud
	// andmed.
	var forToken forTokenType

	// Selgita, millise identiteedi kasutaja valis. Kui valis etteantute
	// hulgast, siis Form submit saatis elemendi isik=<nr> (0-based).
	isikunr := getPtr("isik", r)
	log.Debugf("    Kasutaja valis isiku: %v", isikunr)

	if isikunr != "" {
		// Teisenda int-ks
		i, err := strconv.Atoi(isikunr)
		if err != nil {
			log.Errorf("    Viga vormist saadud andmete kasutamisel: %v", err)
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

	log.Debugf("--- Id-tõendi andmed talletatud: %+v", forToken)

	// Moodusta tagasisuunamis-URL
	ru := returnURI +
		"?code=" + string(c) +
		"&state=" + state +
		"&nonce=" + nonce

	log.Debugf("--- Suunan kasutaja tagasi: %v", ru)

	// Suuna kasutaja tagasi
	http.Redirect(w, r, ru, 301)
}
