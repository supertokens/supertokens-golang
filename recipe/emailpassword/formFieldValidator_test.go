/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package emailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEmailValidator(t *testing.T) {
	assert.Nil(t, defaultEmailValidator("test@supertokens.io"))
	assert.Nil(t, defaultEmailValidator("nsdafa@gmail.com"))
	assert.Nil(t, defaultEmailValidator("fewf3r_fdkj@gmaildsfa.co.uk"))
	assert.Nil(t, defaultEmailValidator("dafk.adfa@gmail.com"))
	assert.Nil(t, defaultEmailValidator("skjlblc3f3@fnldsks.co"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("sdkjfnas34@gmail.com.c"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("d@c"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("fasd"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator("dfa@@@abc.com"))
	assert.Equal(t, "Email is invalid", *defaultEmailValidator(""))
}

func TestDefaultPasswordValidator(t *testing.T) {
	assert.Nil(t, defaultPasswordValidator("dsknfkf38H"))
	assert.Nil(t, defaultPasswordValidator("lasdkf*787~sdfskj"))
	assert.Nil(t, defaultPasswordValidator("L0493434505"))
	assert.Nil(t, defaultPasswordValidator("3453342422L"))
	assert.Nil(t, defaultPasswordValidator("1sdfsdfsdfsd"))
	assert.Nil(t, defaultPasswordValidator("dksjnlvsnl2"))
	assert.Nil(t, defaultPasswordValidator("abcgftr8"))
	assert.Nil(t, defaultPasswordValidator("  d3    "))
	assert.Nil(t, defaultPasswordValidator("abc!@#$%^&*()gftr8"))
	assert.Nil(t, defaultPasswordValidator("    dskj3"))
	assert.Nil(t, defaultPasswordValidator("    dsk  3"))

	assert.Equal(t, "Password must contain at least 8 characters, including a number", *defaultPasswordValidator("asd"))

	assert.Equal(t, "Password's length must be lesser than 100 characters", *defaultPasswordValidator("asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4"))

	assert.Equal(t, "Password must contain at least one number", *defaultPasswordValidator("ascdvsdfvsIUOO"))

	assert.Equal(t, "Password must contain at least one alphabet", *defaultPasswordValidator("234235234523"))
}
