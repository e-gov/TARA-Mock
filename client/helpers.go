package main

import (
	// "math/rand"
	"net/http"
)

// randSeq tagastab tähtedest koosneva juhusõne, pikkusega n.
/* func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
*/

// getParameter tagastab päringu r parameetri p väärtuse; parameetri puudumisel tagastab "".
func getP(p string, r *http.Request) string {
	if v, ok := r.Form[p]; ok {
		// Parameeter võib korduda. Võtame esimese.
		return v[0]
	}
	return ""
}
