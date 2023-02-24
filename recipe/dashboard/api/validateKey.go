package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/validationUtils"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ValidateKey(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions) error {
	isKeyValid, err := validationUtils.ValidateApiKey(options.Req, options.Config, supertokens.MakeDefaultUserContextFromAPI(options.Req))

	if err != nil {
		return err
	}

	if isKeyValid {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
		})
	} else {
		return supertokens.SendUnauthorisedAccess(options.Res)
	}
}
