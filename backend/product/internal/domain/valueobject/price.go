package valueobject

// Price 价格值对象
type Price struct {
	Value     float64
	Currency  string
	IsDefault bool
}

// NewPrice 创建价格值对象
func NewPrice(value float64, currency string) *Price {
	return &Price{
		Value:     value,
		Currency:  currency,
		IsDefault: true,
	}
}

// NewCNYPrice 创建人民币价格值对象
func NewCNYPrice(value float64) *Price {
	return NewPrice(value, "CNY")
}

// Equals 判断两个价格是否相等
func (p *Price) Equals(other *Price) bool {
	if p == nil || other == nil {
		return false
	}
	
	return p.Value == other.Value && p.Currency == other.Currency
}

// LessThan 判断价格是否小于另一个价格
func (p *Price) LessThan(other *Price) bool {
	if p == nil || other == nil || p.Currency != other.Currency {
		return false
	}
	
	return p.Value < other.Value
}

// Add 价格加法
func (p *Price) Add(other *Price) *Price {
	if p == nil || other == nil || p.Currency != other.Currency {
		return p
	}
	
	return &Price{
		Value:     p.Value + other.Value,
		Currency:  p.Currency,
		IsDefault: p.IsDefault,
	}
}

// Multiply 价格乘法
func (p *Price) Multiply(quantity int) *Price {
	if p == nil || quantity < 0 {
		return p
	}
	
	return &Price{
		Value:     p.Value * float64(quantity),
		Currency:  p.Currency,
		IsDefault: p.IsDefault,
	}
}

// Discount 价格折扣
func (p *Price) Discount(percentage float64) *Price {
	if p == nil {
		return p
	}
	
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}
	
	discountFactor := (100 - percentage) / 100
	
	return &Price{
		Value:     p.Value * discountFactor,
		Currency:  p.Currency,
		IsDefault: p.IsDefault,
	}
}
