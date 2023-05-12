package utilities

import (
	"bytes"
	"config/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type AuthClient struct {
	authServiceEndpoint string
	authServiceUserId   string
	authServicePassword string
}

func NewAuthClient(authServiceEndpoint string, authServiceUserId string, authServicePassword string) *AuthClient {
	return &AuthClient{authServiceEndpoint, authServiceUserId, authServicePassword}
}

func (ac *AuthClient) getAuthToken() (string, error) {
	request := struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}{ac.authServiceUserId, ac.authServicePassword}
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(ac.authServiceEndpoint+"/public/logins", "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(respBytes) == 0 {
		return "", errors.New("auth service returned zero-length token")
	}
	var result *models.User
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return "", err
	}
	return result.AuthToken, nil
}

func (ac *AuthClient) IsAuthorized(subjectToken string, permission string, applicationToken string) (bool, error) {
	requestStruct := struct {
		BearerToken string `json:"bearerToken"`
		Permission  string `json:"permission"`
	}{subjectToken, permission}
	requestJSON, err := json.Marshal(requestStruct)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest(http.MethodPost, ac.authServiceEndpoint+"/secure/authz-checks", bytes.NewBuffer(requestJSON))
	if err != nil {
		return false, err
	}
	client := &http.Client{}
	//Not sure if Go still requires this
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for key, val := range via[0].Header {
			req.Header[key] = val
		}
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", applicationToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var result *bool
	json.Unmarshal(respBytes, &result)
	if result == nil {
		return false, errors.New("unable to marshal response to bool")
	}
	return *result, nil
}
