package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailverification/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GenerateEmailVerifyToken(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.GenerateEmailVerifyTokenPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	response, err := apiImplementation.GenerateEmailVerifyTokenPOST(options)
	if err != nil {
		return err
	}
	supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": response.Status,
	})
	return nil
}
