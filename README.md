# TARA-Mock
TARA teenust etendav makett | TARA mock-up, a testing tool

TARA-Mock on rakendus, mis etendab TARA autentimist. 

**Kasutusstsenaarium**. TARA-Mock on mõeldud kasutamiseks siis, kui TARA testteenuse võimalused jäävad klientrakenduse funktsionaalsuste testimisel napiks. TARA testteenuse abil saab autentida väga väikese hulga TARA poolt ette antud testkasutajatega. TARA-Mock seevastu võimaldab kasutajal valida autentimise dialoogis identiteet kas TARA-Mock seadistuses etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida s.t logida sisse suvalise identiteediga.

TARA-Mock pakub samu otspunkte kui TARA, s.t TARA-Mock ja TARA on üksteisega vahetatavad lihtsa seadistusega.

**Käivitamine**. `go run cmd/TARA-Mock/main.go`

**Seadistamine**. Sea failis `main.go` järgmised väärtused:
- `returnURL` - klientrakendusse tagasisuunamise aadress

Nt, kui TARA-Mock ja klientrakendus on samas masinas:
```
const returnURL = "https://localhost:8080/client/return"
```
Genereeri ja sea TARA-Mock HTTPS serveri serdid. Sertide genereerimise näiteskript on failis `keys/genkeys.sh`.

Etteantud identiteedid on koodi sissekirjutatud TARA-Mock koodi. Muuda identiteedid oma vajadustele vastavaks.

**Paigaldamine**. Paigalda TARA-Mock sobivasse masinasse. TARA-Mock on kasutatav ka oma masinas (`localhost`).

**Otspunktid**:
- `/health` - elutukse
- `/` - avaleht teabega TARA-Makett kohta
- `/oidc/authorize` - autentimisele suunamine
- `/token` - identsustõendi väljastamine
- `/oidc/jwks` - identsustõendi avalik võti

Nt, TARA-Mock kasutamisel oma masinas:

`https://localhost:8080/oidc/token` - identsustõendi väljastamine
