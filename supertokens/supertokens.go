package supertokens

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type SuperTokens struct {
	AppInfo           NormalisedAppinfo
	IsInServerlessEnv bool
	RecipeModules     []RecipeModule
}

func NewSuperTokens(config TypeInput) (*SuperTokens, error) {
	var err error
	var s *SuperTokens
	s.AppInfo, err = NormaliseInputAppInfoOrThrowError(config.AppInfo)
	if err != nil {
		return nil, err
	}

	hostList := strings.Split(config.Supertoken.ConnectionURI, ";")
	var hosts []NormalisedURLDomain
	for _, h := range hostList {
		host, err := NewNormalisedURLDomain(h, false)
		if err != nil {
			fmt.Println(err)
		}
		hosts = append(hosts, *host)
	}

	q := Querier{}
	q.Init(hosts, config.Supertoken.APIKey)

	if config.RecipeList == nil || len(config.RecipeList) == 0 {
		return nil, errors.New("Please provide at least one recipe to the supertokens.init function call")
	}

	s.IsInServerlessEnv = config.IsInServerlessEnv

	for _, elem := range config.RecipeList {
		s.RecipeModules = append(s.RecipeModules, elem(s.AppInfo, s.IsInServerlessEnv))
	}

	telemetry := config.Telemetry

	if telemetry {
		err := s.SendTelemetry()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *SuperTokens) SendTelemetry() error {
	q := &Querier{}
	querier, err := q.GetNewInstanceOrThrowError(s.IsInServerlessEnv, "")
	if err != nil {
		return err
	}

	response, err := querier.SendGetRequest(NormalisedURLPath{Value: "/telemetry"}, map[string]string{})
	if err != nil {
		return err
	}
	var telemetryID string
	if response["exists"] == true {
		telemetryID = response["telemetryId"].(string)
	}

	url := "https://api.supertokens.io/0/st/telemetry"

	data := map[string]interface{}{
		"appName":       s.AppInfo.AppName,
		"websiteDomain": s.AppInfo.WebsiteDomain.Value,
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

func (s *SuperTokens) GetInstanceOrThrowError() *SuperTokens {
	return s
}

func (s *SuperTokens) middleware() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, _ := NewNormalisedURLPath(r.URL.Path)        //error handle
		path := s.AppInfo.APIGatewayPath.AppendPath(*reqURL) //doubt
		// method := normaliseHttpMethod(r.Method)

		if strings.HasPrefix(path.Value, s.AppInfo.APIBasePath.Value) == false {
			return
		}
		requestRID := getRIDFromRequest(r)
		if requestRID == "" {
			var matchedRecipe RecipeModule
			for _, recipeModule := range s.RecipeModules {
				if recipeModule.GetRecipeID() == requestRID {
					matchedRecipe = recipeModule
					break
				}
			}
			if reflect.DeepEqual(matchedRecipe, RecipeModule{}) {
				return
			}
		}
	})
}
