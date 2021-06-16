package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignOutAPI(apiImplementation models.APIImplementation, options models.APIOptions) error {
	if apiImplementation.SignOutPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}
	result, err := apiImplementation.SignOutPOST(options)
	if err != nil {
		return err
	}
	supertokens.Send200Response(options.Res, result)
	return nil
}
