# Etteantud identiteedid

TARA-Mock-is määrab kasutaja ise identiteedi (isikukoodi, ees- ja perekonnanime), millega ta autenditakse. Selleks ta kas valib etteantud identiteetide hulgast või sisestab ise identiteeti.

Tarkvaraga on kaasas 3 etteantud identiteeti. Etteantud identiteetide muutmiseks redigeeri faili `service/data.go`:

```
...
		Identity{"Isikukood1", "Eesnimi1", "Perekonnanimi1"},
		Identity{"Isikukood2", "Eesnimi2", "Perekonnanimi2"},
		Identity{"Isikukood3", "Eesnimi3", "Perekonnanimi3"},
...
```
