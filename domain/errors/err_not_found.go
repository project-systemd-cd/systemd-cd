package errors

import "fmt"

type ErrNotFound struct {
	Object string
	IdName string
	Id     string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("\"%s\" not found (%s: \"%s\")", e.Object, e.IdName, e.Id)
}

func (e *ErrNotFound) As(t any) bool {
	switch t.(type) {
	case **ErrNotFound, *ErrNotFound, ErrNotFound:
		return true
	default:
		return false
	}
}
