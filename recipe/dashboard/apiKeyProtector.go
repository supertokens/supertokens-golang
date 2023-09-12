package dashboard

import (
	"errors"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	errors2 "github.com/supertokens/supertokens-golang/recipe/dashboard/errors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func apiKeyProtector(apiImpl dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext, call func() (interface{}, error)) error {
	shouldAllowAccess, err := (*options.RecipeImplementation.ShouldAllowAccess)(options.Req, options.Config, userContext)
	if err != nil {
		if errors.As(err, &errors2.ForbiddenAccessError{}) {
			return supertokens.SendNon200Response(options.Res, 403, map[string]interface{}{
				"message": err.Error(),
			})
		}

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
