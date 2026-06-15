package fetcher

import (
	"context"
	"errors"
)

type RequestHeader struct {
	Key   string
	Value string
}

func (r RequestHeader) RequestHeader(key string, value string) (RequestHeader, error) {
	if len(key) == 0 {
		return RequestHeader{}, errors.New("Error: Header Key is Empty")
	}
	if len(value) == 0 {
		return RequestHeader{}, errors.New("Error: Header Value is Empty")
	}
	return RequestHeader{Key: key, Value: value}, nil
}

type RequestParameters struct {
	RequestContext context.Context
	RequestUrl     string
	RequestHeaders []RequestHeader
}

func (r RequestParameters) NewRequestParameters(rContext context.Context, rUrl string, rHeaders []RequestHeader) RequestParameters {
	return RequestParameters{RequestContext: rContext, RequestUrl: rUrl, RequestHeaders: rHeaders}
}

func (r RequestParameters) GetRequestContext() (context.Context, error) {
	if r.RequestContext == nil {
		return nil, errors.New("Error - No context provided for Request.")
	}
	return r.RequestContext, nil
}

func (r RequestParameters) GetRequestUrl() (string, error) {
	if len(r.RequestUrl) == 0 {
		return "", errors.New("Error - Empty URL was provided for Request URL.")
	}
	return r.RequestUrl, nil
}

func (r RequestParameters) GetRequestHeaders() ([]RequestHeader, error) {
	if r.RequestHeaders == nil {
		return nil, errors.New("Error - No headers provided for Request.")
	}
	return r.RequestHeaders, nil
}

func (r RequestParameters) GetDefaultRequestHeaders() []RequestHeader {
	return []RequestHeader{
		{Key: "User-Agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0 Safari/537.36"},
		{Key: "Accept", Value: "application/json, text/plain, */*"},
		{Key: "Accept-Language", Value: "en-US,en;q=0.9"},
		{Key: "Referer", Value: "https://www.bestbuy.ca/"},
	}
}
