package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignOutAPI(apiImplementation sessmodels.APIInterface, options sessmodels.APIOptions) error {
	if apiImplementation.SignOutPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}
	_, err := apiImplementation.SignOutPOST(options)
	if err != nil {
		return err
	}

	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
	})
}
