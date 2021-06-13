package session

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"

	"github.com/supertokens/supertokens-golang/errors"
)

/*
{
	"alg":     "RS256",
	"typ":     "JWT",
	"version": "2",
}
*/
const header = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsInZlcnNpb24iOiIyIn0="

func verifyJWTAndGetPayload(jwt string, jwtSigningPublicKey string) (map[string]interface{}, error) {
	var splitted = strings.Split(jwt, ".")
	if len(splitted) != 3 {
		return nil, errors.GeneralError{
			Msg: "Invalid JWT",
		}
	}
	if header != splitted[0] {
		return nil, errors.GeneralError{
			Msg: "JWT header mismatch",
		}
	}
	var payload = splitted[1]

	var publicKey, publicKeyError = getPublicKeyFromStr("-----BEGIN PUBLIC KEY-----\n" + jwtSigningPublicKey + "\n-----END PUBLIC KEY-----")
	if publicKeyError != nil {
		return nil, publicKeyError
	}

	h := sha256.New()
	h.Write([]byte(header + "." + payload))
	digest := h.Sum(nil)

	var decodedSignature, decodedSignatureError = b64.StdEncoding.DecodeString(splitted[2])
	if decodedSignatureError != nil {
		return nil, decodedSignatureError
	}

	verificationError := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest, decodedSignature)
	if verificationError != nil {
		return nil, verificationError
	}

	var decodedPayload, base64Error = b64.StdEncoding.DecodeString(payload)
	if base64Error != nil {
		return nil, base64Error
	}

	var result map[string]interface{}
	jsonError := json.Unmarshal(decodedPayload, &result)
	if jsonError != nil {
		return nil, jsonError
	}
	return result, nil
}

func getPublicKeyFromStr(str string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(str))
	if block == nil {
		return nil, errors.GeneralError{
			Msg: "failed to parse PEM block containing the public key",
		}
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.GeneralError{
			Msg: "failed to parse DER encoded public key:" + err.Error(),
		}
	}

	return pub.(*rsa.PublicKey), nil
}
