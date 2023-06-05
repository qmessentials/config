package models

type Unit struct {
	FullName          string `json:"fullName"`
	FullNamePlural    string `json:"fullNamePlural"`
	Abbreviation      string `json:"abbreviation"`
	MeasurementSystem string `json:"measurementSystem"`
	UnitType          string `json:"unitType"`
}
