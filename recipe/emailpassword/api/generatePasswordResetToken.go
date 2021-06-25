package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GeneratePasswordResetToken(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.GeneratePasswordResetTokenPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}
	formFields, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm, nil) // TODO: need help
	if err != nil {
		return err
	}
	result := apiImplementation.GeneratePasswordResetTokenPOST(formFields, options)
	supertokens.Send200Response(options.Res, result)
	return nil
}
