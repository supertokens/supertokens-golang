/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package api

import (
	"strconv"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/constants"
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() dashboardmodels.APIInterface {
	dashboardGET := func(options dashboardmodels.APIOptions, userContext supertokens.UserContext) (string, error) {
		bundleBasePathString, err := (*options.RecipeImplementation.GetDashboardBundleLocation)(userContext)
		if err != nil {
			return "", err
		}

		normalizedDomain, err := supertokens.NewNormalisedURLDomain(bundleBasePathString)
		if err != nil {
			return "", err
		}
		normalizedPath, err := supertokens.NewNormalisedURLPath(bundleBasePathString)
		if err != nil {
			return "", err
		}

		bundleDomain := normalizedDomain.GetAsStringDangerous() + normalizedPath.GetAsStringDangerous()

		stInstance, err := supertokens.GetInstanceOrThrowError()
		if err != nil {
			return "", err
		}

		// We are splitting the passed URI here so that if multiple URI's are passed
		// separated by a colon, the first one is returned.
		connectionURIToNormalize := strings.Split(stInstance.SuperTokens.ConnectionURI, ";")[0]

		// This normalizes the URI to make sure that it has things like protocol etc
		// injected into it before it is returned.
		var normalizationError error
		normalizedConnectionURI, normalizationError := supertokens.NewNormalisedURLDomain(connectionURIToNormalize)
		if normalizationError != nil {
			// In case of failures, we want to return a 500 here, mainly because that
			// is what we return if the connectionURI is invalid which is the case here
			// if normalization fails.
			return "", normalizationError
		}
		connectionURI := normalizedConnectionURI.GetAsStringDangerous()

		normalizedDashboardPath, err := supertokens.NewNormalisedURLPath(constants.DashboardAPI)
		if err != nil {
			return "", err
		}
		dashboardAppPath := options.AppInfo.APIBasePath.AppendPath(normalizedDashboardPath).GetAsStringDangerous()

		authMode := string(options.Config.AuthMode)

		isSearchEnabled := false
		querier, err := supertokens.GetNewQuerierInstanceOrThrowError(options.RecipeID)
		if err != nil {
			return "", err
		}
		cdiVersion, err := querier.GetQuerierAPIVersion(userContext)
		if err != nil {
			return "", err
		}
		if supertokens.MaxVersion(cdiVersion, "2.20") == cdiVersion {
			// Only enable search for CDI version 2.20 and above
			isSearchEnabled = true
		}

		return `
		<html>
		<head>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<script>
						window.staticBasePath = "` + bundleDomain + `/static"
						window.dashboardAppPath = "` + dashboardAppPath + `"
						window.connectionURI = "` + connectionURI + `"
						window.authMode = "` + authMode + `"
						window.isSearchEnabled = "` + strconv.FormatBool(isSearchEnabled) + `"
				</script>
				<script defer src="` + bundleDomain + `/static/js/bundle.js"></script></head>
				<link href="` + bundleDomain + `/static/css/main.css" rel="stylesheet" type="text/css">
				<link rel="icon" type="image/x-icon" href="` + bundleDomain + `/static/media/favicon.ico">
		</head>
		<body>
				<noscript>You need to enable JavaScript to run this app.</noscript>
				<div id="root"></div>
		</body>
		</html>
		`, nil
	}

	return dashboardmodels.APIInterface{
		DashboardGET: &dashboardGET,
	}
}
