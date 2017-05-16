package logs

type Hook interface {
	Formatter
	Levels() []Level
	Fire(Element) error
}


