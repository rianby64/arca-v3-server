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
	Name                          *string
	MaterialCategory              *string
	CutBackgroundPatternColor     *string
	CutBackgroundPatternId        *string
	CutForegroundPatternColor     *string
	CutForegroundPatternId        *string
	SurfaceForegroundPatternColor *string
	SurfaceForegroundPatternId    *string
	Mark                          *string
	Keynote                       *string
	Description                   *string
	Manufacturer                  *string
}

type Materials []*Material

type WallMaterial struct {
	Thickness    float64
	Function     string
	IsStructural bool
	Material     *Material
}

type WallMaterials []*WallMaterial

type Area struct {
	Name string
}

type Areas []*Area

type AreaMaterials struct {
	Area      *Area
	Materials WallMaterials
}

type AreasMaterials []*AreaMaterials

type AreaRelation struct {
	AreaInternal *Area
	AreaExternal *Area
	Material     *WallMaterial
	WallKeynote  *string
	SameArea     bool
}

type AreasRelations []*AreaRelation
