package model

// FilterOption — одна опция фильтра.
type FilterOption struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// FilterCategory — категория фильтров со списком опций.
type FilterCategory struct {
	Code    string         `json:"code"`
	Label   string         `json:"label"`
	Options []FilterOption `json:"options"`
}
