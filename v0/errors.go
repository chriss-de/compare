package compare

import "errors"

var (
	// ErrTypeMismatch Compared types do not match
	ErrTypeMismatch = errors.New("types do not match")
	// ErrInvalidChangeType The specified change values areKind not unsupported
	ErrInvalidChangeType = errors.New("change type must be one of 'create' or 'delete'")

	ErrNotCombinableIdentifier = errors.New("not combinable")
)
