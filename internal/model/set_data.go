package model

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	regexSet = "[0-9]+[dD][0-9]+"
	rSet     = regexp.MustCompile(regexSet)
	regexRd  = "[0-9]+"
	rRd      = regexp.MustCompile(regexRd)
)

type DiceSet struct {
	Dices []*RepeatingDice
}

func NewDiceSet(str string) *DiceSet {
	return &DiceSet{
		Dices: parseSet(str),
	}
}

func NewDiceSetEmpty() *DiceSet {
	return &DiceSet{
		Dices: []*RepeatingDice{},
	}
}

func (x *DiceSet) Add(d int) *DiceSet {
	for _, rd := range x.Dices {
		if rd.Size == d {
			rd.Num++
			return x
		}
	}
	x.Dices = append(x.Dices, NewRepeatingDice(1, NewDice(d)))
	return x
}

func (x *DiceSet) Sub(d int) *DiceSet {
	for i, rd := range x.Dices {
		if rd.Size == d {
			if rd.Num > 1 {
				rd.Num--
			} else {
				newDices := make([]*RepeatingDice, len(x.Dices)-1)
				copy(newDices, x.Dices[:i])
				copy(newDices[i:], x.Dices[i+1:])
				x.Dices = newDices
				return x
			}
		}
	}
	return x
}

func (x *DiceSet) Roll() *Roll {
	res := make([][]int, 0, len(x.Dices))
	for _, rd := range x.Dices {
		res = append(res, rd.Roll())
	}
	return NewRoll(res, x)
}

func (x *DiceSet) layout() []string {
	layout := make([]string, 0, len(x.Dices)+1)
	for _, rd := range x.Dices {
		layout = append(layout, rd.Header())
	}
	return layout
}

func (x *DiceSet) String() string {
	if len(x.Dices) == 0 {
		return "Empty"
	}
	strs := make([]string, 0, len(x.Dices))
	for _, rd := range x.Dices {
		strs = append(strs, rd.String())
	}
	return strings.Join(strs, " ")
}

func parseSet(str string) []*RepeatingDice {
	rdsRaw := rSet.FindAllString(str, -1)

	uniqueSize := map[int]int{}
	for _, rd := range rdsRaw {
		n, s := parseRd(rd)
		uniqueSize[s] += n
	}
	res := make([]*RepeatingDice, 0, len(uniqueSize))
	for s, n := range uniqueSize {
		res = append(res, NewRepeatingDice(n, NewDice(s)))
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Size < res[j].Size
	})
	return res
}

func parseRd(str string) (int, int) {
	numsRaw := rRd.FindAllString(str, -1)
	if len(numsRaw) != 2 {
		panic("failed to parse rd")
	}
	num, err := strconv.Atoi(numsRaw[0])
	if err != nil {
		panic("failed to parse dice num")
	}
	size, err := strconv.Atoi(numsRaw[1])
	if err != nil {
		panic("failed to parse dice size")
	}
	return num, size
}
