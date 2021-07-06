package sort

import (
	"strconv"
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
	case "new":
		return New, nil
	case "hot":
		return Hot, nil
	case "rising":
		return Rising, nil
	case "controversial":
		return Controversial, nil
	case "top":
		return Top, nil
	default:
		return New, &strconv.NumError{Func: "Parse", Num: s, Err: strconv.ErrRange}
	}
}

func ParseDefault(s string) Sort {
	sort, _ := Parse(s)
	return sort
}

func (s Sort) String() string {
	switch s {
	case Hot:
		return "hot"
	case Rising:
		return "rising"
	case Top:
		return "Top"
	case Controversial:
		return "controversial"
	default:
		return "new"
	}
}

func (s Sort) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s Sort) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
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
