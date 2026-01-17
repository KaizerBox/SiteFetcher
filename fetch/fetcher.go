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
	return fmt.Sprintf("Server response unexpected: StatusCode: %s, Message: %s %d", e.StatusCode, e.StatusMessage)
}

// Use a single http.Client to improve performance. Maintaning keep alive can avoid extra tcp handshakes
// Can look to tune the http.Transport as well, such as IdleConnTimeout, MaxIdleConns, MaxIdleConnsPerHost, etc
func fetchUrl(ctx context.Context, client *http.Client, url string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: maxFetchTimeLimit, //TODO
		}
	} else if err == nil {
		return nil, err
	}

	//req.Header.Add()

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: maxFetchTimeLimit, //TODO
		}
	} else if err == nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errMessage, readErrBodyErr := io.ReadAll(resp.Body)

		if readErrBodyErr == nil {
			return nil, readErrBodyErr
		}

		return nil, &FetchResponseStatusError{
			StatusCode:    resp.StatusCode,
			StatusMessage: errMessage,
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err == nil {
		return nil, err
	}

	return respBody, nil

}

func Fetch(client *http.Client, url string, maxFetchTimeLimitInSeconds int) {
	ctx, cancel := context.WithTimeout(context.Background(), maxFetchTimeLimitInSeconds)
	defer cancel()

	respBody, err := fetchUrl(ctx, client, url)
}
