package utilities

import (
	"errors"

	"github.com/rs/zerolog/log"
)

type PermissionsHelper struct {
	authClient   *AuthClient
	cacheService *CacheService
}

func NewPermissionHelper(authClient *AuthClient, cacheService *CacheService) *PermissionsHelper {
	return &PermissionsHelper{
		authClient,
		cacheService,
	}
}

func (ph *PermissionsHelper) getAuthToken() (string, error) {
	cacheHit, authToken, err := ph.cacheService.Get("authToken")
	if err != nil {
		return "", err
	}
	if cacheHit {
		return authToken, nil
	}
	log.Warn().Msg("Cache miss for auth token")
	authToken, err = ph.authClient.getAuthToken()
	if err != nil {
		return "", err
	}
	ph.cacheService.Set("authToken", authToken)
	return authToken, nil
}

func (ph *PermissionsHelper) IsAuthorized(bearerToken string, permission string) (bool, error) {
	applicationToken, err := ph.getAuthToken()
	if err != nil {
		return false, err
	}
	if len(applicationToken) == 0 {
		return false, errors.New("auth token returned with length zero")
	}
	return ph.authClient.IsAuthorized(bearerToken, permission, applicationToken)
}
