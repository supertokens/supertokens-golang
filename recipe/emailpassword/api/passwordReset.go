package api

import "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"

func PasswordReset(apiImplementation models.APIImplementation, options models.APIOptions)error {
	if apiImplementation.PasswordResetPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}
	_, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm, options.Req) // TODO: need help
	if err != nil {
		return err
	}
	// TODO: 
	// token := 
	// result := apiImplementation.GeneratePasswordResetTokenPOST(formFields, options)
	// supertokens.Send200Response(options.Res, result)
	return nil

}
