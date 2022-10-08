package calculate

import (
	"fmt"

	"github.com/DesistDaydream/tcg-probability/pkg/combination"
	"github.com/sirupsen/logrus"
)

type CardsInfo struct {
	WantCards []WantCardInfo
	MiscHave  int
	MiscMin   int
	MiscMax   int
}

type WantCardInfo struct {
	Name string
	Have int // 在卡组中有多少张卡
	Min  int
	Max  int
}

func NewCardsInfo(deckSize, handSize int) *CardsInfo {
	var (
		miscHave int = deckSize
		miscMin  int
		miscMax  int = handSize
	)
	wantCards := []WantCardInfo{
		{Name: "Lv.3", Have: 13, Min: 1, Max: 13},
		// {Name: "Lv.4", Have: 10, Min: 1, Max: 10},
	}

	for _, card := range wantCards {
		miscHave = miscHave - card.Have
		miscMax = miscMax - card.Min
	}

	if miscMin < 0 {
		logrus.Fatalln("数据错误，无法计算，请修正数据")
	}

	return &CardsInfo{
		WantCards: wantCards,
		MiscHave:  miscHave,
		MiscMin:   miscMin,
		MiscMax:   miscMax,
	}
}

// 使用纯数学计算的方式获取指定条件下的组合数
func (c *CardsInfo) RecursiveCalculate(currentHand []int, currentHandSize int, handSize int) int64 {
	if len(c.WantCards) == 0 || currentHandSize >= handSize {
		if currentHandSize == handSize {
			logrus.Debugf("当前手牌容量已经等于手牌容量，检查想要卡片长度：%v", len(c.WantCards))
			var noChance = false
			for i := 0; i < len(c.WantCards); i++ {
				if c.WantCards[i].Min != 0 {
					noChance = true
					break
				}
			}

			if noChance {
				return 0
			}
		} else if currentHandSize > handSize {
			return 0
		}

		var newChance int64 = 1
		var output string = ""

		// ！！计算部分！！
		for i := 0; i < len(currentHand); i += 2 {
			output += fmt.Sprintf("(%v choose  %v) * ", currentHand[i], currentHand[i+1])
			newChance *= combination.Combination(currentHand[i], currentHand[i+1])
		}

		if currentHandSize < handSize {
			output += fmt.Sprintf("(%v choose %v *)", c.MiscHave, handSize-currentHandSize)
			newChance *= combination.Combination(c.MiscHave, handSize-currentHandSize)
		}

		logrus.Debugf(output)

		return newChance
	}

	obj := c.WantCards[len(c.WantCards)-1]
	c.WantCards = c.WantCards[:len(c.WantCards)-1]
	var chance int64 = 0

	for i := obj.Min; i <= obj.Max; i++ {
		currentHand = append(currentHand, obj.Have)
		currentHand = append(currentHand, i)

		chance += c.RecursiveCalculate(currentHand, currentHandSize+i, handSize)

		currentHand = currentHand[:len(currentHand)-1]
		currentHand = currentHand[:len(currentHand)-1]
	}

	c.WantCards = append(c.WantCards, obj)

	return chance
}
