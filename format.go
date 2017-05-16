package logs

type Formatter interface {
	Format(Element) []byte
}
