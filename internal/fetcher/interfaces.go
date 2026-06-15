package fetcher

import (
	"context"
	"time"
)

type Store interface {
	Name() string
	GetResponse(productNum string, postalCodeFSA string) (Product, time.Time, error)
}

type RequestUrlBuilder interface {
	GetProductUrl(string) (string, error)
}

type RequestBuilder interface {
	FetchWithRequestParam(RequestParameter) ([]byte, time.Time, error)
}

type ResponseParser interface {
	ParseResponseMessage([]byte) (Product, error)
}

type ResponseFilter interface {
	Filter(ParseResponseBodyResult, ItemQueryCondition) (bool, error)
}

type RequestParameter interface {
	GetRequestContext() (context.Context, error)
	GetRequestUrl() (string, error)
	GetRequestHeaders() ([]RequestHeader, error)
}
