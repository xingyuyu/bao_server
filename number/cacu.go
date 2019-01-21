package number

//MyNumber type
type MyNumber struct {
	X int32
	Y int32
}

//Cacaldate interface
type Cacaldate interface {
	Add() int32
	Mutil() int32
	Sub() int32
}

//Add number
func (p *MyNumber) Add() int32 {
	return p.X + p.Y
}

func (p *MyNumber) Mutil() int32 {
	return p.Y * p.Y
}

func (p *MyNumber) Sub() int32 {
	return p.X - p.Y
}
