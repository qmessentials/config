package utilities

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
	isAuthorized, err := ph.authClient.IsAuthorized(bearerToken, permission)
	if err != nil {
		return false, err
	}
	return isAuthorized, nil
}
