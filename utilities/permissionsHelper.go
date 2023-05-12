package utilities

import (
	"bytes"
	"config/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

type PermissionsHelper struct {
	authServiceEndpoint string
	authServiceUserId   string
	authServicePassword string
	redisClient         *redis.Client
}

func NewPermissionHelper(authServiceEndpoint string, authServiceUserId string, authServicePassword string, redisClient *redis.Client) *PermissionsHelper {
	return &PermissionsHelper{
		authServiceEndpoint,
		authServiceUserId,
		authServicePassword,
		redisClient,
	}
}

func (ph *PermissionsHelper) getAuthToken() (string, error) {
	authTokenCacheResult, err := ph.redisClient.Get("authToken").Result()
	if err == redis.Nil || len(authTokenCacheResult) == 0 {
		log.Warn().Msg("Cache miss for auth token")
		request := struct {
			UserId   string `json:"userId"`
			Password string `json:"password"`
		}{ph.authServiceUserId, ph.authServicePassword}
		requestJSON, err := json.Marshal(request)
		if err != nil {
			return "", err
		}
		resp, err := http.Post(ph.authServiceEndpoint+"/public/logins", "application/json", bytes.NewBuffer(requestJSON))
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
		ph.redisClient.Set("authToken", result.AuthToken, 0)
		return result.AuthToken, nil
	}
	if err != nil {
		return "", err
	}
	return authTokenCacheResult, nil
}

func (ph *PermissionsHelper) IsAuthorized(bearerToken string, permission string) (bool, error) {
	authToken, err := ph.getAuthToken()
	if err != nil {
		return false, err
	}
	if len(authToken) == 0 {
		return false, errors.New("auth token returned with length zero")
	}
	requestStruct := struct {
		BearerToken string `json:"bearerToken"`
		Permission  string `json:"permission"`
	}{bearerToken, permission}
	requestJSON, err := json.Marshal(requestStruct)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest(http.MethodPost, ph.authServiceEndpoint+"/secure/authz-checks", bytes.NewBuffer(requestJSON))
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	// dump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	return false, err
	// }
	// log.Info().Msg(string(dump))
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
