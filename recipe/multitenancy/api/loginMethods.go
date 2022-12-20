package api

import (
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func LoginMethodsAPI(apiImplementation multitenancymodels.APIInterface, options multitenancymodels.APIOptions) error {
	if apiImplementation.LoginMethodsGET == nil || (*apiImplementation.LoginMethodsGET) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	queryParams := options.Req.URL.Query()

	var tenantId *string = nil
	if tenantIdStrFromQueryParams := queryParams.Get("tenantId"); tenantIdStrFromQueryParams != "" {
		tenantId = &tenantIdStrFromQueryParams
	}

	userContext := supertokens.MakeDefaultUserContextFromAPI(options.Req)

	tenantId, err := (*options.RecipeImplementation.GetTenantId)(tenantId, userContext)
	if err != nil {
		return err
	}

	result, err := (*apiImplementation.LoginMethodsGET)(tenantId, options, userContext)
	if err != nil {
		return err
	}

	if result.OK != nil {
		return supertokens.Send200Response(options.Res, result.OK)
	} else if result.TenantDoesNotExistError != nil {
		return supertokens.SendNon200ResponseWithMessage(options.Res, "Tenant does not exist", 422)
	} else if result.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
