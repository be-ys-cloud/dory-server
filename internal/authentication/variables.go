package authentication

import (
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"sync"
)

var Keys = map[string]structures.Key{}
var Mutex = sync.Mutex{}
