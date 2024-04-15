package rerrs

import (
	td "github.com/swxctx/malatd"
)

var (
	// RerrInvalidParameter error
	RerrInvalidParameter = td.NewRerror(100001, "Invalid Parameter", "Contains invalid request parameters")
)
