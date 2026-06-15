package fetcher

// Product is the shared response model used by store implementations.
type Product struct {
	Name                string                `json:"name"`
	SKU                 string                `json:"sku"`
	RegularPrice        float64               `json:"regularPrice"`
	SalePrice           float64               `json:"salePrice"`
	Seller              string                `json:"sellerId"`
	OnlineAvailability  OnlineAvailability    `json:"availability"`
	InStoreAvailability []InStoreAvailability `json:"inStoreAvailability"`
}

func (p *Product) GetDiscountAmount() float64 {
	return p.RegularPrice - p.SalePrice
}
