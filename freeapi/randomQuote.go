package freeapi

import (
	"fmt"
	"io"
	"net/http"
)

func RandomQuote() {
	myURL := "https://api.freeapi.app/api/v1/public/quotes/quote/random"

	req, _ := http.NewRequest("GET", myURL, nil)
	req.Header.Add("Accept", "application/json")
	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	fmt.Println(res)
	fmt.Println(string(body))
}
