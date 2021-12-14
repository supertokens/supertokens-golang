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

package session

import (
	"context"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func VerifySessionHelper(recipeInstance Recipe, options *sessmodels.VerifySessionOptions, otherHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dw := &supertokens.DoneWriter{ResponseWriter: w}
		session, err := (*recipeInstance.APIImpl.VerifySession)(options, sessmodels.APIOptions{
			Config:               recipeInstance.Config,
			OtherHandler:         otherHandler,
			Req:                  r,
			Res:                  dw,
			RecipeID:             recipeInstance.RecipeModule.GetRecipeID(),
			RecipeImplementation: recipeInstance.RecipeImpl,
		}, &map[string]interface{}{})
		if err != nil {
			err = supertokens.ErrorHandler(err, r, dw)
			if err != nil {
				recipeInstance.RecipeModule.OnGeneralError(err, r, dw)
			}
			return
		}
		if session != nil {
			ctx := context.WithValue(r.Context(), sessmodels.SessionContext, session)
			otherHandler(dw, r.WithContext(ctx))
		} else {
			otherHandler(dw, r)
		}
	})
}
