package combination

import (
	"log"
	"math/big"
)

// 阶乘
func Factorial(n uint64) uint64 {
	if n > 0 {
		return n * Factorial(n-1)
	}
	return 1
}

// 大数阶乘
func BigFactorial(n *big.Int) *big.Int {
	if n.Int64() == 1 {
		return big.NewInt(1)
	} else {
		return n.Mul(n, BigFactorial(big.NewInt(n.Int64()-1)))
	}
}

// 组合
func Combination(n, k int) *big.Int {
	result := big.NewInt(0)
	if n <= k {
		log.Println(n, k)
		return result
	}
	nF := BigFactorial(big.NewInt(int64(n)))
	kF := BigFactorial(big.NewInt(int64(k)))
	nkF := BigFactorial(big.NewInt(int64(n - k)))
	return result.Div(result.Div(nF, kF), nkF)
}
