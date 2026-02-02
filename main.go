package main

import (
	"golang-day23/restapi"
)

func main() {
	// freeapi.RandomQuote()
	restapi.ConnectDB()
	restapi.Handler()
}
