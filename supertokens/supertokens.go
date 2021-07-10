package supertokens

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type SuperTokens struct {
	AppInfo        NormalisedAppinfo
	RecipeModules  []RecipeModule
	OnGeneralError func(err error, req *http.Request, res http.ResponseWriter)
}

var superTokensInstance *SuperTokens

func SupertokensInit(config TypeInput) error {
	if superTokensInstance != nil {
		return nil
	}
	superTokensInstance := &SuperTokens{}

	superTokensInstance.OnGeneralError = onGeneralError
	if config.OnGeneralError != nil {
		superTokensInstance.OnGeneralError = config.OnGeneralError
	}

	var err error
	superTokensInstance.AppInfo, err = NormaliseInputAppInfoOrThrowError(config.AppInfo)
	if err != nil {
		return err
	}

	if config.Supertokens != nil {
		hostList := strings.Split(config.Supertokens.ConnectionURI, ";")
		var hosts []NormalisedURLDomain
		for _, h := range hostList {
			host, err := NewNormalisedURLDomain(h, false)
			if err != nil {
				return err
			}
			hosts = append(hosts, *host)
		}

		initQuerier(hosts, config.Supertokens.APIKey)
	} else {
		initQuerier(nil, nil)
	}

	if config.RecipeList == nil || len(config.RecipeList) == 0 {
		return errors.New("Please provide at least one recipe to the supertokens.init function call")
	}

	for _, elem := range config.RecipeList {
		recipeModule, err := elem(superTokensInstance.AppInfo)
		if err != nil {
			return err
		}
		superTokensInstance.RecipeModules = append(superTokensInstance.RecipeModules, *recipeModule)
	}

	if config.Telemetry != nil && *config.Telemetry {
		// TODO: Telemetry
		// err := SendTelemetry()
		// if err != nil {
		// 	return err
		// }
	}

	return nil
}

func getInstanceOrThrowError() (*SuperTokens, error) {
	if superTokensInstance != nil {
		return superTokensInstance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the SuperTokens.init function?")
}

func SendTelemetry() error {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return err
	}

	response, err := querier.SendGetRequest(NormalisedURLPath{value: "/telemetry"}, nil)
	if err != nil {
		return err
	}
	var telemetryID string
	exists := response["exists"].(bool)
	if exists {
		telemetryID = response["telemetryId"].(string)
	}

	url := "https://api.supertokens.io/0/st/telemetry"

	// TODO: Add SDK name
	data := map[string]interface{}{
		"appName":       superTokensInstance.AppInfo.AppName,
		"websiteDomain": superTokensInstance.AppInfo.WebsiteDomain.GetAsStringDangerous(),
		"telemetryId":   telemetryID,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("api-version", "2")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func (s *SuperTokens) Middleware(theirHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, err := NewNormalisedURLPath(r.RemoteAddr)
		if err != nil {
			s.ErrorHandler(err, r, w)
		}
		path := s.AppInfo.APIGatewayPath.AppendPath(*reqURL)
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
				s.ErrorHandler(err, r, w)
				return
			}

			if id == nil {
				theirHandler.ServeHTTP(w, r)
				return
			}
			apiErr := matchedRecipe.HandleAPIRequest(*id, r, w, theirHandler, path, method)
			if apiErr != nil {
				s.ErrorHandler(err, r, w)
				return
			}
		} else {
			for _, recipeModule := range s.RecipeModules {
				id, err := recipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
				if err != nil {
					s.ErrorHandler(err, r, w)
					return
				}

				if id != nil {
					err := recipeModule.HandleAPIRequest(*id, r, w, theirHandler, path, method)
					if err != nil {
						s.ErrorHandler(err, r, w)
						return
					}
					return
				}
			}
			theirHandler.ServeHTTP(w, r)
		}
	})
}

func (s *SuperTokens) GetAllCORSHeaders() []string {
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

func (s *SuperTokens) ErrorHandler(err error, req *http.Request, res http.ResponseWriter) bool {
	if errors.As(err, &BadInputError{}) {
		if catcherr := SendNon200Response(res, err.Error(), 400); catcherr != nil {
			panic("internal server err" + catcherr.Error())
		}
		return true
	}
	for _, recipe := range s.RecipeModules {
		handled := recipe.HandleError(err, req, res)
		if handled {
			return true
		}
	}
	superTokensInstance.OnGeneralError(err, req, res)
	return true
}

func onGeneralError(err error, req *http.Request, res http.ResponseWriter) {
	SendNon200Response(res, err.Error(), 500)
}
