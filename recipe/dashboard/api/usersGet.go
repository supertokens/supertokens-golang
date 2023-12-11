package api

import (
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type UsersGetResponse struct {
	Status              string                     `json:"status"`
	NextPaginationToken *string                    `json:"nextPaginationToken,omitempty"`
	Users               []UserWithFirstAndLastName `json:"users"`
}

type UserWithFirstAndLastName struct {
	supertokens.User
	firstName string
	lastName  string
}

type UserPaginationResultWithFirstAndLastName struct {
	Users               []UserWithFirstAndLastName
	NextPaginationToken *string
}

func UsersGet(apiImplementation dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (UsersGetResponse, error) {
	req := options.Req
	limitStr := req.URL.Query().Get("limit")

	if limitStr == "" {
		return UsersGetResponse{}, supertokens.BadInputError{
			Msg: "Missing required parameter 'limit'",
		}
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return UsersGetResponse{}, err
	}

	timeJoinedOrder := req.URL.Query().Get("timeJoinedOrder")
	if timeJoinedOrder == "" {
		timeJoinedOrder = "DESC"
	}

	if timeJoinedOrder != "ASC" && timeJoinedOrder != "DESC" {
		return UsersGetResponse{}, supertokens.BadInputError{
			Msg: "Invalid value recieved for 'timeJoinedOrder'",
		}
	}

	paginationToken := req.URL.Query().Get("paginationToken")
	var paginationTokenPtr *string

	if paginationToken != "" {
		paginationTokenPtr = &paginationToken
	}

	var usersResponse supertokens.UserPaginationResult

	u, err := url.Parse(req.URL.String())
	if err != nil {
		return UsersGetResponse{}, err
	}

	queryParams, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return UsersGetResponse{}, err
	}

	queryParamsObject := map[string]string{}
	for i, s := range queryParams {
		queryParamsObject[i] = strings.Join(s, ";")
	}

	if len(queryParamsObject) != 0 {
		// the oder here doesn't matter cause in search, we return all users anyway.
		usersResponse, err = supertokens.GetUsersNewestFirst(tenantId, paginationTokenPtr, &limit, nil, queryParamsObject)
	} else if timeJoinedOrder == "ASC" {
		usersResponse, err = supertokens.GetUsersOldestFirst(tenantId, paginationTokenPtr, &limit, nil, nil)
	} else {
		usersResponse, err = supertokens.GetUsersNewestFirst(tenantId, paginationTokenPtr, &limit, nil, nil)
	}
	if err != nil {
		return UsersGetResponse{}, err
	}

	var userResponseWithFirstAndLastName UserPaginationResultWithFirstAndLastName = UserPaginationResultWithFirstAndLastName{}

	// copy userResponse into userResponseWithFirstAndLastName
	userResponseWithFirstAndLastName.NextPaginationToken = usersResponse.NextPaginationToken
	for _, userObj := range usersResponse.Users {
		userResponseWithFirstAndLastName.Users = append(userResponseWithFirstAndLastName.Users, struct {
			supertokens.User
			firstName string
			lastName  string
		}{
			User:      userObj,
			firstName: "",
			lastName:  "",
		})
	}

	_, err = usermetadata.GetRecipeInstanceOrThrowError()
	if err != nil {
		return UsersGetResponse{
			Status:              "OK",
			NextPaginationToken: usersResponse.NextPaginationToken,
			Users:               userResponseWithFirstAndLastName.Users,
		}, nil
	}

	var processingGroup sync.WaitGroup
	processingGroup.Add(len(usersResponse.Users))

	batchSize := 5
	var sem = make(chan int, batchSize)
	var errInBackground error

	for i, userObj := range userResponseWithFirstAndLastName.Users {
		sem <- 1

		if errInBackground != nil {
			return UsersGetResponse{}, errInBackground
		}

		go func(i int, userObj struct {
			supertokens.User
			firstName string
			lastName  string
		}) {
			defer processingGroup.Done()
			userMetadataResponse, err := usermetadata.GetUserMetadata(userObj.ID, userContext)
			<-sem
			if err != nil {
				errInBackground = err
				return
			}
			firstName, ok := userMetadataResponse["first_name"]
			lastName, ok2 := userMetadataResponse["last_name"]
			if ok {
				userResponseWithFirstAndLastName.Users[i].firstName = firstName.(string)
			}
			if ok2 {
				userResponseWithFirstAndLastName.Users[i].lastName = lastName.(string)
			}
		}(i, userObj)
	}

	if errInBackground != nil {
		return UsersGetResponse{}, errInBackground
	}

	processingGroup.Wait()

	return UsersGetResponse{
		Status:              "OK",
		NextPaginationToken: userResponseWithFirstAndLastName.NextPaginationToken,
		Users:               userResponseWithFirstAndLastName.Users,
	}, nil
}
