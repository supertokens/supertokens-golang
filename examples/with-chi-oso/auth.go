package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/osohq/go-oso"
	"github.com/supertokens/supertokens-golang/examples/with-chi-oso/models"
)

func NewAuth() oso.Oso {
	osoClient, err := oso.NewOso()
	if err != nil {
		fmt.Sprintf("Failed to set up Oso: %v", err)
		os.Exit(1)
	}
	osoClient.RegisterClass(reflect.TypeOf(models.Repository{}), nil)
	osoClient.RegisterClass(reflect.TypeOf(models.User{}), nil)
	if err = osoClient.LoadFiles([]string{"main.polar"}); err != nil {
		fmt.Sprintf("Failed to start: %s", err)
		os.Exit(1)
	}
	return osoClient
}
