package combination

import (
	"fmt"

	"github.com/DesistDaydream/yu-gi-oh/pkg/subset"
	"github.com/sirupsen/logrus"
)

//组合算法(从 n 中取出 k 个数)
func CombinationIndexs(n int, k int) [][]int {
	if k < 1 || k > n {
		fmt.Println("Illegal argument. Param m must between 1 and len(nums).")
		return [][]int{}
	}

	//保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]int, 0, Combination(n, k).Int64())
	//保存每一个组合的索引的数组，1表示选中，0表示未选中
	indexs := make([]int, n)
	for i := 0; i < n; i++ {
		if i < k {
			indexs[i] = 1
		} else {
			indexs[i] = 0
		}
	}

	//第一个结果
	result = addTo(result, indexs)
	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < n-1; i++ {
			if indexs[i] == 1 && indexs[i+1] == 0 {
				find = true

				indexs[i], indexs[i+1] = 0, 1
				if i > 1 {
					moveOneToLeft(indexs[:i])
				}
				result = addTo(result, indexs)

				break
			}
		}

		//本次循环没有找到 1 0 ，说明已经取到了最后一种情况
		if !find {
			break
		}
	}

	return result
}

//将ele复制后添加到arr中，返回新的数组
func addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)

	return arr
}

func moveOneToLeft(leftNums []int) {
	//计算有几个1
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}

	//将前sum个改为1，之后的改为0
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}

// 遍历牌组，获取所有组合的列表
func TraversalDeckCombination(nums []string, indexs [][]int) [][]string {
	if len(indexs) == 0 {
		return [][]string{}
	}

	result := make([][]string, len(indexs))

	for i, v := range indexs {
		line := make([]string, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, nums[j])
			}
		}
		result[i] = line
	}

	return result
}

// 统计
func ConditionCount(combination []string, condition []string) bool {
	return subset.IsSubset(condition, combination)
}

// 判断遍历所有组合数的结果是否正确
func CheckResult(n, k int, combinations [][]string) {
	rightCount := Combination(n, k).Int64()
	if int(rightCount) == len(combinations) {
		logrus.Debugln("数学计算结果与遍历结果相同")
	} else {
		logrus.WithFields(logrus.Fields{
			"数学计算结果": rightCount,
			"遍历结果":   len(combinations),
		}).Errorln("数学计算结果与遍历结果不同")
	}

}
