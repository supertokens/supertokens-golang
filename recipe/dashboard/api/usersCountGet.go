package api

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type usersCountGetResponse struct {
	Status string  `json:"status"`
	Count  float64 `json:"count"`
}

func UsersCountGet(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (usersCountGetResponse, error) {
	count, err := supertokens.GetUserCount(nil)
	if err != nil {
		return usersCountGetResponse{}, err
	}

	return usersCountGetResponse{
		Status: "OK",
		Count:  count,
	}, nil
}
