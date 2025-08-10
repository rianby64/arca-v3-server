package models

type Material struct {
	Name         string  `json:"name"`
	Thickness    float64 `json:"thickness"`
	Keynote      string  `json:"keynote"`
	IsStructural bool    `json:"is_structural"`
}

type Materials []*Material

type AreaMaterials struct {
	Area      string    `json:"name"`
	Materials Materials `json:"materials"`
}

type AreasMaterials []*AreaMaterials

type AreaKey struct {
	AreaInternal string `json:"area_internal"`
	AreaExternal string `json:"area_external"`
	Keynote      string `json:"keynote"`
}

type AreasKeys []*AreaKey

type Relation struct {
	AreaInternal string    `json:"area_internal"`
	AreaExternal string    `json:"area_external"`
	Material     *Material `json:"material"`
	SameArea     bool      `json:"same_area"`
}

type Relations []*Relation
