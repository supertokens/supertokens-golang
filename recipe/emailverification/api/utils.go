package api

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetEmailVerifyLink(appInfo supertokens.NormalisedAppinfo, token string, recipeID string, tenantId string) string {
	return fmt.Sprintf(
		"%s%s/verify-email?token=%s&rid=%s&tenantId=%s",
		appInfo.WebsiteDomain.GetAsStringDangerous(),
		appInfo.WebsiteBasePath.GetAsStringDangerous(),
		token,
		recipeID,
		tenantId,
	)
}
