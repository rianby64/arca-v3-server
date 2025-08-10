package models

type CustomError string

func (e CustomError) Error() string {
	return string(e)
}

const (
	ErrNotFound    = CustomError("not found")
	ErrUnavailable = CustomError("unavailable")
)

type Material struct {
	Name         string  `json:"name"`
	Thickness    float64 `json:"thickness"`
	Keynote      string  `json:"keynote"`
	IsStructural bool    `json:"is_structural"`
}

type Materials []*Material

type Area struct {
	Name string `json:"name"`
}

type Areas []*Area

type AreaMaterials struct {
	Area      *Area     `json:"area"`
	Materials Materials `json:"materials"`
}

type AreasMaterials []*AreaMaterials

type AreaKey struct {
	AreaInternal *Area  `json:"area_internal"`
	AreaExternal *Area  `json:"area_external"`
	Keynote      string `json:"keynote"`
}

type AreasKeys []*AreaKey

type Relation struct {
	AreaInternal *Area     `json:"area_internal"`
	AreaExternal *Area     `json:"area_external"`
	Material     *Material `json:"material"`
	SameArea     bool      `json:"same_area"`
}

type Relations []*Relation
