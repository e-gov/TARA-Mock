# TARA-Mock

TARA-Mock on rakendus, mis etendab TARA autentimist. 

TARA-Mock on mõeldud kasutamiseks siis, kui TARA testteenuse võimalused jäävad klientrakenduse funktsionaalsuste testimisel napiks. TARA testteenuse abil saab autentida väga väikese hulga TARA poolt ette antud testkasutajatega. TARA-Mock seevastu võimaldab kasutajal valida autentimise dialoogis identiteet kas TARA-Mock seadistuses etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida s.t logida sisse suvalise identiteediga.

TARA-Mock pakub samu otspunkte kui TARA, s.t TARA-Mock ja TARA on üksteisega vahetatavad lihtsa seadistusega.

TARA-Mock juures on ka klientrakenduse näidis.

## Kasutusstsenaarium

1) Suunamine autentimisele

---
<img src="docs/Rakendus_01.PNG" width="400">
---

2) Autentimisdialoog

---

![](docs/TARA-Mock_01.PNG)

---

3) Tagasi autentimiselt

---

![](docs/Rakendus_02.PNG)

---

## TARA-Mock

TARA-Mock töötab pordil 8080.

**Käivitamine**.

```
cd service
go run .
```

**Seadistamine**.

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

## Klientrakenduse näidis

**Otspunktid**:
- `/health` - elutukse
- `/` - avaleht; kasutaja saab sealt minna TARA-Mock-i autentima
- `/login` - kasutaja suunamine TARA-Mock-i autentima
- `/return` - autentimiselt tagasi suunatud kasutaja vastuvõtmine, identsustõendi pärimine TARA-Mock-st ja sisselogimise lõpuleviimine 

**Seadistamine**: Klientrakenduse näidis töötab lokaalses masinas, pordil 8081. Genereeri ja sea TARA-Mock HTTPS serveri serdid. Sertide genereerimise näiteskript on failis `keys/genkeys.sh`.

**Käivitamine**:

```
cd service
go run .
```
