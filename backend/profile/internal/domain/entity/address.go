package entity

import "time"

// Address represents a user's shipping address
type Address struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Province     string    `json:"province"`
	City         string    `json:"city"`
	District     string    `json:"district"`
	Address      string    `json:"address"`
	SignerName   string    `json:"signer_name"`
	SignerMobile string    `json:"signer_mobile"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// FullAddress returns the complete address as a string
func (a *Address) FullAddress() string {
	return a.Province + a.City + a.District + a.Address
}
