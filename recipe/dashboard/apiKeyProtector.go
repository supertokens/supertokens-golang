package dashboard

import (
	defaultErrors "errors"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/errors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func apiKeyProtector(apiImpl dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext, call func() (interface{}, error)) error {
	shouldAllowAccess, err := (*options.RecipeImplementation.ShouldAllowAccess)(options.Req, options.Config, userContext)
	if err != nil {
		if defaultErrors.As(err, &errors.OperationNotAllowedError{}) {
			body := map[string]string{
				"message": err.Error(),
			}
			return supertokens.Send403Response(options.Res, body)
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
