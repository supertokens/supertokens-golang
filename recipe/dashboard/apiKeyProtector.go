package dashboard

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func apiKeyProtector(apiImpl dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext, call func() (interface{}, error)) error {
	shouldAllowAccess, err := (*options.RecipeImplementation.ShouldAllowAccess)(options.Req, options.Config, userContext)
	if err != nil {
		return err
	}

	if !shouldAllowAccess {
		return supertokens.SendUnauthorisedAccess(options.Res)
	}

	resp, err := call()
	if err != nil {
		return err
	}

	return supertokens.Send200Response(options.Res, resp)
}
