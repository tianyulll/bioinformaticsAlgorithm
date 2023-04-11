package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
 Find the length of a longest path in the Manhattan Tourist Problem.

Input: Integers n and m, followed by an n × (m + 1) matrix Down and an (n + 1) × m matrix Right.
The two matrices are separated by the "-" symbol.
Output:
The length of a longest path from source (0, 0) to sink (n, m)
in the rectangular grid whose edges are defined by the matrices Down and Right.


*/

/*
ManhattanTourist(n, m, Down, Right)
    s0, 0 ← 0
    for i ← 1 to n
        si, 0 ← si-1, 0 + downi-1, 0
    for j ← 1 to m
        s0, j ← s0, j−1 + right0, j-1
    for i ← 1 to n
        for j ← 1 to m
            si, j ← max{si - 1, j + downi-1, j, si, j - 1 + righti, j-1}
    return sn, m
*/
func manhattantourist(n, m int, down, right [][]int) int {
	path := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		path[i] = make([]int, m+1)
	}
	path[0][0] = 0
	for i := 1; i <= n; i++ {
		path[i][0] = path[i-1][0] + down[i-1][0]
	}
	for j := 1; j <= m; j++ {
		path[0][j] = path[0][j-1] + right[0][j-1]
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			path[i][j] = max(path[i-1][j]+down[i-1][j], path[i][j-1]+right[i][j-1])
		}
	}
	fmt.Println(path)
	return path[n][m]
}

func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i, n, m := 0, 0, 0
	nextMatrix := false
	matrixN := make([][]int, 0)
	matrixM := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if i == 0 {
			tmp := strings.Split(line, " ")
			n, err = strconv.Atoi(tmp[0])
			m, err = strconv.Atoi(tmp[1])
			i++
			continue
		}
		if line == "-" {
			nextMatrix = true
			continue
		}
		if nextMatrix {
			matrixM = append(matrixM, strToInt(strings.Split(line, " ")))
		} else {
			matrixN = append(matrixN, strToInt(strings.Split(line, " ")))
		}
	}

	fmt.Println(n, m)
	fmt.Println(matrixM)
	fmt.Println(matrixN)

	//finished handling input data

	res := manhattantourist(n, m, matrixN, matrixM)
	fmt.Println(res)

}

func strToInt(input []string) []int {
	res := make([]int, 0)
	for _, val := range input {
		tmp, err := strconv.Atoi(val)
		check(err)
		res = append(res, tmp)
	}
	return res
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
