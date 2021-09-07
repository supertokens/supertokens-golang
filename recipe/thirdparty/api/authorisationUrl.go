package api

import (
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func AuthorisationUrlAPI(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.AuthorisationUrlGET == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	queryParams := options.Req.URL.Query()
	thirdPartyId := queryParams.Get("thirdPartyId")

	if len(thirdPartyId) == 0 {
		return supertokens.BadInputError{Msg: "Please provide the thirdPartyId as a GET param"}
	}

	var provider models.TypeProvider
	for _, prov := range options.Providers {
		if prov.ID == thirdPartyId {
			provider = prov
		}
	}
	if reflect.DeepEqual(provider, models.TypeProvider{}) {
		return supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to not be configured on the backend. Please check your frontend and backend configs."}
	}

	result, err := apiImplementation.AuthorisationUrlGET(provider, options)
	if err != nil {
		return err
	}
	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
		"url":    result.OK.Url,
	})
}
