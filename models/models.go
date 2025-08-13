package models

type CustomError string

func (e CustomError) Error() string {
	return string(e)
}

const (
	ErrNoData      = CustomError("no data")
	ErrInvalid     = CustomError("invalid")
	ErrNotFound    = CustomError("not found")
	ErrUnavailable = CustomError("unavailable")
)

type Material struct {
	Name         string
	Thickness    float64
	Keynote      string
	IsStructural bool
}

type Materials []*Material

type Area struct {
	Name string
}

type Areas []*Area

type AreaMaterials struct {
	Area      *Area
	Materials Materials
}

type AreasMaterials []*AreaMaterials

type AreaRelation struct {
	AreaInternal *Area
	AreaExternal *Area
	Material     *Material
	WallKeynote  string
	SameArea     bool
}

type AreasRelations []*AreaRelation
