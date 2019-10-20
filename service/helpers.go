package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
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
	p := filepath.Join("templates", "index.html")
	t, err := template.ParseFiles(p)
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
		return
	}
	t.Execute(w, nil)

}

// healthCheck pakub elutukset (/health).
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"name":"TARA-Mock", "status":"ok"}`)
}

// personCodeToDoB tagastab isikukoodi põhjal arvutatud sünnikuupäeva.
func personCodeToDoB(c string) (dob string, err error) {
	if len(c) > 6 {
		// Leia sajand
		var s string
		switch string(c[0]) {
		case "1", "2":
			s = "18"
		case "3", "4":
			s = "19"
		case "5", "6":
			s = "20"
		default:
			return "", errors.New("Sajand vale")
		}
		dob = s + c[1:3] + "-" + c[3:5] + "-" + c[5:7]
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
