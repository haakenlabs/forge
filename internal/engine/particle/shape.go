package particle

type ShapeType int

const (
	ShapeSphere ShapeType = iota
	ShapeCone
	ShapeBox
)

type ModuleShape struct {
	Shape ShapeType
}

func NewModuleShape() *ModuleShape {
	m := &ModuleShape{}

	return m
}
