package models

type User struct {
	UserID                   string   `json:"userId"`
	GivenNames               []string `json:"givenNames"`
	FamilyNames              []string `json:"familyNames"`
	Roles                    []string `json:"roles"`
	EmailAddress             string   `json:"emailAddress"`
	IsActive                 bool     `json:"isActive"`
	IsPasswordChangeRequired bool     `json:"isPasswordChangeRequired"`
	HashedPassword           string   `json:"-"`
	AuthToken                string   `json:"authToken"`
}
