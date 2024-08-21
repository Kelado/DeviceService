package models

type Filter interface {
	GetField() string
	GetValue() string
}

type BrandFilter struct {
	Field string
	Value string
}

func NewBrandFilter(v string) BrandFilter {
	return BrandFilter{
		Field: "brand",
		Value: v,
	}
}

func (f BrandFilter) GetField() string {
	return f.Field
}

func (f BrandFilter) GetValue() string {
	return f.Value
}
