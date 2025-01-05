package workshop

import "fmt"

type Section struct {
	Number int
	Slug   string
}

func (s *Section) Descriptor() string {
	return fmt.Sprintf("%0.2d-%s", s.Number, s.Slug)
}
