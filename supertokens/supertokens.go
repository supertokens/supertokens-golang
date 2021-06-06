package supertokens

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SuperTokens struct {
	Instance      *SuperTokens
	AppInfo       NormalisedAppinfo
	RecipeModules []RecipeModule
}

var s SuperTokens

func newSuperTokens(config TypeInput) (*SuperTokens, error) {
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
			return nil, err
		}
		hosts = append(hosts, *host)
	}

	q := Querier{}
	q.Init(hosts, config.Supertoken.APIKey)

	if config.RecipeList == nil || len(config.RecipeList) == 0 {
		return nil, errors.New("Please provide at least one recipe to the supertokens.init function call")
	}

	for _, elem := range config.RecipeList {
		s.RecipeModules = append(s.RecipeModules, elem(s.AppInfo))
	}

	telemetry := config.Telemetry

	if telemetry {
		// err := s.SendTelemetry()
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func SupertokensInit(config TypeInput) (err error) {
	if s.Instance == nil {
		s.Instance, err = newSuperTokens(config)
		return err
	}
	return nil
}

func (s *SuperTokens) GetInstanceOrThrowError() (*SuperTokens, error) {
	if s.Instance != nil {
		return s.Instance, nil
	}
	return nil, errors.New("Initialisation not done. Did you forget to call the SuperTokens.init function?")
}

func (s *SuperTokens) SendTelemetry() {
	q := &Querier{}
	querier, err := q.GetNewInstanceOrThrowError("")
	if err != nil {
		fmt.Println(err) // todo: handle error
		return
	}

	response, err := querier.SendGetRequest(NormalisedURLPath{value: "/telemetry"}, map[string]string{})
	if err != nil {
		fmt.Println(err) // todo: handle error
		return
	}
	var telemetryID string
	exists, err := strconv.ParseBool(response["exists"])
	if err == nil && exists == true {
		telemetryID = response["telemetryId"]
	}

	url := "https://api.supertokens.io/0/st/telemetry"

	data := map[string]interface{}{
		"appName":       s.AppInfo.AppName,
		"websiteDomain": s.AppInfo.WebsiteDomain.GetAsStringDangerous(),
		"telemetryId":   telemetryID,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err) // todo: handle error
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err) // todo: handle error
		return
	}
	req.Header.Set("api-version", "2")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err) // todo: handle error
		return
	}
}

func (s *SuperTokens) Middleware() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqURL, _ := NewNormalisedURLPath(r.RemoteAddr) // todo: error handle
		path := s.AppInfo.APIGatewayPath.AppendPath(*reqURL)
		method := r.Method

		if !strings.HasPrefix(path.GetAsStringDangerous(), s.AppInfo.APIBasePath.GetAsStringDangerous()) {
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
				return
			}

			id := matchedRecipe.ReturnAPIIdIfCanHandleRequest(path, method)
			if id == "" {
				return
			}
			s.HandleAPI(*matchedRecipe, id, r, w, path, method)
		} else {
			for _, recipeModule := range s.RecipeModules {
				id := recipeModule.ReturnAPIIdIfCanHandleRequest(path, method)
				if id != "" {
					s.HandleAPI(recipeModule, id, r, w, path, method)
					return
				}
			}
		}
	})
}

func (s *SuperTokens) HandleAPI(matchedRecipe RecipeModule,
	id string,
	r *http.Request,
	w http.ResponseWriter,
	path NormalisedURLPath,
	method string) {
	matchedRecipe.HandleAPIRequest(id, r, w, path, method)
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
