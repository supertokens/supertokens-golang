package api

import (
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func LoginMethodsAPI(apiImplementation multitenancymodels.APIInterface, options multitenancymodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.LoginMethodsGET == nil || (*apiImplementation.LoginMethodsGET) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	queryParams := options.Req.URL.Query()

	var clientType *string = nil
	if clientTypeStrFromQueryParams := queryParams.Get("clientType"); clientTypeStrFromQueryParams != "" {
		clientType = &clientTypeStrFromQueryParams
	}

	result, err := (*apiImplementation.LoginMethodsGET)("public", clientType, options, userContext)
	if err != nil {
		return err
	}

	if result.OK != nil {
		return supertokens.Send200Response(options.Res, result.OK)
	} else if result.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
