package sort

import (
	"errors"
	"strings"
)

type Sort string

const (
	New           Sort = "new"
	Hot           Sort = "hot"
	Rising        Sort = "rising"
	Controversial Sort = "controversial"
	Top           Sort = "top"
)

func Parse(s string) (Sort, error) {
	switch strings.ToLower(s) {
	case "new", "hot", "rising", "controversial", "top":
		return Sort(s), nil
	default:
		return New, errors.New("failed to parse sorting value")
	}
}

func ParseDefault(s string) Sort {
	sort, _ := Parse(s)
	return sort
}

func (s Sort) MarshalText() ([]byte, error) {
	return []byte(s), nil
}

func (s Sort) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s + `"`), nil
}

func (s *Sort) UnmarshalJSON(b []byte) error {
	if string(b) == "" {
		*s = New
		return nil
	}
	sort, err := Parse(string(b))
	if err != nil {
		return err
	}
	*s = sort
	return nil
}

func (s *Sort) UnmarshalText(b []byte) error {
	if string(b) == "" {
		*s = New
		return nil
	}
	sort, err := Parse(string(b))
	if err != nil {
		return err
	}
	*s = sort
	return nil
}
