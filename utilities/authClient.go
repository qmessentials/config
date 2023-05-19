package utilities

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type AuthClient struct {
	authServiceEndpoint string
}

func NewAuthClient(authServiceEndpoint string) *AuthClient {
	return &AuthClient{authServiceEndpoint}
}

func (ac *AuthClient) IsAuthorized(subjectToken string, permission string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, ac.authServiceEndpoint+"/secure/authz-checks/"+permission, nil)
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", subjectToken))
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
