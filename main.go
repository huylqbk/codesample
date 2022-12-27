package main

import (
	"fmt"
	"time"

	"github.com/huylqbk/codesample/httpclient"
)

func main() {
	r, err := httpclient.New().
		SetURL("https://api.agify.io/?name=bella").
		SetMethod("GET").
		SetTimeout(60 * time.Second).
		Execute()
	fmt.Println(string(r), err)
}
