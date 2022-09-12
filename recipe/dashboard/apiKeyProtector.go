package dashboard

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func apiKeyProtector(apiImpl dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext, call func() error) error {
	shouldAllowAccess, err := options.RecipeImplementation.ShouldAllowAccess(options.Req, userContext)
	if err != nil {
		return err
	}

	if !shouldAllowAccess {
		return supertokens.SendUnauthorisedAccess(options.Res)
	}

	return call()
}
