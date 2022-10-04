package util

import (
	"math/rand"
	"sync"
	"time"
)

var randSeed *rand.Rand
var randLock sync.Mutex

func init() {
	randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
}
