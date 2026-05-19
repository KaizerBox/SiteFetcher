// You can edit this code!
// Click here and start typing.
package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)


type Fetcher struct {
	Client *http.Client //Overall Client shared among the sites, manage content, share connections etc.
	RequestBuilder *RequestBuilder //Per Site, each site implements this to customize header, auth, query params, etc Sents the query, and returns the Response Result or error
	ReponseParser *ReponseParser //Per Site, each site implements this to parse the Response Result from RequestBuilder's call to website. 
	                             //Uses the response values to determine if the condition for this item is satisfied, if it is then call a function on the Reponse Result to generate a email template 
								 //struct that can be passed to gmailer to send.
	ReponseFilter *ResFilter 

type ReqBuilder interface {
	//Makes request using URL and context provided in RequestParameters struct
	FetchWithRequestParam(RequestParameters) (byte[], time.Time, error)
}


type ResParser interface {
	//Parses the response body from Fetch (byte[]) and find all if it satisfies result of the query condition ItemQueryCondition. Gives ParseResponseBodyResponse if it does not error out.
	ParseResponseMessage(byte[]) (ParseResponseBodyResult, error)
}

//Will be Implemented by ItemQueryCondition?
type ResFilter interface {
	PriceFilter(ParseResponseBodyResult) (bool, error)
	StockFilter(ParseResponseBodyResult) (bool, error)
	DiscountFilter(ParseResponseBodyResult) (bool, error)
	FullFilter(ParseResponseBodyResponse) (bool, error)
	GenerateEmailMessage(ParseResponseBodyResponse) (GmailMessage, error)
	FilterAndGenerateEmailMessage(ParseResponseBodyResponse) (GmailMessage, error)
}
	
type FetcherSetting struct {
	MaxFetchTimeLimit int
	FetcherSiteName string
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

type RequestParameters struct {
	requestContext context.Context
	requestUrl string
}

func (r *RequestParameters) RequestParameters(rContext, rUrl) RequestParameters {
	return RequestParameters{requestContext: rContext, requestUrl:rUrl}
}

func (r *requestParameters) GetRequestContext() (context.Context, error) {
	return r.requestContext
}

func (r *requestParameters) GetRequestUrl() (string, error) {
	if len(r.requestURL) == nil {
		return "", errors.New("Error - Empty URL was provided for Request URL.")
	}
}

// Use a single http.Client to improve performance. Maintaning keep alive can avoid extra tcp handshakes
// Can look to tune the http.Transport as well, such as IdleConnTimeout, MaxIdleConns, MaxIdleConnsPerHost, etc
func (f *Fetcher) FetchWithRequestParam(requestParams RequestParameters) ([]byte, time.Time, error) {
	reqStartTime = time.Now()
	
	reqContext, err := requestParams.GetRequestContext()
	if err != nil {
		return nil, reqStartTime, err
	}

	reqUrl, err := requestParams.GetRequestUrl()
	if err != nil {
		return nil, reqStartTime, err
	}
	
	req, err := http.NewRequestWithContext(reqContext, http.MethodGet, reqUrl, nil)

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, reqStartTime, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: f.MaxFetchTimeLimit, //TODO
		}
	} else if err != nil {
		return nil, reqStartTime, err
	}

	//req.Header.Add()

	resp, err := f.Client.Do(req)

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: f.MaxFetchTimeLimit, //TODO
		}
	} else if err != nil {
		return nil, reqStartTime, err
	}
	//defer after checking for error to avoid Body.Close() error which may occur when Do fails
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errMessage, readErrBodyErr := io.ReadAll(resp.Body)

		if readErrBodyErr != nil {
			return nil, reqStartTime, readErrBodyErr
		}

		return nil, reqStartTime, &FetchResponseStatusError{
			StatusCode:    resp.StatusCode,
			StatusMessage: errMessage,
		}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, reqStartTime, err
	}

	return respBody, reqStartTime, nil

}

type ParseResponseBodyResult struct {
	ItemName string
	OriginalPrice string
	DiscountAmount string
	PriceAfterDiscount string
	Stock string
	QueryExecutionTimestamp time.Time
	QueryURL string
}

//func (p *ParseResponseBodyResult) GenerateEmailTemplate() {
// --TODO Generates the necessary string fields gmailer to send email to user when query condition of the ResponseBody as parsed sucessfully.
//}

//Parse the ResponseBody based on query conditions (i.e. price, stock, etc).
//Might need to implement a sub function for each condition whic reuslt in boolean, then ParseResponseMessage evaluates the overall condition based on the sub functions.
func (f *Fetcher) ParseResponseMessage(reponseBody []byte) (*ParseResponseBodyResult, error) {
	return nil, error
}

//Contains all the conditions for the item, each item have this struct. Site and URL then uniquely identifies each tracker.
type ItemQueryCondition struct {
	QueryURL string
	PriceLowerBound float32
	PriceUpperBound float32
	InStockStatus bool
	StockLowerBound int
	StockUpperBound int
	DiscountStatus bool
	DiscountLowerBound float32
	DiscountUpperBound float32
}

func (i *ItemQueryCondition) ItemQueryStatus() (ItemQueryConditions, error) {
	//creates tis based on json or other structured file format parsed for the file that stores this info
}

