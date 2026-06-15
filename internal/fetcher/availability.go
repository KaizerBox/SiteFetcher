package fetcher

// OnlineAvailability captures the online fulfillment details for a product.
type OnlineAvailability struct {
	IsAvailableOnline       bool   `json:"isAvailableOnline"`
	OnlineAvailability      string `json:"onlineAvailability"`
	OnlineAvailabilityCount int    `json:"onlineAvailabilityCount"`
	EstimatedDeveliveryDate string `json:"estimatedDeliveryDate"`
	DeliveryDateExpiresOn   string `json:"deliveryDateExpiresOn"`
	DeliveryDatePrecision   string `json:"deliveryDatePrecision"`
	IsBackOrderable         bool   `json:"isBackOrderable"`
}

// InStoreAvailability captures the in-store pickup details for a product.
type InStoreAvailability struct {
	IsAvailableInStore  bool   `json:"isAvailableInStore"`
	StoreName           string `json:"storeName"`
	StoreID             string `json:"storeId"`
	StoreAddress        string `json:"storeAddress"`
	QuantityOnHand      int    `json:"quantityOnHand"`
	HasInventory        bool   `json:"hasInventory"`
	SupportsFulfillment bool   `json:"supportsFulfillment"`
	IsReservable        bool   `json:"isReservable"`
}
