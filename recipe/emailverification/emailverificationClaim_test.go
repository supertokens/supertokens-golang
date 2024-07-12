/*
 * Copyright (c) 2024, VRAI Labs and/or its affiliates. All rights reserved.
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

package emailverification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evclaims"
)

func TestEmailVerificationClaim(t *testing.T) {
	t.Run("value should be fetched if it is nil", func(t *testing.T) {
		validator := evclaims.EmailVerificationClaimValidators.IsVerified(nil, nil)

		shouldRefreshNil := validator.ShouldRefetch(nil, nil)

		assert.True(t, shouldRefreshNil)
	})

	t.Run("value should be fetched as per maxAgeInSeconds if it is provided", func(t *testing.T) {
		refetchTimeOnFalseInSeconds := int64(10)
		maxAgeInSeconds := int64(200)
		validator := evclaims.EmailVerificationClaimValidators.IsVerified(&refetchTimeOnFalseInSeconds, &maxAgeInSeconds)

		payload := map[string]interface{}{
			"st-ev": map[string]interface{}{
				"v": true,
				"t": time.Now().UnixMilli() - 199*1000,
			},
		}

		shouldRefreshValid := validator.ShouldRefetch(payload, nil)

		assert.False(t, shouldRefreshValid)

		payload = map[string]interface{}{
			"st-ev": map[string]interface{}{
				"v": true,
				"t": time.Now().UnixMilli() - 201*1000,
			},
		}

		shouldRefreshExpired := validator.ShouldRefetch(payload, nil)
		assert.True(t, shouldRefreshExpired)
	})

	t.Run("value should be fetched as per refetchTimeOnFalseInSeconds if it is provided", func(t *testing.T) {
		refetchTimeOnFalseInSeconds := int64(8)
		validator := evclaims.EmailVerificationClaimValidators.IsVerified(&refetchTimeOnFalseInSeconds, nil)

		payload := map[string]interface{}{
			"st-ev": map[string]interface{}{
				"v": false,
				"t": time.Now().UnixMilli() - 7*1000,
			},
		}

		shouldRefreshValid := validator.ShouldRefetch(payload, nil)

		assert.False(t, shouldRefreshValid)

		payload = map[string]interface{}{
			"st-ev": map[string]interface{}{
				"v": false,
				"t": time.Now().UnixMilli() - 9*1000,
			},
		}

		shouldRefreshExpired := validator.ShouldRefetch(payload, nil)
		assert.True(t, shouldRefreshExpired)
	})

	t.Run("value should be fetched as per default the refetchTimeOnFalseInSeconds if it is not provided", func(t *testing.T) {
		validator := evclaims.EmailVerificationClaimValidators.IsVerified(nil, nil)

		// NOTE: the default value of refetchTimeOnFalseInSeconds is 10 seconds
		payload := map[string]interface{}{
			"st-ev": map[string]interface{}{
				"v": false,
				"t": time.Now().UnixMilli() - 9*1000,
			},
		}

		shouldRefreshValid := validator.ShouldRefetch(payload, nil)

		assert.False(t, shouldRefreshValid)

		payload = map[string]interface{}{
			"st-ev": map[string]interface{}{
				"v": false,
				"t": time.Now().UnixMilli() - 11*1000,
			},
		}

		shouldRefreshExpired := validator.ShouldRefetch(payload, nil)
		assert.True(t, shouldRefreshExpired)
	})
}
