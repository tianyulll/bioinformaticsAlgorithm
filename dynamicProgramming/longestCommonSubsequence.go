package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

/*
OutputLCS(backtrack, v, i, j)
    if i = 0 or j = 0
        return ""
    if backtracki, j = "↓"
        return OutputLCS(backtrack, v, i - 1, j)
    else if backtracki, j = "→"
        return OutputLCS(backtrack, v, i, j - 1)
    else
        return OutputLCS(backtrack, v, i - 1, j - 1) + vi
down -1
right 0
corner 1
*/
func outputlcs(backtrack [][]int, v []string, i, j int) string {
	if i == 0 || j == 0 {
		return ""
	}
	if backtrack[i][j] == -1 {
		return outputlcs(backtrack, v,i-1,j)
	} else if  backtrack[i][j] == 0 {
		return outputlcs(backtrack, v, i, j-1)
	} else {
		return outputlcs(backtrack, v, i-1,j-1) + v[i-1]
	}
}

/*
LCSBackTrack(v, w)
    for i ← 0 to |v|
        si, 0 ← 0
    for j ← 0 to |w|
        s0, j ← 0
    for i ← 1 to |v|
        for j ← 1 to |w|
            match ← 0
            if vi-1 = wj-1
                match ← 1
            si, j ← max{si-1, j , si,j-1 , si-1, j-1 + match }
            if si,j = si-1,j
                Backtracki, j ← "↓"
            else if si, j = si, j-1
                Backtracki, j ← "→"
            else if si, j = si-1, j-1 + match
                Backtracki, j ← "↘"
    return Backtrack
*/

func LCSBackTrack(v, w []string) [][]int {
	n := len(v)
	m := len(w)

	S := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		S[i] = make([]int, m+1)
	}
	backtrack := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		backtrack[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		S[i][0] = 0
	}
	for j := 0; j <= m; j++ {
		S[0][j] = 0
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			match := 0
			if v[i-1] == w[j-1] {
				match = 1
			}
			S[i][j] = max([]int{S[i-1][j], S[i][j-1], S[i-1][j-1] + match})
			if S[i][j] == S[i-1][j] {
				backtrack[i][j] = -1
			} else if S[i][j] == S[i][j-1] {
				backtrack[i][j] = 0
			} else if S[i][j] == S[i-1][j-1]+match {
				backtrack[i][j] = 1
			}
		}
	}
	return backtrack
}

func main() {
	//load data
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i := 0
	var s, t string
	for scanner.Scan() {
		if i == 0 {
			s = scanner.Text()
			i++
		}
		t = scanner.Text()
	}

	//
	v := strings.Split(s, "")
	w := strings.Split(t, "")
	
	bt := LCSBackTrack(v,w )
	fmt.Println(bt)
	//fmt.Println(bt)
	fmt.Println(outputlcs(bt, v, len(v), len(w)))
//
}

func max(m []int) int {
	res := math.MinInt
	for _, val := range m {
		if val > res {
			res = val
		}
	}
	return res
}
