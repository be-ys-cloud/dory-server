package authentication

import (
	"structures"
	"sync"
)

var Keys = map[string]structures.Key{}
var Mutex = sync.Mutex{}