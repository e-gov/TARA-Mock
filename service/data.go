package main

// Identity on autenditav identiteet (isik).
type Identity struct {
	Isikukood     string
	Eesnimi       string
	Perekonnanimi string
}

// identities hoiab kasutajale valimiseks pakutavaid identiteete (isikuid).
var identities []Identity

func init() {

	// Kasutajale valikuks pakutavad konkreetsed identiteedid (isikud).
	identities = append(identities,
		Identity{"Isikukood1", "Eesnimi1", "Perekonnanimi1"},
		Identity{"Isikukood2", "Eesnimi2", "Perekonnanimi2"},
		Identity{"Isikukood3", "Eesnimi3", "Perekonnanimi3"},
	)

	// TO DO: Võimalik on kompaktsemalt
}
