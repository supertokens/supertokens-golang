package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GenerateEmailVerifyToken(apiImplementation models.APIImplementation, options models.APIOptions) error {
	var result map[string]interface{}

	if apiImplementation.GenerateEmailVerifyTokenPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	response, err := apiImplementation.GenerateEmailVerifyTokenPOST(options)
	if err != nil {
		return err
	}

	if response.OK {
		result = map[string]interface{}{
			"status": "OK",
		}
	} else {
		result = map[string]interface{}{
			"status": "EMAIL_ALREADY_VERIFIED_ERROR",
		}
	}

	supertokens.Send200Response(options.Res, result)
	return nil
}
