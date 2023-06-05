package models

type Test struct {
	TestName           string   `json:"testName"`
	UnitType           string   `json:"unitType"`
	References         []string `json:"references"`
	Standards          []string `json:"standards"`
	AvailableModifiers []string `json:"availableModifiers"`
}
