package contract

type Input interface {
	Write([]byte) (int, error)
}

type Output interface {
	Read([]byte) (int, error)
}

type IO interface {
	Input
	Output
}
