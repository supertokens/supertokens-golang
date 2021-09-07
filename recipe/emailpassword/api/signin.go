package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignInAPI(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.SignInPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := ioutil.ReadAll(options.Req.Body)
	if err != nil {
		panic(err)
	}
	var formFieldsRaw map[string]interface{}
	err = json.Unmarshal(body, &formFieldsRaw)
	if err != nil {
		panic(err)
	}

	formFields, err := validateFormFieldsOrThrowError(options.Config.SignInFeature.FormFields, formFieldsRaw["formFields"].([]interface{}))
	if err != nil {
		return err
	}

	result, err := apiImplementation.SignInPOST(formFields, options)
	if err != nil {
		return err
	}
	if result.WrongCredentialsError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "WRONG_CREDENTIALS_ERROR",
		})
	} else {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
			"user":   result.OK.User,
		})
	}
}
