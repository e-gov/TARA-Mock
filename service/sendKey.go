package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	// JWK (veebivõtme) toiminguteks
	// Dok-n: https://godoc.org/github.com/lestrrat-go/jwx
	"github.com/lestrrat-go/jwx/jwk"
)

// sendKey väljastab klientrakendusele identsustõendi allkirjastamisel
// kasutatavale privaatvõtmele vastava avaliku võtme.
// sendKey teostab OpenID Connect avaliku võtme otspunkti oidc/jwks.
func sendKey(w http.ResponseWriter, r *http.Request) {
	key, err := jwk.New(&signKey.PublicKey)
	if err != nil {
		log.Printf("Viga JWK moodustamisel: %s", err)
		return
	}
	// key sisaldab: {"kty":..., "n":..., "e":...}

	// Moodusta võtmehulk, vrdl:
	// https://tara-test.ria.ee/oidc/jwks
	type keySet struct {
		Keys []jwk.Key `json:"keys"`
	}
	var ks keySet
	ks.Keys = append(ks.Keys, key)

	jsonbuf, err := json.Marshal(ks)
	if err != nil {
		log.Printf("Viga JWK serialiseerimisel: %s", err)
		return
	}

	// Lisa kid, teatava häkina. Seda oleks pidanud tegema juba varem,
	// kuid tüüpidega on raskusi.
	nk := []byte(
		strings.Replace(string(jsonbuf), "[{", "[{\"kid\":\"taramock\",", 1),
	)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(nk)
}

// Vt näide: https://github.com/lestrrat-go/jwx
