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
		connectionURI := stInstance.SuperTokens.ConnectionURI

		normalizedDashboardPath, err := supertokens.NewNormalisedURLPath(dashboardAPI)
		if err != nil {
			return "", err
		}
		dashboardAppPath := options.AppInfo.APIBasePath.AppendPath(normalizedDashboardPath).GetAsStringDangerous()

		authMode := string(options.Config.AuthMode)

		return `
		<html>
		<head>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<script>
						window.staticBasePath = "` + bundleDomain + `/static"
						window.dashboardAppPath = "` + dashboardAppPath + `"
						window.connectionURI = "` + connectionURI + `"
						window.authMode = "` + authMode + `"
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
