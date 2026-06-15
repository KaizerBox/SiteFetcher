package bestbuy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"siteFetcher/internal/fetcher"
)

// BestBuyFetcher ties together the URL builder, request builder, and parser to satisfy fetcher.Store.
type BestBuyFetcher struct {
	RequestUrlBuilder *BestBuyRequestUrlBuilder
	RequestBuilder    fetcher.RequestBuilder
	ResponseParser    *BestBuyResponseParser
}

func NewBestBuyFetcher(requestUrlBuilder *BestBuyRequestUrlBuilder, requestBuilder fetcher.RequestBuilder, responseParser *BestBuyResponseParser) (*BestBuyFetcher, error) {
	if requestUrlBuilder == nil || requestBuilder == nil || responseParser == nil {
		return nil, errors.New("nil parameter provided")
	}
	return &BestBuyFetcher{RequestUrlBuilder: requestUrlBuilder, RequestBuilder: requestBuilder, ResponseParser: responseParser}, nil
}

func (b *BestBuyFetcher) Name() string { return "bestbuy" }

func (b *BestBuyFetcher) GetResponse(productNum string, postalCodeFSA string) (fetcher.Product, time.Time, error) {
	productUrl, err := b.RequestUrlBuilder.GetProductUrl(productNum)
	if err != nil {
		return fetcher.Product{}, time.Time{}, err
	}

	requestParams := fetcher.RequestParameters{}.NewRequestParameters(context.Background(), productUrl, fetcher.RequestParameters{}.GetDefaultRequestHeaders())
	respBody, reqStartTime, err := b.RequestBuilder.FetchWithRequestParam(requestParams)
	if err != nil {
		return fetcher.Product{}, reqStartTime, err
	}

	locationUrl, err := b.RequestUrlBuilder.GetLocationsUrl(postalCodeFSA)
	if err != nil {
		return fetcher.Product{}, time.Time{}, err
	}

	requestParams = fetcher.RequestParameters{}.NewRequestParameters(context.Background(), locationUrl, fetcher.RequestParameters{}.GetDefaultRequestHeaders())
	storeRespBody, _, err := b.RequestBuilder.FetchWithRequestParam(requestParams)
	if err != nil {
		return fetcher.Product{}, reqStartTime, err
	}

	storeIds, storesMap, err := b.ResponseParser.ParseLocationResponseMessage(storeRespBody)
	if err != nil {
		return fetcher.Product{}, reqStartTime, err
	}

	availabilityUrl, err := b.RequestUrlBuilder.GetAvailabilityUrl([]string{productNum}, postalCodeFSA, storeIds)
	if err != nil {
		return fetcher.Product{}, time.Time{}, err
	}

	requestParams = fetcher.RequestParameters{}.NewRequestParameters(context.Background(), availabilityUrl, fetcher.RequestParameters{}.GetDefaultRequestHeaders())
	availabilityRespBody, _, err := b.RequestBuilder.FetchWithRequestParam(requestParams)
	if err != nil {
		return fetcher.Product{}, reqStartTime, err
	}

	availabilityInfo, err := b.ResponseParser.ParseAvailabilityResponseMessage(availabilityRespBody)
	if err != nil {
		return fetcher.Product{}, reqStartTime, err
	}

	productAvailabilityInfo, ok := availabilityInfo[productNum]
	if !ok {
		return fetcher.Product{}, reqStartTime, errors.New("Error: Product availability information not found.")
	}

	product, err := b.ResponseParser.ParseResponseMessage(respBody)
	if err != nil {
		return fetcher.Product{}, reqStartTime, err
	}

	inStoreAvailability := make([]fetcher.InStoreAvailability, len(productAvailabilityInfo.PickupLocationAvailability.PickupLocationAvailabilityDetails))
	for i, location := range productAvailabilityInfo.PickupLocationAvailability.PickupLocationAvailabilityDetails {
		storeInfo, ok := storesMap[location.StoreID]
		if !ok {
			return fetcher.Product{}, reqStartTime, fmt.Errorf("Error: Store information not found for store ID: %s", location.StoreID)
		}
		inStoreAvailability[i] = fetcher.InStoreAvailability{
			IsAvailableInStore:  location.HasInventory,
			StoreName:           storeInfo.StoreName,
			StoreID:             location.StoreID,
			StoreAddress:        storeInfo.StoreAddress,
			QuantityOnHand:      location.QuantityOnHand,
			HasInventory:        location.HasInventory,
			SupportsFulfillment: location.SupportsFulfillment,
			IsReservable:        location.IsReservable,
		}
	}
	product.InStoreAvailability = inStoreAvailability

	if len(productAvailabilityInfo.ShippingAvailability.ShippingAvailabilityDetails) == 0 {
		return fetcher.Product{}, reqStartTime, errors.New("Error: Shipping availability information is nil.")
	}

	product.OnlineAvailability = fetcher.OnlineAvailability{
		IsAvailableOnline:       product.OnlineAvailability.IsAvailableOnline,
		OnlineAvailability:      productAvailabilityInfo.ShippingAvailability.Status,
		OnlineAvailabilityCount: productAvailabilityInfo.ShippingAvailability.QuantityRemaining,
		EstimatedDeveliveryDate: productAvailabilityInfo.ShippingAvailability.ShippingAvailabilityDetails[0].DeliveryDate,
		DeliveryDateExpiresOn:   productAvailabilityInfo.ShippingAvailability.ShippingAvailabilityDetails[0].DeliveryDateExpiresOn,
		DeliveryDatePrecision:   productAvailabilityInfo.ShippingAvailability.ShippingAvailabilityDetails[0].DeliveryDatePrecision,
		IsBackOrderable:         productAvailabilityInfo.ShippingAvailability.IsBackorderable,
	}

	return product, reqStartTime, nil
}

func NewBestBuyStore(client *http.Client) (fetcher.Store, error) {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	availabilityQuery, err := LoadAvailabilityQueryConfig("internal/fetcher/bestbuy/availability_query_config.json")
	if err != nil {
		return nil, err
	}

	urlBuilder, err := NewBestBuyRequestUrlBuilderWithConfig("BestBuy", "https", "www.bestbuy.ca", "/api/v2/json/product/%s", "/api/v3/json/locations", "/ecomm-api/availability/products", availabilityQuery)
	if err != nil {
		return nil, err
	}

	parser := &BestBuyResponseParser{}
	requestBuilder := &fetcher.DefaultRequestBuilder{Client: client}
	return NewBestBuyFetcher(urlBuilder, requestBuilder, parser)
}
