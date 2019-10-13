package routex

import "errors"

var ErrPropConflict = errors.New("same key exists in route props")
var ErrPropIsNil = errors.New("props is not initialized")
