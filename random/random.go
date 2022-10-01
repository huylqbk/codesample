package random

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func Rand(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n) + 1
}

func RandDistribution(n int) int {
	x := Rand(n / 2)
	y := Rand(n / 2)
	return x + ((y&1)*n)/2
}

func NewUUID() string {
	return uuid.New().String()
}
