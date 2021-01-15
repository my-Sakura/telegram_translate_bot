package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://openapi.youdao.com/api?q=中国&from=auto&to=en&appKey=03f9ad9d055d2a3d", nil)
	if err != nil {
		fmt.Println(err)
	}

	//req.Header.Add("Content-Type", "application/text")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println()
	}

	fmt.Println(string(body))
}
