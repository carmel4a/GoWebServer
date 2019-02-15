package main

type StringError struct {
	s string
}

func (p StringError) Error() string {
	return p.s
}
