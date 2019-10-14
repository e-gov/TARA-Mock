# TARA-Mock
TARA teenust etendav makett | TARA mock-up, a testing tool

TARA-Mock on rakendus, mis etendab TARA autentimist. 

**Kasutusstsenaarium**. TARA-Mock on mõeldud kasutamiseks siis, kui TARA testteenuse võimalused jäävad klientrakenduse funktsionaalsuste testimisel napiks. Kui TARA testteenus võimaldab autentida väga väikese hulga TARA poolt ette antud testkasutajatega, siis TARA-Mock võimaldab kasutajal valida autentimise dialoogis identiteet kas TARA-Mock seadistuses etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida s.t logida sisse suvalise identiteediga.

 TARA-Mock ühendatakse klientrakenduse külge täpselt niisamuti nagu TARA külge

**Käivitamine** `go run cmd/TARA-Mock/main.go`

**Seadistamine**. Sea failis `main.go` järgmised väärtused:
- host, kuhu TARA-Mock on paigaldatud
- klientrakendusse tagasisuunamise aadress

```
const taraMockHost = "localhost"
const returnURL = "https://localhost:8080/client/return"
```
Samuti genereeri ja sea TARA-Mock HTTPS serveri serdid. Sertide genereerimise näiteskript on failis `keys/genkeys.sh`.

Etteantud identiteedid on koodi sissekirjutatud TARA-Mock koodi. Muuda identiteedid oma vajadustele vastavaks.

**Paigaldamine**. Paigalda TARA-Mock sobivasse masinasse. TARA-Mock on kasutatav ka oma masinas (`localhost`).

**Kasutamine (oma masinas)**:
- `https://localhost:8080/health` - elutukse
- `https://localhost:8080/` - avaleht teabega TARA-Makett kohta
- `https://localhost:8080/oidc/authorize` - autentimisele suunamine
- `https://localhost:8080/oidc/token` - identsustõendi väljastamine
- `https://localhost:8080/oidc/jwks` - identsustõendi avalik võti

