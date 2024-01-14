package frontend

import (
	"hexarch/core"
)

type zeroFrontEnd struct{}

func (f zeroFrontEnd) Start(store *core.KeyValueStore) error {
	return nil
}
