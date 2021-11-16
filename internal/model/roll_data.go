package model

import (
	"fmt"
	"strings"
)

type Roll struct {
	Result         [][]int
	RerollIndecies map[[2]int]struct{}
	Set            *DiceSet
}

func NewRoll(r [][]int, set *DiceSet) *Roll {
	return &Roll{
		Result:         r,
		RerollIndecies: map[[2]int]struct{}{},
		Set:            set,
	}
}

func (x *Roll) ToMarkdown() string {
	lines := make([]string, 0, len(x.Set.Dices)+2)
	layout := x.Set.layout()
	lines = append(lines, x.Set.String())
	total := 0
	for i, lo := range layout {
		line := fmt.Sprintf("%s:", lo)
		for j, r := range x.Result[i] {
			if _, ok := x.RerollIndecies[[2]int{i, j}]; ok {
				line += fmt.Sprintf(" ~_%d_~", r)
			} else {
				line += fmt.Sprintf(" *%d*", r)
			}
			total += r
		}
		lines = append(lines, line)
	}
	lines = append(lines, fmt.Sprintf("Total: *%d*", total))
	return strings.Join(lines, "\n")
}

func (x *Roll) RerollDice(i, j int) {
	if _, ok := x.RerollIndecies[[2]int{i, j}]; ok {
		delete(x.RerollIndecies, [2]int{i, j})
	} else {
		x.RerollIndecies[[2]int{i, j}] = struct{}{}
	}
}

func (x *Roll) FormatDice(i, j int) string {
	if _, ok := x.RerollIndecies[[2]int{i, j}]; ok {
		return strikethrough(x.Result[i][j])
	} else {
		return fmt.Sprintf("%d", x.Result[i][j])
	}
}

func strikethrough(n int) string {
	str := fmt.Sprint(n)
	res := ""
	for _, r := range str {
		res += fmt.Sprintf("%s\u0332", string(r))
	}
	return res
}

func (x *Roll) Reroll() {
	for ind := range x.RerollIndecies {
		x.Result[ind[0]][ind[1]] = x.Set.Dices[ind[0]].Dice.Roll()
	}
	x.RerollIndecies = map[[2]int]struct{}{}
}
