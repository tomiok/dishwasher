package server

import (
	"fmt"
	"math/rand"
	"time"
)

func generateID(prefix string) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	v := r.Int()
	id := fmt.Sprintf("%s-%d", prefix, v)

	return id
}
