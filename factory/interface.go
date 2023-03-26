package factory

type IFactory interface {
	Make() interface{}
}
