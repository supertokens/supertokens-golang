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

package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormaliseSessionScope(t *testing.T) {
	t.Run("test with leading dot", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError(".example.com")
		assert.NoError(t, err)
		assert.Equal(t, ".example.com", *result)
	})

	t.Run("test without leading dot", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError("example.com")
		assert.NoError(t, err)
		assert.Equal(t, "example.com", *result)
	})

	t.Run("test with http prefix", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError("http://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "example.com", *result)
	})

	t.Run("test with https prefix", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError("https://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "example.com", *result)
	})

	t.Run("test with IP address", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError("192.168.1.1")
		assert.NoError(t, err)
		assert.Equal(t, "192.168.1.1", *result)
	})

	t.Run("test with localhost", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError("localhost")
		assert.NoError(t, err)
		assert.Equal(t, "localhost", *result)
	})

	t.Run("test with leading and trailing whitespace", func(t *testing.T) {
		result, err := normaliseSessionScopeOrThrowError("  example.com  ")
		assert.NoError(t, err)
		assert.Equal(t, "example.com", *result)
	})
}
