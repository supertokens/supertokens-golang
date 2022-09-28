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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func AppleRedirectHandler(apiImplementation tpmodels.APIInterface, options tpmodels.APIOptions) error {
	if apiImplementation.AppleRedirectHandlerPOST == nil || (*apiImplementation.AppleRedirectHandlerPOST) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	err := options.Req.ParseMultipartForm(0)
	if err != nil {
		return err
	}

	infoFromProvider := map[string]interface{}{}

	for key, value := range options.Req.PostForm {
		infoFromProvider[key] = value[0]
	}

	return (*apiImplementation.AppleRedirectHandlerPOST)(infoFromProvider, options, supertokens.MakeDefaultUserContextFromAPI(options.Req))
}
