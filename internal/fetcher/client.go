package fetcher

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type DefaultRequestBuilder struct {
	Client *http.Client
}

func (b *DefaultRequestBuilder) DefaultRequestBuilder(client *http.Client) (DefaultRequestBuilder, error) {
	if client == nil {
		return DefaultRequestBuilder{}, errors.New("Error: HTTP Client is nil.")
	}
	return DefaultRequestBuilder{Client: client}, nil
}

func (b *DefaultRequestBuilder) FetchWithRequestParam(requestParams RequestParameter) ([]byte, time.Time, error) {
	reqStartTime := time.Now()

	reqContext, err := requestParams.GetRequestContext()
	if err != nil {
		return nil, reqStartTime, err
	}

	reqUrl, err := requestParams.GetRequestUrl()
	if err != nil {
		return nil, reqStartTime, err
	}

	req, err := http.NewRequestWithContext(reqContext, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, reqStartTime, err
	}

	reqHeaders, err := requestParams.GetRequestHeaders()
	if err != nil {
		return nil, reqStartTime, err
	}

	for _, header := range reqHeaders {
		req.Header.Add(header.Key, header.Value)
	}

	client := b.Client
	if client == nil {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, reqStartTime, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, reqStartTime, readErr
		}
		return nil, reqStartTime, &FetchResponseStatusError{
			StatusCode:    resp.StatusCode,
			StatusMessage: errBody,
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, reqStartTime, err
	}

	return respBody, reqStartTime, nil
}
