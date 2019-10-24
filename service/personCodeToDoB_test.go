package main

import (
	"testing"
)

func TestPersonCodeToDoB(t *testing.T) {

	// Liiga lühike
	dob1, err := personCodeToDoB("361071")
	if err == nil {
		t.Errorf("<361071>: ootasin: %v, %v, sain: %v, %v",
			"", nil,
			dob1, err)
	}

	// Õige, eesliitega
	dob2, err := personCodeToDoB("EE36107120334") // Priit Parmakson
	if err != nil || dob2 != "1961-07-12" {
		t.Errorf("<36107120334>: ootasin: %v, %v, sain: %v, %v", "1961-07-12", nil,
			dob2, err)
	}

	// Õige, eesliiteta
	dob2, err := personCodeToDoB("36107120334") // Priit Parmakson
	if err != nil || dob2 != "1961-07-12" {
		t.Errorf("<36107120334>: ootasin: %v, %v, sain: %v, %v", "1961-07-12", nil,
			dob2, err)
	}

	// Kuu või päeva nr väljaspool lubatud vahemikku
	dob3, err := personCodeToDoB("36107130334") // Priit Parmakson
	if err != nil {
		t.Errorf("<36107130334>: ootasin: %v, %v, sain: %v, %v",
			"", nil,
			dob3, err)
	}

	// Ei vasta kuupäevavormingule
	dob4, err := personCodeToDoB("36N07130334") // Priit Parmakson
	if err == nil {
		t.Errorf("<36N07130334>: ootasin: %v, %v, sain: %v, %v",
			"", "<veateade>",
			dob4, err)
	}
}
