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
	"strconv"
)

/*
From Andriod App Config: https://www.bestbuy.ca/ns/mobile/android/configuration/v1.5.0/en
"category-api-url": "http://www.bestbuy.ca/api/v2/json/category/",
"product-api-url": "http://www.bestbuy.ca/api/v2/json/product/",
"search-api-url": "http://www.bestbuy.ca/api/v2/json/search",
"store-api-url": "http://www.bestbuy.ca/api/v2/json/locations",
"pdp-api-url": "http://www.bestbuy.ca/api/v2/json/product/",
"reviews-api-url": "http://www.bestbuy.ca/api/v2/json/reviews/",
"availability-api-url": "http://www.bestbuy.ca/api/v2/json/availability", //getProductAvailabilityForSingleSku.sku
"publication-api-url": "http://www.bestbuy.ca/api/v2/json/publication/",
"media-base-url": "https://multimedia.bbycastatic.ca",
"product-availability-api-url":"http://api.bestbuy.ca/availability/products",

From website Search
"dataSources": {
                        "offerApiUrl": "https://www.bestbuy.ca/api/offers/v1/products",
                        "productDomainApi": "https://www.bestbuy.ca/api/v3/products",
                        "accountApiUrl": "https://www.bestbuy.ca/api/account",
                        "basketServiceApiUrl": "https://www.bestbuy.ca/api/basket",
                        "brandsApiUrl": "https://www.bestbuy.ca/api/v2/json/brands/",
                        "cellPhonesCarrierPlansUrl": "https://www.bestbuy.ca/api/cellphones-plans-pricing/sku/{skuId}?api-version=2022-05-01&sp=%2Ftriggers%2Fmanual%2Frun&sv=1.0&sig=qqpzPnL_WPQXWUV73BbXlPLU0EGP_ZfI0vsIJFccWOE",
                        "locationApiUrl": "https://www.bestbuy.ca/api/v3/json/locations/locate",
                        "availabilityApiUrl": "https://www.bestbuy.ca/ecomm-api/availability/products",
                        "productApiUrl": "https://www.bestbuy.ca/api/v2/json/product/",
                        "presentationCatalogQueryApiUrl": "https://www.bestbuy.ca/api/v1/catalog/query",
                        "productGatewayApiUrl": "https://www.bestbuy.ca/api/v3/",
                        "baseSwatchUrl": "https://multimedia.bbycastatic.ca",
                        "categoryApiUrl": "https://www.bestbuy.ca/api/v2/json/category/",
                        "contentApiUrl": "https://www.bestbuy.ca/api/merch/v1/",
                        "contentFallbackApiUrl": {
                            "en": "https://ecomm-media.bbycastatic.ca/digital-assets/static/homepage/en.json",
                            "fr": "https://ecomm-media.bbycastatic.ca/digital-assets/static/homepage/fr.json"
                        },
                        "collectionApiUrl": "https://www.bestbuy.ca/api/merch/v2/",
                        "priceEhfReadMoreUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/paying-and-purchasing/environmental-handling-fees",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/paiement-et-achat/ecofrais"
                        },
                        "pdpFinancingReadMoreUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/about/best-buy-financing/blt12a5f527f63859b7?icmp=pdp_financing_how_it_works",
                            "fr": "https://www.bestbuy.ca/fr-ca/a-propos/financement-best-buy/blt12a5f527f63859b7?icmp=pdp_financing_how_it_works"
                        },
                        "specialOfferApiUrl": "https://www.bestbuy.ca/api/soc/v1/products",
                        "helpTopicsApiUrl": "https://www.bestbuy.ca/api/merch/v1/help",
                        "reviewApiUrl": "https://www.bestbuy.ca/api/reviews/v2",
                        "searchApiUrl": "https://www.bestbuy.ca/api/v2/json/search",
                        "presentationSubscriptionApiUrl": "https://www.bestbuy.ca/api/presentation-subscription",
                        "sellerApiUrl": "https://www.bestbuy.ca/api/seller/v1/sellers/",
                        "sellerReviewsApiUrl": "https://www.bestbuy.ca/api/v2/json/sellerreviews",
                        "contentEnApiUrl": "https://blog.bestbuy.ca/wp-json/wp/v2/pages/",
                        "contentFrApiUrl": "https://blogue.bestbuy.ca/wp-json/wp/v2/pages/",
                        "categorySeoEnUrl": "https://ecomm-media.bbycastatic.ca/digital-assets/category-seo/en.json",
                        "categorySeoFrUrl": "https://ecomm-media.bbycastatic.ca/digital-assets/category-seo/fr.json",
                        "digitalAssetsMaxCdnUrl": "https://ecomm-media.bbycastatic.ca/digital-assets/",
                        "storeLocationApiUrl": "https://www.bestbuy.ca/api/v3/json/locations",
                        "remoteConfigUrl": "https://www.bestbuy.ca/remote-config/config.json",
                        "navigationLinkGroupsBaseUrl": "https://ecomm-media.bbycastatic.ca/digital-assets/navigation-link-groups/",
                        "salesforceWebToLeadURL": "https://webto.salesforce.com/servlet/servlet.WebToLead?encoding=UTF-8",
                        "inHomeAdvisorTermsAndConditionsUrl": {
                            "en": "https://www.bestbuy.ca/projects/inhomeadvisor/assets/terms_and_conditions.pdf",
                            "fr": "https://www.bestbuy.ca/projects/inhomeadvisor/assets/terms_and_conditions.pdf"
                        },
                        "inHomeAdvisorPrivacyPolicyUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/blt372c78db41358a01/blt12691f41ac6895d7",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/blt372c78db41358a01/blt12691f41ac6895d7"
                        },
                        "warrantyTermsAndConditionsUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/blt372c78db41358a01/blt612eaea73f4477ad",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/blt372c78db41358a01/blt612eaea73f4477ad"
                        },
                        "interestBasedAdsTermsAndConditionsUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/policies-and-terms-and-conditions/interest-based-ads",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/politiques-modalites-et-conditions/publicites-ciblees-par-centres-d-interet"
                        },
                        "quickAndEasyPickupHelpUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/shipping-delivery-and-pick-up/quick-and-easy-store-pickup",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/expedition-livraison-et-ramassage/ramassage-rapide-et-facile"
                        },
                        "quebecLegalWarrantyUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/blt372c78db41358a01/blt9829bb6fbacda8f2",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/blt372c78db41358a01/blt9829bb6fbacda8f2"
                        },
                        "myBestBuyAccountHelpUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/your-account",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/votre-compte"
                        },
                        "shippingAndDeliveryHelpUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/shipping-delivery-and-pick-up/common-shipping-questions",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/expedition-livraison-et-ramassage/questions-courantes-sur-l-expedition"
                        },
                        "returnsAndExchangesHelpUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/returns-and-exchanges/returning-or-exchanging-a-best-buy-product",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/politique-de-retours-et-echanges/retourner-ou-echanger-un-produit-best-buy"
                        },
                        "findOrderUrl": {
                            "en": "https://www.bestbuy.ca/order/en-ca/find-order",
                            "fr": "https://www.bestbuy.ca/order/fr-ca/trouver-commande"
                        },
                        "contactUsUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/help/your-order/contact-us",
                            "fr": "https://www.bestbuy.ca/fr-ca/aide/vos-commandes/contactez-nous"
                        },
                        "salesforceWebToCaseUrl": "https://webto.salesforce.com/servlet/servlet.WebToCase?encoding=UTF-8",
                        "storeLocatorUrl": "https://stores.bestbuy.ca/{locale}/{locationId}",
                        "recommendationsUrl": "https://bestbuycanada.tt.omtrdc.net/m2/bestbuycanada/ubox/raw",
                        "accountDashboardUrl": "https://www.bestbuy.ca/account",
                        "newsLetterApiUrl": "https://www.bestbuy.ca/api/consent-management",
                        "bbProtectionChatUrl": "https://remotesupport.bestbuy.ca/js/chat-entry.js",
                        "bbProtectionCovertOrigin": "https://remotesupport.bestbuy.ca",
                        "communicationApiUrl": "https://www.bestbuy.ca/api/conversations",
                        "storeMessageApiUrl": "https://www.bestbuy.ca/store-messages",
                        "storesStatusApiUrl": "https://www.bestbuy.ca/store-status",
                        "productListApiUrl": "https://www.bestbuy.ca/api/product-list",
                        "relatedProductsApiUrl": "https://www.bestbuy.ca/api/v3/products",
                        "googleScriptUrl": "https://securepubads.g.doubleclick.net/tag/js/gpt.js",
                        "googleAdSenseScriptUrl": "https://www.google.com/adsense/search/ads.js",
                        "adobeScriptUrl": "https://assets.adobedtm.com/launch-EN048473b7b22f4c66a2858404e3c8219d.min.js",
                        "financePlanUrl": "https://www.bestbuy.ca/api/financingplanapi/v1/financingPlans/web",
                        "notificationExperienceApiUrl": "https://www.bestbuy.ca/api/nea"
                    },
                    "staticUrls": {
                        "bbyHealthUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/services/best-buy-health/blt4265f1a522b6ec4e?icmp=ipp_bestbuyhealth_contactform_link_bbyhealth",
                            "fr": "https://www.bestbuy.ca/fr-ca/services/services-de-sante-best-buy/blt4265f1a522b6ec4e?icmp=ipp_bestbuyhealth_contactform_link_bbyhealth"
                        },
                        "assuredLivingUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/services/best-buy-health-assured-living/blt9bee2e53cafb82cd?icmp=ipp_bestbuyhealth_contactform_link_assuredliving",
                            "fr": "https://www.bestbuy.ca/fr-ca/services/services-de-sante-best-buy-vivez-rassures/blt9bee2e53cafb82cd?icmp=ipp_bestbuyhealth_contactform_link_assuredliving"
                        },
                        "bbyHealthBlogUrl": {
                            "en": "https://blog.bestbuy.ca/tag/best-buy-health?icmp=ipp_bestbuyhealth_contactform_link_blogs",
                            "fr": "https://blog.bestbuy.ca/tag/best-buy-health?icmp=ipp_bestbuyhealth_contactform_link_blogs"
                        },
                        "bbyBusinessUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/about/best-buy-business/bltfad9143fefc09dc6",
                            "fr": "https://www.bestbuy.ca/fr-ca/a-propos/best-buy-affaires/bltfad9143fefc09dc6"
                        },
                        "bbyBusinessWorkFromHomeUrl": {
                            "en": "https://www.bestbuy.ca/en-ca/about/best-buy-business-work-from-home/blt42e5f92ecd5dea96?icmp=ipp-bbb-wfh",
                            "fr": "https://www.bestbuy.ca/fr-ca/a-propos/best-buy-affaires-teletravail/blt42e5f92ecd5dea96"
                        },
                        "bestbuyLogoUrl": "https://www.bestbuy.ca/ns-static/img/bestbuy-canada-logo.jpg"
                    },
*/


type Fetcher struct {
	Client *http.Client //Overall Client shared among the sites, manage content, share connections etc.
	RequestUrlBuilder *RequestUrlBuilder
	RequestBuilder *RequestBuilder //Per Site, each site implements this to customize header, auth, query params, etc Sents the query, and returns the Response Result or error
	ReponseParser *ReponseParser //Per Site, each site implements this to parse the Response Result from RequestBuilder's call to website. 
	                             //Uses the response values to determine if the condition for this item is satisfied, if it is then call a function on the Reponse Result to generate a email template 
								 //struct that can be passed to gmailer to send.
	ReponseFilter *ResponseFilter 
}

type RequestUrlBuilder interface {
	GetProductUrl(string) (string, error)
	//GetAvailabilityUrl(string) (string, error)
}

type ReqBuilder interface {
	//Makes request using URL and context provided in ReqParameter interface
	FetchWithRequestParam(*http.Client, ReqParameter) (byte[], time.Time, error)
}

type ResponseParser interface {
	//Parses the response body from Fetch (byte[]) and find all if it satisfies result of the query condition ItemQueryCondition. Gives ParseResponseBodyResponse if it does not error out.
	ParseResponseMessage(byte[]) (ParseResponseBodyResult, error)
}

//Will be Implemented by ItemQueryCondition?
type ResponseFilter interface {
	Filter(ParseResponseBodyResult, ItemQueryCondition) (bool, error)
}

type ReqParameter interface {
	GetRequestContext() (context.Context, error)

	GetRequestUrl() (string, error)
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

type BestBuyRequestUrlBuilder struct {
	Name string = "BestBuy"
	//API Url link GET
	DomainURL string = "https://www.bestbuy.ca/api/v2/json/product/"
}

func (b *BestBuyReqBuilder) BestBuyReqBuilder(name string, domainUrl string) (BestBuyReqBuilder, error) {
	if len(name) == 0 {
		return nil, errors.New("Error: name is Empty")
	}
	if len(domainUrl) == 0 {
		return nil, errors.New("Error: Domain URL is Empty")
	}

	return BestBuyReqBuilder{Name: name, DomainURL:domainUrl,}, nil
}

func (b *BestBuyReqBuilder) GetProductUrl(productNum string) (string, error) {
	val, err := strconv.Atoi(productNum)
	if len(productNum) == 0 || err != nil {
		return "", fmt.Errorf("Error: Product Number: %s is invalid.", productNum)
	}
	return b.DomainURL + productNum, nil
}

type BestBuyRequestBuilder struct {
	Client *http.Client
}

// Use a single http.Client to improve performance. Maintaning keep alive can avoid extra tcp handshakes
// Can look to tune the http.Transport as well, such as IdleConnTimeout, MaxIdleConns, MaxIdleConnsPerHost, etc
func (b *BestBuyReqBuilder) FetchWithRequestParam(client *http.Client, requestParams ReqParameter) ([]byte, time.Time, error) {
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

	resp, err := client.Do(req)

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

type BestBuyResponseParser {

}

func (b * BestBuyResponseParser) ParseResponseMessage(byte[]) (ParseResponseBodyResult, error) {
	
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
