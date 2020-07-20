package main

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randSeq tagastab tähtedest koosneva juhusõne, pikkusega n.
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// getPtr tagastab päringuga r esitatud vormi või
// URL-parameetri p väärtuse. Parameetri puudumisel tagastab "".
func getPtr(p string, r *http.Request) string {
	if v, ok := r.Form[p]; ok {
		// Parameeter võib korduda. Võtame esimese.
		return v[0]
	}
	return ""
}

// landingPage annab teavet TARA-Mock-i kohta (/).
// Siinne vastus saadetakse ka päringutele, mida ei suudeta
// marsruutida.
func landingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	// Loe avalehe mall, täida ja saada sirvikusse.
	t, err := template.ParseFiles(conf.IndexTmpl)
	if err != nil {
		log.Errorf("Unable to load template: %v", err)
		return
	}
	t.Execute(w, nil)
}

// sendConf saadab OpenID Connect seadistuse (otspunkt
// .well-known/openid-configuration).
func sendConf(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	type Conf struct {
		Issuer                           string   `json:"issuer"`
		ScopesSupported                  []string `json:"scopes_supported"`
		ResponseTypesSupported           []string `json:"response_types_supported"`
		SubjectTypesSupported            []string `json:"subject_types_supported"`
		ClaimTypesSupported              []string `json:"claim_types_supported"`
		ClaimsSupported                  []string `json:"claims_supported"`
		GrantTypesSupported              []string `json:"grant_types_supported"`
		IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
		UILocalesSupported               []string `json:"ui_locales_supported"`
		TokenEndpoint                    string   `json:"token_endpoint"`
		UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
		AuthorizationEndpoint            string   `json:"authorization_endpoint"`
		JwksURI                          string   `json:"jwks_uri"`
	}

	conf := Conf{
		Issuer:                           "https://" + conf.TaraMockHost,
		ScopesSupported:                  []string{"openid", "idcard", "mid", "banklink", "smartid", "eidas", "eidasonly", "email"},
		ResponseTypesSupported:           []string{"code"},
		SubjectTypesSupported:            []string{"public", "pairwise"},
		ClaimTypesSupported:              []string{"normal"},
		ClaimsSupported:                  []string{"sub", "given_name", "family_name", "date_of_birth", "email", "email_verified"},
		GrantTypesSupported:              []string{"authorization_code"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
		UILocalesSupported:               []string{"et", "en", "ru"},
		TokenEndpoint:                    "https://" + conf.TaraMockHost + conf.HTTPServerPort + "/oidc/token",
		UserinfoEndpoint:                 "https://" + conf.TaraMockHost + conf.HTTPServerPort + "/oidc/profile",
		AuthorizationEndpoint:            "https://" + conf.TaraMockHost + conf.HTTPServerPort + "/oidc/authorize",
		JwksURI:                          "https://" + conf.TaraMockHost + conf.HTTPServerPort + "/oidc/jwks",
	}

	json.NewEncoder(w).Encode(conf)
}

// personCodeToDoB tagastab isikukoodi põhjal arvutatud sünnikuupäeva.
// Eeldab, et isikukood antakse eesliitega "EE". Kui isikukoodis eesliidet "EE"
// ei leia, siis üritab sünnikuupäeva arvutada eesliiteta. err-is tagastab
// veateate või edu korral nil.
func personCodeToDoB(c string) (dob string, err error) {
	r := []rune(c)
	log.Debugf("First 2", string(r[0:2]))
	if len(c) > 6 {
		// Kas eesliide, milles riigikood?
		var e []rune
		if string(r[0:2]) == "EE" {
			e = r[2:]
		} else {
			e = r
		}
		// Leia sajand
		var s string
		switch string(e[0]) {
		case "1", "2":
			s = "18"
		case "3", "4":
			s = "19"
		case "5", "6":
			s = "20"
		default:
			return "", errors.New("Sajand vale")
		}
		dob = s + string(e[1:3]) + "-" + string(e[3:5]) + "-" + string(e[5:7])
		dobc := dob + "T15:04:05+00:00"
		// Kontrolli, kas legaalne kp
		// RFC3339 - "2012-11-01T22:08:41+00:00"
		_, err := time.Parse(time.RFC3339, dobc)
		if err != nil {
			return "", errors.New("Illegaalne kp")
		}
		return dob, nil
	}
	return "", errors.New("Isikukood liiga lühike")
}

// healthCheck pakub elutukset (/health).
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"name":"TARA-Mock", "status":"ok"}`)
}
