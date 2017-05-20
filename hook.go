package logs

type Hook interface {
	Name() string
	IsAsyncFire() bool
	Levels() []Level
	Fire(Element) error
	OnError(error)
}


