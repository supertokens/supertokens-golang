package api

import (
	"strconv"
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func UsersGet(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions) error {
	req := options.Req
	limitStr := req.URL.Query().Get("limit")

	if limitStr == "" {
		return supertokens.BadInputError{
			Msg: "Missing required parameter 'limit'",
		}
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return err
	}

	timeJoinedOrder := req.URL.Query().Get("timeJoinedOrder")
	if timeJoinedOrder == "" {
		timeJoinedOrder = "DESC"
	}

	if timeJoinedOrder != "ASC" && timeJoinedOrder != "DESC" {
		return supertokens.BadInputError{
			Msg: "Invalid value recieved for 'timeJoinedOrder'",
		}
	}

	paginationToken := req.URL.Query().Get("paginationToken")
	var paginationTokenPtr *string

	if paginationToken != "" {
		paginationTokenPtr = &paginationToken
	}

	var usersResponse supertokens.UserPaginationResult

	if timeJoinedOrder == "ASC" {
		usersResponse, err = supertokens.GetUsersOldestFirst(paginationTokenPtr, &limit, nil)
	} else {
		usersResponse, err = supertokens.GetUsersNewestFirst(paginationTokenPtr, &limit, nil)
	}
	if err != nil {
		return err
	}

	_, err = usermetadata.GetRecipeInstanceOrThrowError()
	if err != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":              "OK",
			"users":               usersResponse.Users,
			"nextPaginationToken": usersResponse.NextPaginationToken,
		})
	}

	var processingGroup sync.WaitGroup
	processingGroup.Add(len(usersResponse.Users))

	batchSize := 5
	var sem = make(chan int, batchSize)
	var errInBackground error

	for i, userObj := range usersResponse.Users {
		sem <- 1

		if errInBackground != nil {
			return errInBackground
		}

		go func(i int, userObj struct {
			RecipeId string                 `json:"recipeId"`
			User     map[string]interface{} `json:"user"`
		}) {
			defer processingGroup.Done()
			userMetadataResponse, err := usermetadata.GetUserMetadata(userObj.User["id"].(string))
			<-sem
			if err != nil {
				errInBackground = err
				return
			}
			usersResponse.Users[i].User["firstName"] = userMetadataResponse["first_name"]
			usersResponse.Users[i].User["lastName"] = userMetadataResponse["last_name"]
		}(i, userObj)
	}

	if errInBackground != nil {
		return errInBackground
	}

	processingGroup.Wait()

	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status":              "OK",
		"users":               usersResponse.Users,
		"nextPaginationToken": usersResponse.NextPaginationToken,
	})
}
