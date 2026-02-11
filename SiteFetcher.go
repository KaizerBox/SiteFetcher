package main

import (
	"context"
	"fmt"
	"net/http"
	"siteFetcher/components/fetcher"
	"siteFetcher/components/gmailer"
)

type FetcherMailer struct {
	Fetcher fetcher.Fetcher
	Emailer gmailer.Gmailer
}

func main() {
	fetcherContext := context.Background()

	fetcherClient := &http.Client{}

	fetcherSetting := fetcher.FetcherSetting{
		MaxFetchTimeLimit: 10,
	}

	fetcher := fetcher.Fetcher{
		FetchContext:   fetcherContext,
		Client:         fetcherClient,
		FetcherSetting: fetcherSetting,
	}

	respBody, err := fetcher.FetchWithUrl("url")

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(respBody)

}
