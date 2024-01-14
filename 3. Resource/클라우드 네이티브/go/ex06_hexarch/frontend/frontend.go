package frontend

import (
	"hexarch/core"
)

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}
