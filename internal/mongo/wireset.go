package mongo

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewClient,
)
