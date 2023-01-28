package errors

type ErrValidationMsg struct {
	Msg string
}

func (e *ErrValidationMsg) Error() string {
	return e.Msg
}

func (e *ErrValidationMsg) As(t any) bool {
	switch t.(type) {
	case **ErrValidation, *ErrValidation, ErrValidation,
		**ErrValidationMsg, *ErrValidationMsg, ErrValidationMsg:
		return true
	default:
		return false
	}
}
