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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumaniseMilliseconds(t *testing.T) {
	assert.Equal(t, "1 second", HumaniseMilliseconds(1000))
	assert.Equal(t, "59 seconds", HumaniseMilliseconds(59000))
	assert.Equal(t, "1 minute", HumaniseMilliseconds(60000))
	assert.Equal(t, "1 minute", HumaniseMilliseconds(119000))
	assert.Equal(t, "2 minutes", HumaniseMilliseconds(120000))
	assert.Equal(t, "1 hour", HumaniseMilliseconds(3600000))
	assert.Equal(t, "1 hour", HumaniseMilliseconds(3660000))
	assert.Equal(t, "1.1 hours", HumaniseMilliseconds(3960000))
	assert.Equal(t, "2 hours", HumaniseMilliseconds(7260000))
	assert.Equal(t, "5 hours", HumaniseMilliseconds(18000000))
}
