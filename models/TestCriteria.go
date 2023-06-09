package models

type TestCriteria struct {
	NamePattern    *string   `json:"namePattern"`
	UnitTypeValues *[]string `json:"unitTypeValues"`
}
