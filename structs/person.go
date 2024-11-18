package structs

type Person struct {
	FirstName string  `json:"vorname"`
	LastName  string  `json:"nachname"`
	BirthDate string  `json:"geburtsdatum"`
	Sex       string  `json:"geschlecht"`
	RvNr      string  `json:"rentenversicherungsnummer"`
	Picture   Picture `json:"picture"`
}
