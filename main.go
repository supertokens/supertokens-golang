package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	userID := "user"
	var JWTPayload interface{} = map[string]interface{}{}
	var sessionData interface{} = nil
	requestBody := map[string]interface{}{
		"userId":             userID,
		"userDataInJWT":      JWTPayload,
		"userDataInDatabase": sessionData,
	}
	res2B, _ := json.Marshal(requestBody)
	fmt.Println(string(res2B))
}
