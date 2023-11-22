package api

import (
	"fmt"
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetMagicLink(appInfo supertokens.NormalisedAppinfo, recipeID string, preAuthSessionID string, linkCode string, tenantId string, request *http.Request, userContext supertokens.UserContext) (string, error) {
	websiteDomain, err := appInfo.GetOrigin(request, userContext)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s%s/verify?rid=%s&preAuthSessionId=%s&tenantId=%s#%s",
		websiteDomain.GetAsStringDangerous(),
		appInfo.WebsiteBasePath.GetAsStringDangerous(),
		recipeID,
		preAuthSessionID,
		tenantId,
		linkCode,
	), nil
}
