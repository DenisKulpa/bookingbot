package model

type City struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

type Zone struct {
	ID             int      `json:"id"`
	CityID         int      `json:"city_id"`
	CityName       string   `json:"city_name,omitempty"`
	Name           string   `json:"name"`
	Emoji          string   `json:"emoji,omitempty"`
	ShortDesc      string   `json:"short_desc,omitempty"`
	FullDesc       string   `json:"full_desc,omitempty"`
	TargetAudience string   `json:"target_audience,omitempty"`
	Pros           []string `json:"pros,omitempty"`
	Cons           []string `json:"cons,omitempty"`
	HousingTypes   []string `json:"housing_types,omitempty"`
	PriceLevel     int      `json:"price_level,omitempty"`
	BestFor        string   `json:"best_for,omitempty"`
	SeasonNote     string   `json:"season_note,omitempty"`
	SortOrder      int      `json:"sort_order"`
	IsActive       bool     `json:"is_active"`
}

type Subzone struct {
	ID        int    `json:"id"`
	ZoneID    int    `json:"zone_id"`
	ZoneName  string `json:"zone_name,omitempty"`
	CityName  string `json:"city_name,omitempty"`
	Name      string `json:"name"`
	Emoji     string `json:"emoji,omitempty"`
	ShortDesc string `json:"short_desc,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

type DistrictDetail struct {
	District *Zone     `json:"district"`
	Subzones []*Subzone `json:"subzones"`
}
