package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func UsersCountGet(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions) error {
	count, err := supertokens.GetUserCount(nil)
	if err != nil {
		return err
	}
	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
		"count":  count,
	})
}
