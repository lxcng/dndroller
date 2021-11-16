package model

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Dice struct {
	Size int
}

func NewDice(s int) *Dice {
	return &Dice{
		Size: s,
	}
}

func (x *Dice) Roll() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res := math.Ceil(r.Float64() * float64(x.Size))
	return int(res)
}

func (x *Dice) Header() string {
	return fmt.Sprintf("d%d", x.Size)
}
