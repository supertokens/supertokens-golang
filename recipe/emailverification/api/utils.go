package api

import (
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetEmailVerifyLink(appInfo supertokens.NormalisedAppinfo, token string, recipeID string, tenantId string, request *http.Request, userContext supertokens.UserContext) (string, error) {
	websiteDomain, err := appInfo.GetOrigin(request, userContext)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s%s/verify-email?token=%s&rid=%s&tenantId=%s",
		websiteDomain.GetAsStringDangerous(),
		appInfo.WebsiteBasePath.GetAsStringDangerous(),
		token,
		recipeID,
		tenantId,
	), nil
}
