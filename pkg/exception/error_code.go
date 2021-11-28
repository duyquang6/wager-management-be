// Package exception provides error code definition
package exception

// Common module (00) error codes definition.
var (
	ErrValidation     = 4000001
	ErrUnauthorized   = 4010001
	ErrInternalServer = 5000001
	ErrUnknownError   = 500
)

// Purchase module (01) error codes definition.
var (
	ErrRelatedWagerNotFound                      = 4040101
	ErrBuyingPriceGreaterThanCurrentSellingPrice = 4000101
)

// Wager module (02) error codes definition.
var (
	ErrWagerNotFound = 4040201
	ErrPaginationParamInvalid = 4000201
)
