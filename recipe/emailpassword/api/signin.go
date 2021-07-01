package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignInAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
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

	formFields, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm, formFieldsRaw["formFields"].([]models.FormFieldValue))
	if err != nil {
		return err
	}
	result := apiImplementation.SignInPOST(formFields, options)
	supertokens.Send200Response(options.Res, result)
	return nil
}
