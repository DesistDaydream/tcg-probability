package main

import (
	"github.com/DesistDaydream/tcg-probability/cmd/calculate"
	"github.com/DesistDaydream/tcg-probability/pkg/combination"
	"github.com/DesistDaydream/tcg-probability/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type Flags struct {
	DeckSize       int
	HandSize       int
	DoNotCalculate bool
}

func (flags *Flags) AddFlags() {
	pflag.IntVarP(&flags.DeckSize, "deck-size", "d", 50, "卡组总数")
	pflag.IntVarP(&flags.HandSize, "hand-size", "h", 5, "手牌总数")
	pflag.BoolVarP(&flags.DoNotCalculate, "calculate", "c", false, "是否只通过数学计算获取结果")
}

func main() {
	// 设置命令行标志
	ygoFlags := &logging.LoggingFlags{}
	ygoFlags.AddFlags()
	flags := &Flags{}
	flags.AddFlags()
	pflag.Parse()

	// 初始化日志
	if err := logging.LogInit(ygoFlags.LogLevel, ygoFlags.LogOutput, ygoFlags.LogFormat); err != nil {
		logrus.Fatal("初始化日志失败", err)
	}

	// 设置变量
	var (
		Deck          []string
		wantHandCards []string // 想要抓到手上的卡牌
	)

	cardsInfo := calculate.NewCardsInfo(flags.DeckSize, flags.HandSize)

	for _, wantcard := range cardsInfo.WantCards {
		logrus.Infof(
			"卡组中有 %v 张【\033[0;31;31m %v \033[0m】，我们想要最少【\033[0;31;31m %v \033[0m】张、最多【\033[0;31;31m %v \033[0m】张",
			wantcard.Have, wantcard.Name, wantcard.Min, wantcard.Max)
		// 将手牌填充到卡组中
		for i := 0; i < wantcard.Have; i++ {
			Deck = append(Deck, wantcard.Name)
		}
		// 将想要的卡牌保存到数组变量中
		for i := 0; i < wantcard.Min; i++ {
			wantHandCards = append(wantHandCards, wantcard.Name)
		}
	}

	// 填充卡组中空余位置
	for i := 0; i < flags.DeckSize; i++ {
		if len(Deck) < flags.DeckSize {
			Deck = append(Deck, "any")
		}
	}

	logrus.Debugf("当前卡组：%v", Deck)
	logrus.Debugf("想要的最少手牌：%v", wantHandCards)
	// ！！！注意：这里暂时只能计算想要手牌中最少存在几张A，几张B的情况，默认最多可以有所有A、B、等等

	if flags.DoNotCalculate {
		var TargetCombination int = 0 // 满足条件的手牌组合数
		// 遍历卡组，获取卡组中所有组合种类的列表
		combinations := combination.TraversalDeckCombination(Deck, combination.CombinationIndexs(flags.DeckSize, flags.HandSize))

		logrus.Debugf("原始组合总数: %v", len(combinations))
		combination.CheckResult(flags.DeckSize, flags.HandSize, combinations)

		// 获取卡组中指定组合的总数
		for _, c := range combinations {
			if combination.ConditionCount(c, wantHandCards) {
				TargetCombination++
			}
		}
		logrus.Infof("从 %v 张牌的卡组中抽 %v 张卡，包含上述想要的最少手牌的概率为 %v。", flags.DeckSize, flags.HandSize, float64(TargetCombination)/float64(len(combinations)))

	} else {
		var currentHand []int
		// 所有可能的组合总数
		all := combination.Combination(flags.DeckSize, flags.HandSize)
		// 计算想要的组合总数
		result := cardsInfo.RecursiveCalculate(currentHand, 0, flags.HandSize)

		logrus.WithFields(logrus.Fields{
			"总组合数":        all,
			"包含想要的卡牌的组合数": result,
			"概率":          float64(result) / float64(all),
		}).Infof("从 %v 张牌的卡组中抽 %v 张卡", flags.DeckSize, flags.HandSize)
	}
}
