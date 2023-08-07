package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type analyticsPostResponse struct {
	Status string `json:"status"`
}

type analyticsPostRequestBody struct {
	Email            *string `json:"email"`
	DashboardVersion *string `json:"dashboardVersion"`
}

func AnalyticsPost(apiInterface dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (analyticsPostResponse, error) {
	supertokensInstance, instanceError := supertokens.GetInstanceOrThrowError()

	if supertokens.IsRunningInTestMode() {
		return analyticsPostResponse{
			Status: "OK",
		}, nil
	}

	if instanceError != nil {
		return analyticsPostResponse{}, instanceError
	}

	if supertokensInstance.Telemetry != nil && !*supertokensInstance.Telemetry {
		return analyticsPostResponse{
			Status: "OK",
		}, nil
	}

	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return analyticsPostResponse{}, err
	}

	var readBody analyticsPostRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return analyticsPostResponse{}, err
	}

	if readBody.Email == nil {
		return analyticsPostResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'email' is missing",
		}
	}

	if readBody.DashboardVersion == nil {
		return analyticsPostResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'dashboardVersion' is missing",
		}
	}

	data := map[string]interface{}{
		"websiteDomain":    supertokensInstance.AppInfo.WebsiteDomain.GetAsStringDangerous(),
		"apiDomain":        supertokensInstance.AppInfo.APIDomain.GetAsStringDangerous(),
		"appName":          supertokensInstance.AppInfo.AppName,
		"sdk":              "golang",
		"sdkVersion":       supertokens.VERSION,
		"email":            *readBody.Email,
		"dashboardVersion": *readBody.DashboardVersion,
	}

	querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return analyticsPostResponse{}, err
	}

	response, err := querier.SendGetRequest("/telemetry", nil)
	if err != nil {
		// We don't send telemetry events if this fails
		return analyticsPostResponse{
			Status: "OK",
		}, nil
	}

	exists := response["exists"].(bool)

	if exists {
		data["telemetryId"] = response["telemetryId"].(string)
	}

	numberOfUsers, err := supertokens.GetUserCount(nil, nil)
	if err != nil {
		// We don't send telemetry events if this fails
		return analyticsPostResponse{
			Status: "OK",
		}, nil
	}

	data["numberOfUsers"] = numberOfUsers

	jsonData, err := json.Marshal(data)
	if err != nil {
		return analyticsPostResponse{}, err
	}

	url := "https://api.supertokens.com/0/st/telemetry"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return analyticsPostResponse{
			Status: "OK",
		}, nil
	}
	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("api-version", "3")
	client := &http.Client{}
	client.Do(req)

	return analyticsPostResponse{
		Status: "OK",
	}, nil
}
