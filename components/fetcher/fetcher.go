// You can edit this code!
// Click here and start typing.
package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type FetcherSetting struct {
	MaxFetchTimeLimit int
}

type Fetcher struct {
	FetchContext context.Context
	Client       *http.Client
	FetcherSetting
}

type FetchTimeOutError struct {
	FailedMaxFetchTimeLimit int
}

func (e *FetchTimeOutError) Error() string {
	return fmt.Sprintf("Fetch exceeded max time limit: %d", e.FailedMaxFetchTimeLimit)
}

type FetchResponseStatusError struct {
	StatusCode    int
	StatusMessage []byte
}

func (e *FetchResponseStatusError) Error() string {
	return fmt.Sprintf("Server response unexpected: StatusCode: %d, Message: %s", e.StatusCode, e.StatusMessage)
}

// Use a single http.Client to improve performance. Maintaning keep alive can avoid extra tcp handshakes
// Can look to tune the http.Transport as well, such as IdleConnTimeout, MaxIdleConns, MaxIdleConnsPerHost, etc
func (f *Fetcher) FetchWithUrl(url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(f.FetchContext, http.MethodGet, url, nil)

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: f.MaxFetchTimeLimit, //TODO
		}
	} else if err != nil {
		return nil, err
	}

	//req.Header.Add()

	resp, err := f.Client.Do(req)

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: f.MaxFetchTimeLimit, //TODO
		}
	} else if err != nil {
		return nil, err
	}
	//defer after checking for error to avoid Body.Close() error which may occur when Do fails
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errMessage, readErrBodyErr := io.ReadAll(resp.Body)

		if readErrBodyErr != nil {
			return nil, readErrBodyErr
		}

		return nil, &FetchResponseStatusError{
			StatusCode:    resp.StatusCode,
			StatusMessage: errMessage,
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil

}
