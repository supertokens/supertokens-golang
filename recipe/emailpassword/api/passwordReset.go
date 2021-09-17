package api

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func PasswordReset(apiImplementation epmodels.APIInterface, options epmodels.APIOptions) error {
	if apiImplementation.PasswordResetPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := ioutil.ReadAll(options.Req.Body)
	if err != nil {
		return err
	}
	var formFieldsRaw map[string]interface{}
	err = json.Unmarshal(body, &formFieldsRaw)
	if err != nil {
		return err
	}

	formFields, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm, formFieldsRaw["formFields"].([]interface{}))
	if err != nil {
		return err
	}

	token, ok := formFieldsRaw["token"]
	if !ok {
		return supertokens.BadInputError{Msg: "Please provide the password reset token"}
	}
	if reflect.TypeOf(token).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "The password reset token must be a string"}
	}

	result, err := apiImplementation.PasswordResetPOST(formFields, token.(string), options)
	if err != nil {
		return err
	}
	if result.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
		})
	} else {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "RESET_PASSWORD_INVALID_TOKEN_ERROR",
		})
	}
}
