package bestbuy

import (
	"encoding/json"
	"errors"

	"siteFetcher/internal/fetcher"
)

// BestBuyResponseParser handles BestBuy JSON parsing into shared product models.
type BestBuyResponseParser struct{}

func (b *BestBuyResponseParser) ParseResponseMessage(respBody []byte) (fetcher.Product, error) {
	if len(respBody) == 0 {
		return fetcher.Product{}, errors.New("Error: Response body is empty.")
	}

	var p fetcher.Product
	if err := json.Unmarshal(respBody, &p); err != nil {
		return fetcher.Product{}, err
	}
	return p, nil
}

// Location parsing types.
type Store struct {
	StoreID      string `json:"locationID"`
	StoreName    string `json:"name"`
	StoreAddress string `json:"address1"`
}

type Location struct {
	Stores []Store `json:"locations"`
}

func (b *BestBuyResponseParser) ParseLocationResponseMessage(respBody []byte) ([]string, map[string]Store, error) {
	if len(respBody) == 0 {
		return nil, nil, errors.New("Error: Response body is empty.")
	}

	var location Location
	if err := json.Unmarshal(respBody, &location); err != nil {
		return nil, nil, err
	}

	stores := location.Stores
	storeIds := make([]string, len(stores))
	storeMap := make(map[string]Store, len(stores))
	for i, store := range stores {
		storeIds[i] = store.StoreID
		storeMap[store.StoreID] = store
	}
	return storeIds, storeMap, nil
}

// Availability parsing types.
type PickupLocationvAvailabilityDetail struct {
	StoreName           string `json:"name"`
	StoreID             string `json:"locationKey"`
	QuantityOnHand      int    `json:"quantityOnHand"`
	HasInventory        bool   `json:"hasInventory"`
	SupportsFulfillment bool   `json:"supportsFulfillment"`
	IsReservable        bool   `json:"isReservable"`
}

type PickupLocationAvailability struct {
	Status                            string                              `json:"status"`
	PickupLocationAvailabilityDetails []PickupLocationvAvailabilityDetail `json:"locations"`
}

type ShippingAvailabilityDetail struct {
	CarrierName           string `json:"carrierName"`
	DeliveryDate          string `json:"deliveryDate"`
	DeliveryDateExpiresOn string `json:"deliveryDateExpiresOn"`
	DeliveryDatePrecision string `json:"deliveryDatePrecision"`
}

type ShippingAvailability struct {
	Status                      string                       `json:"status"`
	QuantityRemaining           int                          `json:"quantityRemaining"`
	OrderLimit                  int                          `json:"orderLimit"`
	IsBackorderable             bool                         `json:"isBackorderable"`
	ShippingAvailabilityDetails []ShippingAvailabilityDetail `json:"levelsOfServices"`
}

type ProductAvailability struct {
	SKU                        string                     `json:"sku"`
	SellerID                   string                     `json:"sellerId"`
	SaleChannelExclusivity     string                     `json:"saleChannelExclusivity"`
	PickupLocationAvailability PickupLocationAvailability `json:"pickup"`
	ShippingAvailability       ShippingAvailability       `json:"shipping"`
}

type ProductAvailabilitiesResponse struct {
	ProductAvailabilities []ProductAvailability `json:"availabilities"`
}

func (b *BestBuyResponseParser) ParseAvailabilityResponseMessage(respBody []byte) (map[string]ProductAvailability, error) {
	if len(respBody) == 0 {
		return nil, errors.New("Error: Response body is empty.")
	}

	var response ProductAvailabilitiesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	availabilityMap := make(map[string]ProductAvailability, len(response.ProductAvailabilities))
	for _, availability := range response.ProductAvailabilities {
		availabilityMap[availability.SKU] = availability
	}
	return availabilityMap, nil
}
