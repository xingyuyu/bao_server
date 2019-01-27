package token

import (
	"log"
	"testing"
)

func Test_GetUnionID(t *testing.T) {
	id := GetUnionID("oi_kx1TKih6PG3mczng5sXULzZBM", "18_h4VYFZ5fx0XBSxPEJuAdjREmGhAp-uQvl9iotyWOOON_iQfAtsOMT5H3q666RTNRsOWu3VTfW6hDYb4T0Hft3sKkT6kdbRezOGaCw1QCnFz--9N3fMWto0N4F1EIJGjABAYZW")
	log.Println("union id=", id)
}
