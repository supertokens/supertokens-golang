/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

// This file implements a minimal in-memory WebAuthn authenticator so the tests
// can drive real registration/authentication ceremonies against the core (which
// does real cryptographic verification). It uses only the standard library: an
// ES256 (P-256) key, hand-rolled CBOR for the attestation object / COSE key, and
// ECDSA assertion signatures.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
)

// --- tiny CBOR encoder (only what the attestation object / COSE key need) ---

func cborHead(major byte, n uint64) []byte {
	mt := major << 5
	switch {
	case n < 24:
		return []byte{mt | byte(n)}
	case n < 1<<8:
		return []byte{mt | 24, byte(n)}
	case n < 1<<16:
		return []byte{mt | 25, byte(n >> 8), byte(n)}
	case n < 1<<32:
		return []byte{mt | 26, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
	default:
		b := []byte{mt | 27}
		return append(b, byte(n>>56), byte(n>>48), byte(n>>40), byte(n>>32), byte(n>>24), byte(n>>16), byte(n>>8), byte(n))
	}
}

func cborUint(n uint64) []byte { return cborHead(0, n) }

// cborNegInt encodes a negative integer (v < 0).
func cborNegInt(v int64) []byte { return cborHead(1, uint64(-1-v)) }

func cborBytes(b []byte) []byte { return append(cborHead(2, uint64(len(b))), b...) }

func cborText(s string) []byte { return append(cborHead(3, uint64(len(s))), []byte(s)...) }

func cborMapHeader(n uint64) []byte { return cborHead(5, n) }

func b64url(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

// virtualAuthenticator is a single ES256 credential.
type virtualAuthenticator struct {
	key    *ecdsa.PrivateKey
	credID []byte
	aaguid []byte
}

func newVirtualAuthenticator() *virtualAuthenticator {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	credID := make([]byte, 20)
	if _, err := rand.Read(credID); err != nil {
		panic(err)
	}
	return &virtualAuthenticator{key: key, credID: credID, aaguid: make([]byte, 16)}
}

func (a *virtualAuthenticator) credentialIDB64() string { return b64url(a.credID) }

// coseKey encodes the public key as a COSE_Key (EC2 / ES256 / P-256).
func (a *virtualAuthenticator) coseKey() []byte {
	x := make([]byte, 32)
	y := make([]byte, 32)
	a.key.X.FillBytes(x)
	a.key.Y.FillBytes(y)

	var b []byte
	b = append(b, cborMapHeader(5)...)
	b = append(b, cborUint(1)...)    // kty
	b = append(b, cborUint(2)...)    // EC2
	b = append(b, cborUint(3)...)    // alg
	b = append(b, cborNegInt(-7)...) // ES256
	b = append(b, cborNegInt(-1)...) // crv
	b = append(b, cborUint(1)...)    // P-256
	b = append(b, cborNegInt(-2)...) // x
	b = append(b, cborBytes(x)...)
	b = append(b, cborNegInt(-3)...) // y
	b = append(b, cborBytes(y)...)
	return b
}

// authData builds authenticatorData. When attested is true the attested
// credential data (aaguid + credId + COSE key) is appended — used at
// registration. UP and UV flags are always set.
func (a *virtualAuthenticator) authData(rpID string, attested bool) []byte {
	rpHash := sha256.Sum256([]byte(rpID))
	flags := byte(0x01 | 0x04) // UP | UV
	if attested {
		flags |= 0x40 // AT
	}

	b := append([]byte{}, rpHash[:]...)
	b = append(b, flags)
	b = append(b, 0, 0, 0, 0) // signCount
	if attested {
		b = append(b, a.aaguid...)
		var l [2]byte
		binary.BigEndian.PutUint16(l[:], uint16(len(a.credID)))
		b = append(b, l[:]...)
		b = append(b, a.credID...)
		b = append(b, a.coseKey()...)
	}
	return b
}

func clientDataJSON(typ, challenge, origin string) []byte {
	// challenge is the base64url string returned by the core's options endpoint
	// and is placed verbatim into clientDataJSON.
	m := map[string]interface{}{
		"type":        typ,
		"challenge":   challenge,
		"origin":      origin,
		"crossOrigin": false,
	}
	out, _ := json.Marshal(m)
	return out
}

// registrationResponse builds the credential payload for /signup (attestation
// "none", so the core does not verify a signature here).
func (a *virtualAuthenticator) registrationResponse(rpID, challenge, origin string) map[string]interface{} {
	clientData := clientDataJSON("webauthn.create", challenge, origin)
	authData := a.authData(rpID, true)

	var attObj []byte
	attObj = append(attObj, cborMapHeader(3)...)
	attObj = append(attObj, cborText("fmt")...)
	attObj = append(attObj, cborText("none")...)
	attObj = append(attObj, cborText("attStmt")...)
	attObj = append(attObj, cborMapHeader(0)...)
	attObj = append(attObj, cborText("authData")...)
	attObj = append(attObj, cborBytes(authData)...)

	return map[string]interface{}{
		"id":                      a.credentialIDB64(),
		"rawId":                   a.credentialIDB64(),
		"authenticatorAttachment": "platform",
		"type":                    "public-key",
		"clientExtensionResults":  map[string]interface{}{},
		"response": map[string]interface{}{
			"clientDataJSON":    b64url(clientData),
			"attestationObject": b64url(attObj),
			"transports":        []string{"internal"},
		},
	}
}

// assertionResponse builds the credential payload for /signin (a real ES256
// signature over authenticatorData || sha256(clientDataJSON)).
func (a *virtualAuthenticator) assertionResponse(rpID, challenge, origin, userHandle string) map[string]interface{} {
	clientData := clientDataJSON("webauthn.get", challenge, origin)
	authData := a.authData(rpID, false)

	clientDataHash := sha256.Sum256(clientData)
	signed := append(append([]byte{}, authData...), clientDataHash[:]...)
	digest := sha256.Sum256(signed)
	sig, err := ecdsa.SignASN1(rand.Reader, a.key, digest[:])
	if err != nil {
		panic(err)
	}

	resp := map[string]interface{}{
		"clientDataJSON":    b64url(clientData),
		"authenticatorData": b64url(authData),
		"signature":         b64url(sig),
	}
	if userHandle != "" {
		resp["userHandle"] = b64url([]byte(userHandle))
	}
	return map[string]interface{}{
		"id":                      a.credentialIDB64(),
		"rawId":                   a.credentialIDB64(),
		"authenticatorAttachment": "platform",
		"type":                    "public-key",
		"clientExtensionResults":  map[string]interface{}{},
		"response":                resp,
	}
}
