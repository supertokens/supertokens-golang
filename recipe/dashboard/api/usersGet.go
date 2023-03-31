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
	Status              string  `json:"status"`
	NextPaginationToken *string `json:"nextPaginationToken,omitempty"`
	Users               []Users `json:"users"`
}

type Users struct {
	RecipeId string `json:"recipeId"`
	User     User   `json:"user"`
}

type User struct {
	Id          string                     `json:"id"`
	TimeJoined  float64                    `json:"timeJoined"`
	FirstName   string                     `json:"firstName,omitempty"`
	LastName    string                     `json:"lastName,omitempty"`
	Email       string                     `json:"email,omitempty"`
	PhoneNumber string                     `json:"phoneNumber,omitempty"`
	ThirdParty  dashboardmodels.ThirdParty `json:"thirdParty,omitempty"`
}

func UsersGet(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions) (UsersGetResponse, error) {
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
		usersResponse, err = supertokens.GetUsersWithSearchParams(timeJoinedOrder, paginationTokenPtr, &limit, nil, queryParamsObject)
	} else if timeJoinedOrder == "ASC" {
		usersResponse, err = supertokens.GetUsersOldestFirst(paginationTokenPtr, &limit, nil)
	} else {
		usersResponse, err = supertokens.GetUsersNewestFirst(paginationTokenPtr, &limit, nil)
	}
	if err != nil {
		return UsersGetResponse{}, err
	}

	_, err = usermetadata.GetRecipeInstanceOrThrowError()
	if err != nil {
		return UsersGetResponse{
			Status:              "OK",
			NextPaginationToken: usersResponse.NextPaginationToken,
			Users:               getUsersTypeFromPaginationResult(usersResponse),
		}, nil
	}

	var processingGroup sync.WaitGroup
	processingGroup.Add(len(usersResponse.Users))

	batchSize := 5
	var sem = make(chan int, batchSize)
	var errInBackground error

	for i, userObj := range usersResponse.Users {
		sem <- 1

		if errInBackground != nil {
			return UsersGetResponse{}, errInBackground
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
		return UsersGetResponse{}, errInBackground
	}

	processingGroup.Wait()

	return UsersGetResponse{
		Status:              "OK",
		NextPaginationToken: usersResponse.NextPaginationToken,
		Users:               getUsersTypeFromPaginationResult(usersResponse),
	}, nil
}

func getUsersTypeFromPaginationResult(usersResponse supertokens.UserPaginationResult) []Users {
	users := []Users{}
	for _, v := range usersResponse.Users {
		user := User{
			Id:         v.User["id"].(string),
			TimeJoined: v.User["timeJoined"].(float64),
		}
		firstName := v.User["firstName"]
		if firstName != nil {
			user.FirstName = firstName.(string)
		}
		lastName := v.User["lastName"]
		if lastName != nil {
			user.LastName = lastName.(string)
		}

		if v.RecipeId == "emailpassword" {
			user.Email = v.User["email"].(string)
		} else if v.RecipeId == "thirdparty" {
			user.Email = v.User["email"].(string)
			user.ThirdParty = dashboardmodels.ThirdParty{
				Id:     v.User["thirdParty"].(map[string]interface{})["id"].(string),
				UserId: v.User["thirdParty"].(map[string]interface{})["userId"].(string),
			}
		} else {
			email := v.User["email"]
			if email != nil {
				user.Email = email.(string)
			}
			phoneNumber := v.User["phoneNumber"]
			if phoneNumber != nil {
				user.PhoneNumber = phoneNumber.(string)
			}
		}

		users = append(users, Users{
			RecipeId: v.RecipeId,
			User:     user,
		})
	}
	return users
}
