package errors

import (
	"fmt"
)

type ErrValidation struct {
	Property    string
	Given       *string
	Description string
}

func (e *ErrValidation) Error() string {
	// e.g.
	// the property "Name" given "" cannot be empty.
	// the property "Email" given "given-email" is invalid email address.
	if e.Given == nil {
		return fmt.Sprintf("the property \"%s\" given 'nil' %s.", e.Property, e.Description)
	}
	return fmt.Sprintf("the property \"%s\" given \"%s\" %s.", e.Property, *e.Given, e.Description)
}

func (e *ErrValidation) As(t any) bool {
	switch t.(type) {
	case **ErrValidation, *ErrValidation, ErrValidation,
		**ErrValidationMsg, *ErrValidationMsg, ErrValidationMsg:
		return true
	default:
		return false
	}
}
