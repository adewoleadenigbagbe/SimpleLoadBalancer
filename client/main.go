package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "http://localhost:4091/"
	for i := 0; i < 500; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(body))
	}
}
