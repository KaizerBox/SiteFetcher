package fetcher

import (
	"net/http"
	"time"
)

type Fetcher struct {
	Client            *http.Client
	RequestUrlBuilder RequestUrlBuilder
	RequestBuilder    RequestBuilder
	ResponseParser    ResponseParser
	ResponseFilter    ResponseFilter
}

type ParseResponseBodyResult struct {
	ItemName                string
	OriginalPrice           string
	DiscountAmount          string
	PriceAfterDiscount      string
	Stock                   string
	Seller                  string
	QueryExecutionTimestamp time.Time
}

func (f *Fetcher) ParseResponseMessage(reponseBody []byte) (*ParseResponseBodyResult, error) {
	return nil, nil
}

type ItemQueryCondition struct {
	QueryURL           string
	PriceLowerBound    float32
	PriceUpperBound    float32
	InStockStatus      bool
	StockLowerBound    int
	StockUpperBound    int
	DiscountStatus     bool
	DiscountLowerBound float32
	DiscountUpperBound float32
}

func (i *ItemQueryCondition) ItemQueryStatus() (ItemQueryCondition, error) {
	return *i, nil
}
