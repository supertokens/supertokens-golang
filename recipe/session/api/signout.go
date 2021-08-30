package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignOutAPI(apiImplementation models.APIInterface, options models.APIOptions) error {
	if apiImplementation.SignOutPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}
	_, err := apiImplementation.SignOutPOST(options)
	if err != nil {
		return err
	}

	supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
	})
	return nil
}
