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

package epunittesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
)

func TestDefaultEmailValidator(t *testing.T) {
	assert.Nil(t, emailpassword.DefaultEmailValidator("test@supertokens.io"))
	assert.Nil(t, emailpassword.DefaultEmailValidator("nsdafa@gmail.com"))
	assert.Nil(t, emailpassword.DefaultEmailValidator("fewf3r_fdkj@gmaildsfa.co.uk"))
	assert.Nil(t, emailpassword.DefaultEmailValidator("dafk.adfa@gmail.com"))
	assert.Nil(t, emailpassword.DefaultEmailValidator("skjlblc3f3@fnldsks.co"))
	assert.Equal(t, "Email is invalid", *emailpassword.DefaultEmailValidator("sdkjfnas34@gmail.com.c"))
	assert.Equal(t, "Email is invalid", *emailpassword.DefaultEmailValidator("d@c"))
	assert.Equal(t, "Email is invalid", *emailpassword.DefaultEmailValidator("fasd"))
	assert.Equal(t, "Email is invalid", *emailpassword.DefaultEmailValidator("dfa@@@abc.com"))
	assert.Equal(t, "Email is invalid", *emailpassword.DefaultEmailValidator(""))
}

func TestDefaultPasswordValidator(t *testing.T) {
	assert.Nil(t, emailpassword.DefaultPasswordValidator("dsknfkf38H"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("lasdkf*787~sdfskj"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("L0493434505"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("3453342422L"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("1sdfsdfsdfsd"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("dksjnlvsnl2"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("abcgftr8"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("  d3    "))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("abc!@#$%^&*()gftr8"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("    dskj3"))
	assert.Nil(t, emailpassword.DefaultPasswordValidator("    dsk  3"))

	assert.Equal(t, "Password must contain at least 8 characters, including a number", *emailpassword.DefaultPasswordValidator("asd"))

	assert.Equal(t, "Password's length must be lesser than 100 characters", *emailpassword.DefaultPasswordValidator("asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4asdfdefrg4"))

	assert.Equal(t, "Password must contain at least one number", *emailpassword.DefaultPasswordValidator("ascdvsdfvsIUOO"))

	assert.Equal(t, "Password must contain at least one alphabet", *emailpassword.DefaultPasswordValidator("234235234523"))
}
