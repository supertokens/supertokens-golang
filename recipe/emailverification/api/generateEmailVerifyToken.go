package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GenerateEmailVerifyToken(apiImplementation evmodels.APIInterface, options evmodels.APIOptions) error {
	if apiImplementation.GenerateEmailVerifyTokenPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	response, err := apiImplementation.GenerateEmailVerifyTokenPOST(options)
	if err != nil {
		return err
	}
	if response.EmailAlreadyVerifiedError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "EMAIL_ALREADY_VERIFIED_ERROR",
		})
	} else {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
		})
	}
}
