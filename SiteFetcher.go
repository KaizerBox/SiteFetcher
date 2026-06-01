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

	reqUrlBuilder := fetcher.BestBuyRequestUrlBuilder{
		Name:            "BestBuy",
		Scheme:          "https",
		Host:            "www.bestbuy.ca",
		ProductURL:      "/api/v2/json/product/%s",
		LocationUrl:     "/api/v3/json/locations",
		AvailabilityUrl: "/ecomm-api/availability/products",
	}

	//Get Items info
	productUrl, err := reqUrlBuilder.GetProductUrl("19869815")

	if err != nil {
		fmt.Printf("Error getting product URL: %v\n", err)
		return
	}

	reqBuilder := fetcher.BestBuyRequestBuilder{Client: fetcherClient}

	reqParam := fetcher.RequestParameters{
		RequestUrl:     productUrl,
		RequestContext: fetcherContext,
	}

	responseData, _, err := reqBuilder.FetchWithRequestParam(fetcherClient, reqParam)

	if err != nil {
		fmt.Printf("Error fetching product data: %v\n", err)
		return
	}

	responseParser := fetcher.BestBuyResponseParser{}

	parsedResult, err := responseParser.ParseResponseMessage(responseData)

	if err != nil {
		fmt.Printf("Error parsing product data: %v\n", err)
		return
	}

	fmt.Printf("Parsed Product Data: %+v\n", parsedResult)
	fmt.Printf("------------------------------\n")

	//Get Locations/Stores Info
	locationsUrl, err := reqUrlBuilder.GetLocationsUrl("V8P")
	if err != nil {
		fmt.Printf("Error getting locations URL: %v\n", err)
		return
	}
	locationsReqParam := fetcher.RequestParameters{
		RequestUrl:     locationsUrl,
		RequestContext: fetcherContext,
	}

	locationsResponseData, _, err := reqBuilder.FetchWithRequestParam(fetcherClient, locationsReqParam)

	if err != nil {
		fmt.Printf("Error fetching locations data: %v\n", err)
		return
	}

	parsedStoreIDs, parsedStores, err := responseParser.ParseLocationResponseMessage(locationsResponseData)

	if err != nil {
		fmt.Printf("Error parsing locations data: %v\n", err)
		return
	}

	fmt.Printf("Parsed Store IDs: %+v\n", parsedStoreIDs)
	fmt.Printf("Parsed Stores: %+v\n", parsedStores)
	fmt.Printf("------------------------------\n")

	//Get Availability Info
	availabilityUrl, err := reqUrlBuilder.GetAvailabilityUrl([]string{"19869815"}, "V8P", parsedStoreIDs)

	if err != nil {
		fmt.Printf("Error getting availability URL: %v\n", err)
		return
	}

	availabilityReqParam := fetcher.RequestParameters{
		RequestUrl:     availabilityUrl,
		RequestContext: fetcherContext,
	}
	availabilityResponseData, _, err := reqBuilder.FetchWithRequestParam(fetcherClient, availabilityReqParam)

	if err != nil {
		fmt.Printf("Error fetching availability data: %v\n", err)
		return
	}

	parsedAvailability, err := responseParser.ParseAvailabilityResponseMessage(availabilityResponseData)

	if err != nil {
		fmt.Printf("Error parsing availability data: %v\n", err)
		return
	}

	fmt.Printf("Parsed Availability: %+v\n", parsedAvailability)
	//fmt.Printf("Parsed Availability Status: %s\n", parsedAvailability["19869815"].PickupLocationAvailability.Status)
	fmt.Printf("------------------------------\n")
}
