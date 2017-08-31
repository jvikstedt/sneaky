package engine

import "sync"

var mu sync.Mutex
var cid int = 1

func nextID() int {
	mu.Lock()
	nextID := cid
	cid = cid + 1
	mu.Unlock()
	return nextID
}
