package fetcher

import "fmt"

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
