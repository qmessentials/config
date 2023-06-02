package routers

import (
	"config/utilities"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func checkPermissions(c *gin.Context, permission string, permissionsHelper *utilities.PermissionsHelper) int {
	authHeader := c.Request.Header.Get("Authorization")
	log.Info().Msgf("Auth header is %s", authHeader)
	bearerPattern := regexp.MustCompile("(?i)^bearer (.*)$")
	tokens := bearerPattern.FindStringSubmatch(authHeader)
	if len(tokens) != 2 {
		log.Warn().Msg("Unauthenticated attempt to retrieve config data")
		return http.StatusUnauthorized
	}
	isAllowed, err := permissionsHelper.IsAuthorized(tokens[1], permission)
	if err != nil {
		if err.Error() == "authentication failure" {
			return http.StatusUnauthorized
		}
		log.Error().Err(err).Msg("Error checking permissions")
		return http.StatusInternalServerError
	}
	if !isAllowed {
		log.Warn().Msgf("Failed permission check")
		return http.StatusForbidden
	}
	return http.StatusOK
}
