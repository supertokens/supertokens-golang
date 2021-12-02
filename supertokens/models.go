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

package supertokens

import "net/http"

type NormalisedAppinfo struct {
	AppName         string
	WebsiteDomain   NormalisedURLDomain
	APIDomain       NormalisedURLDomain
	APIBasePath     NormalisedURLPath
	APIGatewayPath  NormalisedURLPath
	WebsiteBasePath NormalisedURLPath
}

type AppInfo struct {
	AppName         string
	WebsiteDomain   string
	APIDomain       string
	WebsiteBasePath *string
	APIBasePath     *string
	APIGatewayPath  *string
}

type Recipe func(appInfo NormalisedAppinfo, onGeneralError func(err error, req *http.Request, res http.ResponseWriter)) (*RecipeModule, error)

type TypeInput struct {
	Supertokens    *ConnectionInfo
	AppInfo        AppInfo
	RecipeList     []Recipe
	Telemetry      *bool
	OnGeneralError func(err error, req *http.Request, res http.ResponseWriter)
}

type ConnectionInfo struct {
	ConnectionURI string
	APIKey        string
}

type APIHandled struct {
	PathWithoutAPIBasePath NormalisedURLPath
	Method                 string
	ID                     string
	Disabled               bool
}

type DoneWriter struct {
	http.ResponseWriter
	done bool
}

func (w *DoneWriter) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}

type UserContext = *interface{}
