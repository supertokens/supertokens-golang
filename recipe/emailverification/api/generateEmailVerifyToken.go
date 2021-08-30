package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GenerateEmailVerifyToken(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.GenerateEmailVerifyTokenPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	response, err := apiImplementation.GenerateEmailVerifyTokenPOST(options)
	if err != nil {
		return err
	}
	if response.EmailAlreadyVerifiedError != nil {
		supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "EMAIL_ALREADY_VERIFIED_ERROR",
		})
	} else {
		supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
		})
	}
	return nil
}
