package supertokens

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"strings"
)

type superTokens struct {
	AppInfo        NormalisedAppinfo
	RecipeModules  []RecipeModule
	OnGeneralError func(err error, req *http.Request, res http.ResponseWriter)
}

var superTokensInstance *superTokens

func supertokensInit(config TypeInput) error {
	if superTokensInstance != nil {
		return nil
	}
	superTokens := &superTokens{}

	superTokens.OnGeneralError = defaultOnGeneralError
	if config.OnGeneralError != nil {
		superTokens.OnGeneralError = config.OnGeneralError
	}

	var err error
	superTokens.AppInfo, err = NormaliseInputAppInfoOrThrowError(config.AppInfo)
	if err != nil {
		return err
	}

	if config.Supertokens != nil {
		hostList := strings.Split(config.Supertokens.ConnectionURI, ";")
		var hosts []NormalisedURLDomain
		for _, h := range hostList {
			host, err := NewNormalisedURLDomain(h)
			if err != nil {
				return err
			}
			hosts = append(hosts, host)
		}

		initQuerier(hosts, config.Supertokens.APIKey)
	} else {
		// TODO: Add tests for init without supertokens core.
	}

	if config.RecipeList == nil || len(config.RecipeList) == 0 {
		return errors.New("please provide at least one recipe to the supertokens.init function call")
	}

	for _, elem := range config.RecipeList {
		recipeModule, err := elem(superTokens.AppInfo)
		if err != nil {
			return err
		}
		superTokens.RecipeModules = append(superTokens.RecipeModules, *recipeModule)
	}

	if config.Telemetry != nil && *config.Telemetry {
		sendTelemetry()
	}

	superTokensInstance = superTokens
	return nil
}

func defaultOnGeneralError(err error, req *http.Request, res http.ResponseWriter) {
	http.Error(res, err.Error(), 500)
}

func getInstanceOrThrowError() (*superTokens, error) {
	if superTokensInstance != nil {
		return superTokensInstance, nil
	}
	return nil, errors.New("initialisation not done. Did you forget to call the SuperTokens.init function?")
}

// TODO: add test to see query
func sendTelemetry() {
	if flag.Lookup("test.v") != nil {
		// if running in test mode, we do not want to send this.
		return
	}
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return
	}

	response, err := querier.SendGetRequest("/telemetry", nil)
	if err != nil {
		return
	}
	exists := response["exists"].(bool)

	url := "https://api.supertokens.io/0/st/telemetry"

	data := map[string]interface{}{
		"appName":       superTokensInstance.AppInfo.AppName,
		"websiteDomain": superTokensInstance.AppInfo.WebsiteDomain.GetAsStringDangerous(),
		"sdk":           "golang",
	}
	if exists {
		data["telemetryId"] = response["telemetryId"].(string)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("api-version", "2")

	client := &http.Client{}
	client.Do(req)
}

func (s *superTokens) middleware(theirHandler http.Handler) http.Handler {
	if theirHandler == nil {
		theirHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, err := NewNormalisedURLPath(r.URL.Path)
		if err != nil {
			err = s.errorHandler(err, r, w)
			if err != nil {
				s.OnGeneralError(err, r, w)
			}
			return
		}
		path := s.AppInfo.APIGatewayPath.AppendPath(reqURL)
		method := r.Method

		if !strings.HasPrefix(path.GetAsStringDangerous(), s.AppInfo.APIBasePath.GetAsStringDangerous()) {
			theirHandler.ServeHTTP(w, r)
			return
		}
		requestRID := getRIDFromRequest(r)
		if requestRID != "" {
			var matchedRecipe *RecipeModule
			for _, recipeModule := range s.RecipeModules {
				if recipeModule.GetRecipeID() == requestRID {
					matchedRecipe = &recipeModule
					break
				}
			}
			if matchedRecipe == nil {
				theirHandler.ServeHTTP(w, r)
				return
			}
			id, err := matchedRecipe.ReturnAPIIdIfCanHandleRequest(path, method)

			if err != nil {
				err = s.errorHandler(err, r, w)
				if err != nil {
					s.OnGeneralError(err, r, w)
				}
				return
			}

			if id == nil {
				theirHandler.ServeHTTP(w, r)
				return
			}
			apiErr := matchedRecipe.HandleAPIRequest(*id, r, w, theirHandler.ServeHTTP, path, method)
			if apiErr != nil {
				apiErr = s.errorHandler(apiErr, r, w)
				if apiErr != nil {
					s.OnGeneralError(apiErr, r, w)
				}
				return
			}
		} else {
			for _, recipeModule := range s.RecipeModules {
				id, err := recipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
				if err != nil {
					err = s.errorHandler(err, r, w)
					if err != nil {
						s.OnGeneralError(err, r, w)
					}
					return
				}

				if id != nil {
					err := recipeModule.HandleAPIRequest(*id, r, w, theirHandler.ServeHTTP, path, method)
					if err != nil {
						err = s.errorHandler(err, r, w)
						if err != nil {
							s.OnGeneralError(err, r, w)
						}
					}
					return
				}
			}
			theirHandler.ServeHTTP(w, r)
		}
	})
}

func (s *superTokens) getAllCORSHeaders() []string {
	headerMap := map[string]bool{HeaderRID: true, HeaderFDI: true}
	for _, recipe := range s.RecipeModules {
		headers := recipe.GetAllCORSHeaders()
		for _, header := range headers {
			headerMap[header] = true
		}
	}
	var headers []string
	for header := range headerMap {
		headers = append(headers, header)
	}
	return headers
}

func (s *superTokens) errorHandler(originalError error, req *http.Request, res http.ResponseWriter) error {
	// TODO: replace errors.As with errors.Is if we are not casting the error to that specific type.
	if errors.As(originalError, &BadInputError{}) {
		if catcher := SendNon200Response(res, originalError.Error(), 400); catcher != nil {
			s.OnGeneralError(originalError, req, res)
		}
		return nil
	}
	for _, recipe := range s.RecipeModules {
		if recipe.HandleError != nil {
			handled, err := recipe.HandleError(originalError, req, res)
			if err != nil {
				return err
			}
			if handled {
				return nil
			}
		}
	}
	return originalError
}

// TODO: make this an array of users.
type UserPaginationResult struct {
	Users struct {
		recipeId string
		user     map[string]interface{}
	}
	NextPaginationToken *string
}

// TODO: Add tests
func getUsers(timeJoinedOrder string, limit *int, paginationToken *string, includeRecipeIds *[]string) (*UserPaginationResult, error) {

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return nil, err
	}

	requestBody := map[string]interface{}{
		"timeJoinedOrder": timeJoinedOrder,
	}
	if limit != nil {
		requestBody["limit"] = *limit
	}
	if paginationToken != nil {
		requestBody["paginationToken"] = *paginationToken
	}
	if includeRecipeIds != nil {
		requestBody["includeRecipeIds"] = strings.Join((*includeRecipeIds)[:], ",")
	}

	resp, err := querier.SendGetRequest("/users", requestBody)

	if err != nil {
		return nil, err
	}

	// TODO: try not to do marshal and unmarshal
	// TODO: Also, Unmarshal is slow, so try and use something else.
	temporaryVariable, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	var result = UserPaginationResult{}

	err = json.Unmarshal(temporaryVariable, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// TODO: Add tests
func getUserCount(includeRecipeIds *[]string) (int, error) {

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return -1, err
	}

	requestBody := map[string]interface{}{}

	if includeRecipeIds != nil {
		requestBody["includeRecipeIds"] = strings.Join((*includeRecipeIds)[:], ",")
	}

	resp, err := querier.SendGetRequest("/users/count", requestBody)

	if err != nil {
		return -1, err
	}

	return resp["count"].(int), nil
}
