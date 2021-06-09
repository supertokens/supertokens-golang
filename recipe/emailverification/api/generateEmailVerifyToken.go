package api

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/schema"
)

func GenerateEmailVerifyToken(apiImplementation schema.APIImplementation, options schema.APIOptions) error {
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

	// TODO: send200Response(options.res, result);
	fmt.Printf("", result)
	return nil
}
