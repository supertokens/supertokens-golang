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

type test struct {
	someFunc      *func()
	someOtherFunc *func()
}

func TestOverride(t *testing.T) {
	m := 1
	someOtherFunc := func() {
		m = 2
	}
	someFunc := func() {
		someOtherFunc()
	}

	oI := test{
		someFunc:      &someFunc,
		someOtherFunc: &someOtherFunc,
	}

	originalSomeFunc := *oI.someFunc

	(*oI.someFunc) = func() {
		originalSomeFunc()
	}
	(*oI.someOtherFunc) = func() {
		m = 3
	}

	(*oI.someFunc)()
	assert.Equal(t, 3, m)
}

type test2Ep struct {
	SignIn         *func()
	GetUserByEmail *func()
}

type test2TpEp struct {
	SignIn          *func()
	GetUsersByEmail *func()
}

func TestOverride2(t *testing.T) {
	m := 0

	getUserByEmail := func() {
		m = 1
	}

	signIn := func() {
		getUserByEmail()
	}

	ep := test2Ep{
		SignIn:         &signIn,
		GetUserByEmail: &getUserByEmail,
	}

	///////////////////

	oSignIn := *ep.SignIn
	signInTpep := func() {
		oSignIn()
	}

	oGetUserByEmail := *ep.GetUserByEmail
	getUsersByEmailTpep := func() {
		oGetUserByEmail()
	}

	tpep := test2TpEp{
		SignIn:          &signInTpep,
		GetUsersByEmail: &getUsersByEmailTpep,
	}

	derivedEp := func(tpepInstance test2TpEp) test2Ep {
		s := func() {
			(*tpepInstance.SignIn)()
		}
		g := func() {
			(*tpepInstance.GetUsersByEmail)()
		}

		return test2Ep{
			SignIn:         &s,
			GetUserByEmail: &g,
		}
	}

	derived := derivedEp(tpep)
	(*ep.GetUserByEmail) = (*derived.GetUserByEmail)
	(*ep.SignIn) = (*derived.SignIn)

	//////////////////

	// user override
	originalGetUsersByEmail := *tpep.GetUsersByEmail
	(*tpep.GetUsersByEmail) = func() {
		m = 5
		originalGetUsersByEmail()
		if m == 1 {
			m = 2
		}
	}
	originalSignIn := *tpep.SignIn
	(*tpep.SignIn) = func() {
		originalSignIn()
	}

	////////////////
	// usage of tpep by functions exposed by tpep
	(*tpep.SignIn)()
	assert.Equal(t, 2, m)

	//////////////////
	// usage of ep recipe functions inside ep APIs
	m = 1

	toEpRecipe := derivedEp(tpep)

	(*toEpRecipe.GetUserByEmail)()

	assert.Equal(t, 2, m)

}
