package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignUpAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.SignUpPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}
	formFields, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm, options.Req) // TODO: need help
	if err != nil {
		return err
	}
	result := apiImplementation.SignUpPOST(formFields, options)
	if result.Ok != nil {
		supertokens.Send200Response(options.Res, result)
		return nil
	}
	return errors.FieldError{
		Msg: "Error in input formFields",
		Payload: []errors.ErrorPayload{{
			ID:    "email",
			Error: "This email already exists. Please sign in instead.",
		}},
	}
}
