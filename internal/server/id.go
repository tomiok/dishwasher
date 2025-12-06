package server

import (
	"fmt"
	"math/rand"
	"time"
)

func generateID(seed bool) string {
	var prefix string

	if seed {
		prefix = "master"
	} else {
		prefix = "node"
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	v := r.Int()
	id := fmt.Sprintf("%s-%d", prefix, v)

	return id
}
