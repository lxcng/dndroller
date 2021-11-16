package model

import "fmt"

type RepeatingDice struct {
	*Dice
	Num int
}

func NewRepeatingDice(n int, d *Dice) *RepeatingDice {
	return &RepeatingDice{
		Dice: d,
		Num:  n,
	}
}

func (x *RepeatingDice) Roll() []int {
	res := make([]int, 0, x.Num)
	for i := 0; i < x.Num; i++ {
		res = append(res, x.Dice.Roll())
	}
	return res
}

func (x *RepeatingDice) String() string {
	return fmt.Sprintf("%dd%d", x.Num, x.Size)
}
