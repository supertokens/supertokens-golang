package api

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetMagicLink(appInfo supertokens.NormalisedAppinfo, recipeID string, preAuthSessionID string, linkCode string, tenantId string) string {
	return fmt.Sprintf(
		"%s%s/verify?rid=%s&preAuthSessionId=%s&tenantId=%s#%s",
		appInfo.WebsiteDomain.GetAsStringDangerous(),
		appInfo.WebsiteBasePath.GetAsStringDangerous(),
		recipeID,
		preAuthSessionID,
		tenantId,
		linkCode,
	)
}
