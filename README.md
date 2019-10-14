# TARA-Mock
TARA teenust etendav makett | TARA mock-up, a testing tool

Rakendus etendab TARA autentimist ja on sellisena kasutatav klientrakenduste testimisel, kui TARA testteenuse võimalused jäävad napiks. TARA testteenus võimaldab autentida piiratud hulga testkasutajatena. TARA-Mock:
- ühendatakse klientrakenduse külge täpselt niisamuti nagu TARA külge
- autentimise dialoogis võimaldab kasutajal valida identiteet kas etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida.

Käivitamine: `go run cmd/TARA-Mock/main.go`

*Seadistamine* on minimaalne. Sea failis `main.go` järgmised väärtused:
- host, kuhu TARA-Mock on paigaldatud
- klientrakendusse tagasisuunamise aadress

```
const taraMockHost = "localhost"
const returnURL = "https://localhost:8080/client/return"
```
Samuti genereeri ja sea TARA-Mock HTTPS serveri serdid. Sertide genereerimise näiteskript on failis `keys/genkeys.sh`.

*Paigaldamine*. Paigalda TARA-Mock sobivasse masinasse. TARA-Mock on kasutatav ka oma masinas (`localhost`).

Kasutamine (oma masinas):
- `https://localhost:8080/health` - elutukse
- `https://localhost:8080/` - avaleht teabega TARA-Makett kohta
- `https://localhost:8080/oidc/authorize` - autentimisele suunamine
- `https://localhost:8080/oidc/token` - identsustõendi väljastamine
- `https://localhost:8080/oidc/jwks` - identsustõendi avalik võti

