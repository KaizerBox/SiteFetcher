package main

import (
	"fmt"
	"net/http"

	"siteFetcher/internal/stores"
)

func main() {
	fetcherClient := &http.Client{}

	store, err := stores.NewStore("bestbuy", fetcherClient)
	if err != nil {
		fmt.Printf("Error creating store: %v\n", err)
		return
	}

	finalResult, fetchTime, err := store.GetResponse("19869815", "V8P")
	if err != nil {
		fmt.Printf("Error getting final product data: %v\n", err)
		return
	}

	fmt.Printf("Final Product Data: %+v\n", finalResult)
	fmt.Printf("Fetch Time: %v\n", fetchTime)
}
