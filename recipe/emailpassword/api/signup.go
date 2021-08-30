package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignUpAPI(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.SignUpPOST == nil {
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

	formFields, err := validateFormFieldsOrThrowError(options.Config.SignUpFeature.FormFields, formFieldsRaw["formFields"].([]interface{}))
	if err != nil {
		return err
	}

	result, err := apiImplementation.SignUpPOST(formFields, options)
	if err != nil {
		return err
	}
	if result.OK != nil {
		supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
			"user":   result.OK.User,
		})
	} else {
		return errors.FieldError{
			Msg: "Error in input formFields",
			Payload: []errors.ErrorPayload{{
				ID:    "email",
				Error: "This email already exists. Please sign in instead.",
			}},
		}
	}
	return nil
}
