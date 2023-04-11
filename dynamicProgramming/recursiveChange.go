package main

import (
	"fmt"
	"math"
)




/*
Input: An integer money and an array Coins = (coin1, ..., coind).
Output: The minimum number of coins with denominations Coins that changes money.
DPChange(money, Coins)
    MinNumCoins(0) ← 0
    for m ← 1 to money
        MinNumCoins(m) ← ∞
            for i ← 0 to |Coins| - 1
                if m ≥ coini
                    if MinNumCoins(m - coini) + 1 < MinNumCoins(m)
                        MinNumCoins(m) ← MinNumCoins(m - coini) + 1
    output MinNumCoins(money)
*/

func dpChange(money int, coins []int) int {
	minNum := make([]int, money)
	minNum = append(minNum, 0)
	for m := 1; m <=money;m++ {
		minNum[m] = math.MaxInt
		for _, val := range coins {
			if m >= val {
				if minNum[m - val] + 1 < minNum[m] {
					minNum[m] = minNum[m - val] + 1
				}
			}
		}
	}
	
	return minNum[money]
}


func main() {
	coins := []int {24,15, 5, 3, 1}
	fmt.Println(dpChange(17287, coins))
}




/*
RecursiveChange(money, Coins)
    if money = 0
        return 0
    MinNumCoins ← ∞
    for i ← 0 to |Coins| - 1
        if money ≥ coini
            NumCoins ← RecursiveChange(money − coini, Coins)
            if NumCoins + 1 < MinNumCoins
                MinNumCoins ← NumCoins + 1
    return MinNumCoins
*/
func RecursiveChange(money int, coins []int) int {

	if money == 0 {
		return 0
	}

	minCoin := math.MaxInt
	for i := 0; i < len(coins); i++ {
		if money >= coins[i] {
			numCoin := RecursiveChange(money-coins[i], coins)
			if numCoin+1 < minCoin {
				minCoin = numCoin + 1
			}
		}
	}
	return minCoin
}