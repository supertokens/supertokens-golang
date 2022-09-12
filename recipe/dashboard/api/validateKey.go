package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateKey(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions) error {
	shouldAllowAccess, err := options.RecipeImplementation.ShouldAllowAccess(options.Req, supertokens.MakeDefaultUserContextFromAPI(options.Req))
	if err != nil {
		return err
	}

	if !shouldAllowAccess {
		return supertokens.SendUnauthorisedAccess(options.Res)
	}

	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
	})
}
