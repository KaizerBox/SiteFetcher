package bestbuy

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"siteFetcher/internal/fetcher"
)

// BestBuyRequestUrlBuilder creates BestBuy-specific URLs for product, location, and availability requests.
type BestBuyRequestUrlBuilder struct {
	Name              string
	Scheme            string
	Host              string
	ProductURL        string
	LocationUrl       string
	AvailabilityUrl   string
	AvailabilityQuery map[string]string
}

func NewBestBuyRequestUrlBuilder(name, scheme, host, productURL, locationURL, availabilityURL string) (*BestBuyRequestUrlBuilder, error) {
	if name == "" || scheme == "" || host == "" {
		return nil, errors.New("invalid parameters for BestBuyRequestUrlBuilder")
	}
	return &BestBuyRequestUrlBuilder{Name: name, Scheme: scheme, Host: host, ProductURL: productURL, LocationUrl: locationURL, AvailabilityUrl: availabilityURL}, nil
}

func NewBestBuyRequestUrlBuilderWithConfig(name, scheme, host, productURL, locationURL, availabilityURL string, availabilityQuery map[string]string) (*BestBuyRequestUrlBuilder, error) {
	builder, err := NewBestBuyRequestUrlBuilder(name, scheme, host, productURL, locationURL, availabilityURL)
	if err != nil {
		return nil, err
	}

	builder.AvailabilityQuery = availabilityQuery
	return builder, nil
}

func (b *BestBuyRequestUrlBuilder) GetProductUrl(productNum string) (string, error) {
	if _, err := strconv.Atoi(productNum); len(productNum) == 0 || err != nil {
		return "", fmt.Errorf("Error: Product Number: %s is invalid.", productNum)
	}

	u := &url.URL{Scheme: b.Scheme, Host: b.Host, Path: fmt.Sprintf(b.ProductURL, productNum)}
	return u.String(), nil
}

func (b *BestBuyRequestUrlBuilder) GetLocationsUrl(postalCodeFSA string) (string, error) {
	valid, regexString, err := fetcher.ValidatePostalCodeFSA(postalCodeFSA)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", fmt.Errorf("Error: Postal Code FSA: %s is invalid, should match regex: %s", postalCodeFSA, regexString)
	}

	u := &url.URL{Scheme: b.Scheme, Host: b.Host, Path: b.LocationUrl}
	q := u.Query()
	q.Set("lang", "en-CA")
	q.Set("postalCode", postalCodeFSA)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (b *BestBuyRequestUrlBuilder) GetAvailabilityUrl(productNums []string, postalCodeFSA string, storeId []string) (string, error) {
	if len(productNums) == 0 {
		return "", errors.New("Error: Product Number list is Empty")
	}
	if valid, regexString, err := fetcher.ValidatePostalCodeFSA(postalCodeFSA); err != nil {
		return "", err
	} else if !valid {
		return "", fmt.Errorf("Error: Postal Code FSA: %s is invalid, should match regex: %s", postalCodeFSA, regexString)
	}
	if len(storeId) == 0 {
		return "", errors.New("Error: Store ID list is Empty")
	}

	queryConfig := map[string]string{
		"accept":          "application/vnd.bestbuy.standardproduct.v1+json",
		"accept-language": "en-CA",
		"locations":       strings.Join(storeId, "|"),
		"postalCode":      postalCodeFSA,
		"skus":            strings.Join(productNums, "|"),
	}

	for key, value := range b.AvailabilityQuery {
		if value != "" {
			queryConfig[key] = value
		}
	}

	u := &url.URL{Scheme: b.Scheme, Host: b.Host, Path: b.AvailabilityUrl}
	q := u.Query()
	for key, value := range queryConfig {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
