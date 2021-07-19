package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type bodyParams struct {
	ThirdPartyId string `json:"thirdPartyId"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirectURI"`
}

func SignInUpAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.SignInUpPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := ioutil.ReadAll(options.Req.Body)
	if err != nil {
		return err
	}
	var bodyParams bodyParams
	err = json.Unmarshal(body, &bodyParams)
	if err != nil {
		return err
	}

	if bodyParams.ThirdPartyId == "" {
		return supertokens.BadInputError{Msg: "Please provide the thirdPartyId in request body"}
	}

	if bodyParams.Code == "" {
		return supertokens.BadInputError{Msg: "Please provide the code in request body"}
	}

	if bodyParams.ThirdPartyId == "" {
		return supertokens.BadInputError{Msg: "Please provide the redirectURI in request body"}
	}

	var provider models.TypeProvider
	for _, prov := range options.Providers {
		if prov.ID == bodyParams.ThirdPartyId {
			provider = prov
		}
	}

	result := apiImplementation.SignInUpPOST(provider, bodyParams.Code, bodyParams.RedirectURI, options)

	if result.Status == "OK" {
		supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":         result.Status,
			"user":           result.User,
			"createdNewUser": result.CreatedNewUser,
		})
	} else if result.Status == "NO_EMAIL_GIVEN_BY_PROVIDER" {
		supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "NO_EMAIL_GIVEN_BY_PROVIDER",
		})
	} else {
		supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "FIELD_ERROR",
			"error":  result.Error,
		})
	}
	return nil
}
