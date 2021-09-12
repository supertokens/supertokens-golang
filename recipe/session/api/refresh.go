package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func HandleRefreshAPI(apiImplementation sessmodels.APIInterface, options sessmodels.APIOptions) error {
	if apiImplementation.RefreshPOST == nil {
		options.OtherHandler.ServeHTTP(options.Res, options.Req)
		return nil
	}
	err := apiImplementation.RefreshPOST(options)
	if err != nil {
		return err
	}
	return supertokens.Send200Response(options.Res, nil)
}
