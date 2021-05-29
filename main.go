package main

import (
	"fmt"
	"net/url"
)

func main() {
	input := "localhost"
	urlObj, err := url.Parse(input)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", *urlObj)
}
