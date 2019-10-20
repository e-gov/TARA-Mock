package main

import (
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
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
