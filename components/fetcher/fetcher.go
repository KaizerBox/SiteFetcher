// You can edit this code!
// Click here and start typing.
package fetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	Client            *http.Client //Overall Client shared among the sites, manage content, share connections etc.
	RequestUrlBuilder *RequestUrlBuilder
	RequestBuilder    *RequestBuilder //Per Site, each site implements this to customize header, auth, query params, etc Sents the query, and returns the Response Result or error
	ResponseParser    *ResponseParser //Per Site, each site implements this to parse the Response Result from RequestBuilder's call to website.
	//Uses the response values to determine if the condition for this item is satisfied, if it is then call a function on the Reponse Result to generate a email template
	//struct that can be passed to gmailer to send.
	ResponseFilter *ResponseFilter
}

type RequestUrlBuilder interface {
	GetProductUrl(string) (string, error)
	//GetAvailabilityUrl(string) (string, error)
}

type RequestBuilder interface {
	//Makes request using URL and context provided in ReqParameter interface
	FetchWithRequestParam(*http.Client, RequestParameter) ([]byte, time.Time, error)
}

type ResponseParser interface {
	//Parses the response body from Fetch (byte[]) and find all if it satisfies result of the query condition ItemQueryCondition. Gives ParseResponseBodyResponse if it does not error out.
	ParseResponseMessage([]byte) (ParseResponseBodyResult, error)
}

// Will be Implemented by ItemQueryCondition?
type ResponseFilter interface {
	Filter(ParseResponseBodyResult, ItemQueryCondition) (bool, error)
}

type RequestParameter interface {
	GetRequestContext() (context.Context, error)

	GetRequestUrl() (string, error)
}

type RequestParameters struct {
	RequestContext context.Context
	RequestUrl     string
}

func (r RequestParameters) RequestParameters(rContext context.Context, rUrl string) RequestParameters {
	return RequestParameters{RequestContext: rContext, RequestUrl: rUrl}
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

type BestBuyRequestUrlBuilder struct {
	Name   string //"BestBuy",
	Scheme string //"https",
	Host   string //"www.bestbuy.ca",
	//API Url link GET
	ProductURL      string //"/api/v2/json/product/{productNum}",
	LocationUrl     string //"/api/v3/json/locations?lang=en-CA&postalCode={postalCodeFSA}",
	AvailabilityUrl string //"/ecomm-api/availability/products?sku={productNum}&storeIds={storeId}&postalCode={postalCodeFSA}",
}

func (b *BestBuyRequestUrlBuilder) BestBuyRequestUrlBuilder(name string, scheme string, host string, ProductURL string, LocationUrl string, AvailabilityUrl string) (BestBuyRequestUrlBuilder, error) {
	if len(name) == 0 {
		return BestBuyRequestUrlBuilder{}, errors.New("Error: name is Empty")
	}
	if len(scheme) == 0 {
		return BestBuyRequestUrlBuilder{}, errors.New("Error: Scheme is Empty")
	}
	if len(host) == 0 {
		return BestBuyRequestUrlBuilder{}, errors.New("Error: Host is Empty")
	}

	return BestBuyRequestUrlBuilder{Name: name, Scheme: scheme, Host: host, ProductURL: ProductURL, LocationUrl: LocationUrl, AvailabilityUrl: AvailabilityUrl}, nil
}

func (b *BestBuyRequestUrlBuilder) GetProductUrl(productNum string) (string, error) {
	_, err := strconv.Atoi(productNum)
	if len(productNum) == 0 || err != nil {
		return "", fmt.Errorf("Error: Product Number: %s is invalid.", productNum)
	}
	u := &url.URL{
		Scheme: b.Scheme,
		Host:   b.Host,
		Path:   fmt.Sprintf(b.ProductURL, productNum),
	}
	return u.String(), nil
}

func validateFSAConfigString(configStr string) error {
	// Validate one or more allowed characters or ranges inside square brackets.
	regexPattern := `^\[[A-Za-z0-9\-]+\]$`
	if !regexp.MustCompile(regexPattern).MatchString(configStr) {
		return fmt.Errorf("Error: Invalid FSA config string: %s", configStr)
	}
	return nil
}

func validatePostalCodeFSA(postalCodeFSA string) (bool, string, error) {
	if len(postalCodeFSA) == 0 {
		return false, "", errors.New("Error: Postal Code FSA is Empty")
	}
	if len(postalCodeFSA) != 3 {
		return false, "", fmt.Errorf("Error: Postal Code FSA: %s is invalid, should be 3 characters.", postalCodeFSA)
	}

	//Move to config parameters
	FSAAllowedFirstChar := "[ABCEGHJKLMNPRSTVXY]"
	FSAAllowedSecondChar := "[0-9]"
	FSAAllowedThirdChar := "[A-Z]"
	//Validate the config string for FSA is of type [A-Za-z0-9]]
	if err := validateFSAConfigString(FSAAllowedFirstChar); err != nil {
		return false, "", err
	}
	if err := validateFSAConfigString(FSAAllowedSecondChar); err != nil {
		return false, "", err
	}
	if err := validateFSAConfigString(FSAAllowedThirdChar); err != nil {
		return false, "", err
	}

	//Create regex to validate FSA based on configed valid FSA values.
	postalCodeFSARegex := regexp.MustCompile(fmt.Sprintf(`^%s%s%s$`, FSAAllowedFirstChar, FSAAllowedSecondChar, FSAAllowedThirdChar))
	//String to Upper in case of user input lower case, and check regex match for valid FSA format.
	result := postalCodeFSARegex.MatchString(strings.ToUpper(postalCodeFSA))

	return result, postalCodeFSARegex.String(), nil
}

func (b *BestBuyRequestUrlBuilder) GetLocationsUrl(postalCodeFSA string) (string, error) {
	valid, regexString, err := validatePostalCodeFSA(postalCodeFSA)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("Error: Postal Code FSA: %s is invalid, should match regex: %s", postalCodeFSA, regexString)
	}

	u := &url.URL{
		Scheme: b.Scheme,
		Host:   b.Host,
		Path:   b.LocationUrl,
	}
	q := u.Query()
	q.Set("lang", "en-CA")
	q.Set("postalCode", postalCodeFSA)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Only uses the FSA of the postal code. FSA is the first 3 characters of the postal code which identifies the general area of the location.
func (b *BestBuyRequestUrlBuilder) GetAvailabilityUrl(productNums []string, postalCodeFSA string, storeId []string) (string, error) {
	if len(productNums) == 0 {
		return "", errors.New("Error: Product Number list is Empty")
	}

	if valid, regexString, err := validatePostalCodeFSA(postalCodeFSA); err != nil {
		return "", err
	} else if !valid {
		return "", fmt.Errorf("Error: Postal Code FSA: %s is invalid, should match regex: %s", postalCodeFSA, regexString)
	}

	if len(storeId) == 0 {
		return "", errors.New("Error: Store ID list is Empty")
	}

	u := &url.URL{
		Scheme: b.Scheme,
		Host:   b.Host,
		Path:   b.AvailabilityUrl,
	}
	q := u.Query()
	q.Set("accept", "application/vnd.bestbuy.standardproduct.v1+json")
	q.Set("accept-language", "en-CA")
	q.Set("locations", strings.Join(storeId, "|"))
	q.Set("postalCode", postalCodeFSA)
	q.Set("skus", strings.Join(productNums, "|"))
	q.Encode()
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type BestBuyRequestBuilder struct {
	Client *http.Client
}

func (b *BestBuyRequestBuilder) BestBuyRequestBuilder(client *http.Client) (BestBuyRequestBuilder, error) {
	if client == nil {
		return BestBuyRequestBuilder{}, errors.New("Error: HTTP Client is nil.")
	}
	return BestBuyRequestBuilder{Client: client}, nil
}

// Use a single http.Client to improve performance. Maintaning keep alive can avoid extra tcp handshakes
// Can look to tune the http.Transport as well, such as IdleConnTimeout, MaxIdleConns, MaxIdleConnsPerHost, etc
func (b *BestBuyRequestBuilder) FetchWithRequestParam(client *http.Client, requestParams RequestParameter) ([]byte, time.Time, error) {
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

	/* if errors.Is(err, context.DeadlineExceeded) {
		return nil, reqStartTime, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: f.MaxFetchTimeLimit, //TODO
		}
	} else if err != nil {
		return nil, reqStartTime, err
	} */

	//req.Header.Add()

	// Add common browser headers to reduce chance the server rejects the request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-CA,en;q=0.9")
	req.Header.Set("Referer", "https://www.bestbuy.ca/")

	resp, err := client.Do(req)

	if err != nil {
		return nil, reqStartTime, err
	}

	/* if errors.Is(err, context.DeadlineExceeded) {
		return nil, &FetchTimeOutError{
			FailedMaxFetchTimeLimit: f.MaxFetchTimeLimit, //TODO
		}
	} else if err != nil {
		return nil, reqStartTime, err
	} */
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

type BestBuyResponseParser struct {
}

type Availability struct {
	SKU                     string `json:"sku"`
	InStoreAvailability     string `json:"inStoreAvailability"`
	IsAvailableOnline       bool   `json:"isAvailableOnline"`
	OnlineAvailability      string `json:"onlineAvailability"`
	OnlineAvailabilityCount int    `json:"onlineAvailabilityCount"`
}

type Product struct {
	Name         string       `json:"name"`
	SKU          string       `json:"sku"`
	RegularPrice float64      `json:"regularPrice"`
	SalePrice    float64      `json:"salePrice"`
	Seller       string       `json:"sellerId"`
	Availability Availability `json:"availability"`
}

func (b *BestBuyResponseParser) ParseResponseMessage(respBody []byte) (ParseResponseBodyResult, error) {
	if len(respBody) == 0 {
		return ParseResponseBodyResult{}, errors.New("Error: Response body is empty.")
	}
	//ToDo
	var p Product
	if err := json.Unmarshal(respBody, &p); err != nil {
		return ParseResponseBodyResult{}, err
	}
	parsedResponseResult := ParseResponseBodyResult{
		ItemName:                p.Name,
		OriginalPrice:           fmt.Sprintf("%.2f", p.RegularPrice),
		DiscountAmount:          fmt.Sprintf("%.2f", p.RegularPrice-p.SalePrice),
		PriceAfterDiscount:      fmt.Sprintf("%.2f", p.SalePrice),
		Stock:                   fmt.Sprintf("In Store: %s, Online: %s (%d available)", p.Availability.InStoreAvailability, p.Availability.OnlineAvailability, p.Availability.OnlineAvailabilityCount),
		QueryExecutionTimestamp: time.Now(),
		Seller:                  p.Seller,
	}
	return parsedResponseResult, nil
}

type Store struct {
	StoreID      string `json:"locationID"`
	StoreName    string `json:"name"`
	StoreAddress string `json:"address1"`
}

type Location struct {
	Stores []Store `json:"locations"`
}

func (b *BestBuyResponseParser) ParseLocationResponseMessage(respBody []byte) ([]string, []Store, error) {
	if len(respBody) == 0 {
		return nil, nil, errors.New("Error: Response body is empty.")
	}

	// Process locations to extract store IDs
	var location Location
	if err := json.Unmarshal(respBody, &location); err != nil {
		return nil, nil, err
	}
	stores := location.Stores
	storeIds := make([]string, len(stores))
	for i, store := range stores {
		storeIds[i] = store.StoreID
	}
	return storeIds, stores, nil
}

type PickupLocationvAvailabilityDetail struct {
	StoreName           string `json:"name"`
	StoreID             string `json:"locationKey"`
	QuantityOnHand      int    `json:"quantityOnHand"`
	HasInventory        bool   `json:"hasInventory"`
	SupportsFulfillment bool   `json:"supportsFulfillment"`
	IsReservable        bool   `json:"isReservable"`
}

type PickupLocationAvailability struct {
	Status                            string                              `json:"status"`
	PickupLocationAvailabilityDetails []PickupLocationvAvailabilityDetail `json:"locations"`
}

type ShippingAvailabilityDetail struct {
	CarrierName           string `json:"carrierName"`
	DeliveryDtae          string `json:"deliveryDate"`
	DeliveryDateExpiresOn string `json:"deliveryDateExpiresOn"`
	DeliveryDatePrecision string `json:"deliveryDatePrecision"`
}

type ShippingAvailability struct {
	Status                      string                       `json:"status"`
	QuantityRemaining           int                          `json:"quantityRemaining"`
	ShippingAvailabilityDetails []ShippingAvailabilityDetail `json:"levelsOfServices"`
}

type ProductAvailability struct {
	SKU                        string                     `json:"sku"`
	SellerID                   string                     `json:"sellerId"`
	SaleChannelExclusivity     string                     `json:"saleChannelExclusivity"`
	PickupLocationAvailability PickupLocationAvailability `json:"pickup"`
	ShippingAvailability       ShippingAvailability       `json:"shipping"`
}

type ProductAvailabilitiesResponse struct {
	ProductAvailabilities []ProductAvailability `json:"availabilities"`
}

func (b *BestBuyResponseParser) ParseAvailabilityResponseMessage(respBody []byte) (map[string]ProductAvailability, error) {
	if len(respBody) == 0 {
		return nil, errors.New("Error: Response body is empty.")
	}

	var response ProductAvailabilitiesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	availabilityMap := make(map[string]ProductAvailability)
	for _, availability := range response.ProductAvailabilities {
		availabilityMap[availability.SKU] = availability
	}

	return availabilityMap, nil
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
	ItemName                string
	OriginalPrice           string
	DiscountAmount          string
	PriceAfterDiscount      string
	Stock                   string
	Seller                  string
	QueryExecutionTimestamp time.Time
}

//func (p *ParseResponseBodyResult) GenerateEmailTemplate() {
// --TODO Generates the necessary string fields gmailer to send email to user when query condition of the ResponseBody as parsed sucessfully.
//}

// Parse the ResponseBody based on query conditions (i.e. price, stock, etc).
// Might need to implement a sub function for each condition whic reuslt in boolean, then ParseResponseMessage evaluates the overall condition based on the sub functions.
func (f *Fetcher) ParseResponseMessage(reponseBody []byte) (*ParseResponseBodyResult, error) {
	return nil, nil
}

// Contains all the conditions for the item, each item have this struct. Site and URL then uniquely identifies each tracker.
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
	//creates tis based on json or other structured file format parsed for the file that stores this info
	return *i, nil
}
