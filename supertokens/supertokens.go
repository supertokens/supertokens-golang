package supertokens

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/errors"
)

type SuperTokens struct {
	AppInfo       NormalisedAppinfo
	RecipeModules []RecipeModule
}

var superTokensInstance *SuperTokens = nil

func SupertokensInit(config TypeInput) error {
	if superTokensInstance != nil {
		return nil
	}

	superTokensInstance := SuperTokens{}

	var err error
	superTokensInstance.AppInfo, err = NormaliseInputAppInfoOrThrowError(config.AppInfo)
	if err != nil {
		return err
	}

	if config.Supertoken != nil {
		hostList := strings.Split(config.Supertoken.ConnectionURI, ";")
		var hosts []NormalisedURLDomain
		for _, h := range hostList {
			host, err := NewNormalisedURLDomain(h, false)
			if err != nil {
				return err
			}
			hosts = append(hosts, *host)
		}

		InitQuerier(hosts, config.Supertoken.APIKey)
	} else {
		InitQuerier(nil, nil)
	}

	if config.RecipeList == nil || len(config.RecipeList) == 0 {
		return errors.BadInputError{Msg: "Please provide at least one recipe to the supertokens.init function call"}
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
		// err := s.SendTelemetry()
		// if err != nil {
		// 	return nil, err
		// }
	}

	return nil
}

func GetInstanceOrThrowError() (*SuperTokens, error) {
	if superTokensInstance != nil {
		return superTokensInstance, nil
	}
	return nil, errors.BadInputError{Msg: "Initialisation not done. Did you forget to call the SuperTokens.init function?"}
}

func (s *SuperTokens) SendTelemetry() {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return
	}

	response, err := querier.SendGetRequest(NormalisedURLPath{value: "/telemetry"}, map[string]interface{}{})
	if err != nil {
		return
	}
	var telemetryID string
	exists := response["exists"].(bool)
	if exists {
		telemetryID = response["telemetryId"].(string)
	}

	url := "https://api.supertokens.io/0/st/telemetry"

	// TODO: Add SDK name
	data := map[string]interface{}{
		"appName":       s.AppInfo.AppName,
		"websiteDomain": s.AppInfo.WebsiteDomain.GetAsStringDangerous(),
		"telemetryId":   telemetryID,
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
	_, err = client.Do(req)
	if err != nil {
		return
	}
}

func (s *SuperTokens) Middleware(theirHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, _ := NewNormalisedURLPath(r.RemoteAddr) // TODO: error handle
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
				// TODO: handle error...
				return
			}

			if id == nil {
				theirHandler.ServeHTTP(w, r)
				return
			}
			apiErr := matchedRecipe.HandleAPIRequest(*id, r, w, theirHandler, path, method)
			if apiErr != nil {
				// TODO: handle error
				return
			}
		} else {
			for _, recipeModule := range s.RecipeModules {
				id, err := recipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
				if err != nil {
					// TODO: handle error...
					return
				}

				if id != nil {
					err := recipeModule.HandleAPIRequest(*id, r, w, theirHandler, path, method)
					if err != nil {
						// TODO: handle error
						return
					}
					return
				}
			}
			theirHandler.ServeHTTP(w, r)
		}
	})
}

func (s *SuperTokens) getAllCORSHeaders() []string {
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
