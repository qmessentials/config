package utilities

import "strings"

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
func (ph *PermissionsHelper) IsAuthorized(bearerToken string, permission string) (bool, error) {
	foundInCache, valuesPipeJoined, err := ph.cacheService.Get("PERMISSIONS|" + bearerToken)
	if foundInCache {
		values := strings.Split(valuesPipeJoined, "|")
		for _, value := range values {
			if value == permission {
				return true, nil
			}
		}
	}
	if err != nil {
		return false, err
	}
	isAuthorized, err := ph.authClient.IsAuthorized(bearerToken, permission)
	if err != nil {
		return false, err
	}
	if isAuthorized {
		if foundInCache {
			valuesPipeJoined += "|" + permission
		} else {
			valuesPipeJoined = permission
		}
		ph.cacheService.Set("PERMISSIONS|"+bearerToken, valuesPipeJoined)
	}
	return isAuthorized, nil
}
