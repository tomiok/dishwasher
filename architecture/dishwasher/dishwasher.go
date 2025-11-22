package dishwasher

const (
	SizeMid     DishSize = 768
	PathDefault          = "/dishwasher.db"
)

type DishSize int32

type Dishwasher struct {
	Path          string
	DimensionSize DishSize
	InMemory      bool
}

func New() Dishwasher {
	return Dishwasher{
		Path:          PathDefault,
		DimensionSize: SizeMid,
	}
}

func (d Dishwasher) WithPath(s string) Dishwasher {
	res := d
	res.Path = s
	return res
}

func (d Dishwasher) WithDimensionSize(size DishSize) Dishwasher {
	res := d
	res.DimensionSize = size

	return res
}
