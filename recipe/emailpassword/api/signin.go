package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignInAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.SignInPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}
	formFields, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm, options.Req) // TODO: need help
	if err != nil {
		return err
	}
	result := apiImplementation.SignInPOST(formFields, options)
	supertokens.Send200Response(options.Res, result)
	return nil
}
