package main

import (
	"github.com/DesistDaydream/tcg-probability/pkg/combination"
	cbn "github.com/DesistDaydream/tcg-probability/pkg/combination"
	"github.com/DesistDaydream/tcg-probability/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
// HandCards = []string{"lv.3"}
)

type Flags struct {
	DeckSize      int
	HandSize      int
	OnlyCalculate bool
}

func (flags *Flags) AddFlags() {
	pflag.IntVarP(&flags.DeckSize, "deck-size", "d", 50, "卡组总数")
	pflag.IntVarP(&flags.HandSize, "hand-size", "h", 5, "手牌总数")
	pflag.BoolVarP(&flags.OnlyCalculate, "calculate", "c", false, "是否只通过数学计算获取结果")
}

type WantCardInfo struct {
	Name string
	// 在卡组中有多少张卡
	Have int
	Min  int
	// 暂时不知道如何计算最多有多少张想要的牌
	Max int
}

type WantCardsInfo struct {
	WantCards []WantCardInfo
}

func NewWantCardsInfo() *WantCardsInfo {
	wantCards := []WantCardInfo{
		{Name: "Lv.3", Have: 13, Min: 1, Max: 13},
		// {Name: "Lv.4", Have: 10, Min: 1},
	}

	return &WantCardsInfo{
		WantCards: wantCards,
	}
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
		Deck              []string
		TargetCombination int      = 0 // 满足条件的手牌组合数
		wantHandCards     []string     // 想要抓到手上的卡牌
	)

	wantCardsInfo := NewWantCardsInfo()

	for _, wantcard := range wantCardsInfo.WantCards {
		logrus.Infof(
			"牌组中有 %v 张【\033[0;31;31m %v \033[0m】，我们想要最少【\033[0;31;31m %v \033[0m】张、最多【\033[0;31;31m %v \033[0m】张",
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

	// 填充牌组中空余位置
	for i := 0; i < flags.DeckSize; i++ {
		if len(Deck) < flags.DeckSize {
			Deck = append(Deck, "any")
		}
	}

	logrus.Debugf("当前牌组：%v", Deck)
	logrus.Debugf("想要的最少手牌：%v", wantHandCards)
	// ！！！注意：这里暂时只能计算想要手牌中最少存在几张A，几张B的情况，默认最多可以有所有A、B、等等

	if flags.OnlyCalculate {
		// 所有可能的组合总数
		all := combination.Combination(flags.DeckSize, flags.HandSize).Int64()
		logrus.Infof("%v", all)
		// TODO: 计算想要的组合总数

	} else {
		// 遍历牌组，获取牌组中所有组合种类的列表
		combinations := cbn.TraversalDeckCombination(Deck, cbn.CombinationIndexs(flags.DeckSize, flags.HandSize))

		logrus.Debugf("原始组合总数: %v", len(combinations))
		cbn.CheckResult(flags.DeckSize, flags.HandSize, combinations)

		// 获取牌组中指定组合的总数
		for _, combination := range combinations {
			if cbn.ConditionCount(combination, wantHandCards) {
				TargetCombination++
			}
		}
		logrus.Infof("从 %v 张牌的卡组中抽 %v 张卡，包含上述想要的最少手牌的概率为 %v。", flags.DeckSize, flags.HandSize, float64(TargetCombination)/float64(len(combinations)))
	}
}
