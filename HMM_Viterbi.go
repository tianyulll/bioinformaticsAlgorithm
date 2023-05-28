package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type node struct {
	prevX int
	prevY int
	state string
}

func viterbi(emitString string, transitLetter []string) ([][]float64, [][]node) {
	rowNum := len(transitMap)
	seq := strings.Split(emitString, "")

	dp := make([][]float64, rowNum)
	for i := range dp {
		dp[i] = make([]float64, len(emitString))
	}
	backtrack := make([][]node, rowNum)
	for i := range backtrack {
		backtrack[i] = make([]node,len(emitString))
	}

	//initialize first column
	for i := range dp {
		fmt.Println(transitMap[transitLetter[i]][seq[0]], float64(len(transitLetter)))
		dp[i][0] = math.Log(1/float64(len(transitLetter))*transitMap[transitLetter[i]][seq[0]])
		initNode := node {
			prevX: -1,
			prevY: -1,
			state: transitLetter[i],
		}
		backtrack[i][0] = initNode
	}
	//complete matrix
	for j:=1;j<len(dp[0]); j++ {
		for i:=0; i<len(dp); i++ {
			//iterate through previous column
			weights := make([]float64, 0)
			for k:=0; k<len(dp); k++ {
				transit := stateMap[transitLetter[k]][transitLetter[i]]
				emit := transitMap[transitLetter[i]][seq[j]]

				tmp := dp[k][j-1] + math.Log(transit * emit)
				weights = append(weights, tmp)
			}
			pos, grade := max(weights) //returns the row of the prevNode
			prev := node {
				prevX: pos,
				prevY: j-1, //at previous column
				state: transitLetter[i], //state of new Node
			}
			dp[i][j] = grade
			backtrack[i][j] = prev
		}
	}

	return dp, backtrack
}

func backTrack(dp [][]float64, backtrack [][]node) string {
	dpT := transpose(dp)
	x,_ := max(dpT[len(dpT)-1]) //row of the last maximum grade
	y := len(backtrack[0])-1
	
	res := make([]string, 0)
	for y!= -1 {
		res = append(res, backtrack[x][y].state)
		tmpx := backtrack[x][y].prevX
		tmpy := backtrack[x][y].prevY
		x = tmpx
		y = tmpy
	}

	tmp := strings.Join(res, "")
	return Reverse(tmp)
}

//global variable
var stateMap = map[string]map[string]float64{}
var transitMap = map[string]map[string]float64{}

func main() {
	inputfile := os.Args[1]
	emitString, transitLetter := parseData(inputfile)

	fmt.Println(transitMap)
	dp, backtrack := viterbi(emitString, transitLetter)
	output := backTrack(dp, backtrack)
	fmt.Println(output)
}


func parseData(inputfile string) (string, []string){
	inputFile, err := os.Open(inputfile)
	check(err)
	scanner := bufio.NewScanner(inputFile)

	i := 0
	var emitString string
	emitLetter := make([]string, 0)
	transitLetter := make([]string, 0)

	for scanner.Scan() {
		if scanner.Text() == "--------" {
			i++
			continue
		}
		if i == 0 {
			emitString = scanner.Text()
		} else if i == 1 {
			tmp := strings.Fields(scanner.Text())
			for _, val := range tmp {
				emitLetter = append(emitLetter, val)
			}
		} else if i == 2 {
			tmp := strings.Fields(scanner.Text())
			for _, val := range tmp {
				transitLetter = append(transitLetter, val)
			}

		} else if i == 3 {
			tmp := strings.Fields(strings.TrimSpace(scanner.Text()))
			tmpstate := tmp[0]
			tmp = tmp[1:]
			tmpMap := make(map[string]float64, 0)
			for i, val := range tmp {
				conv, err := strconv.ParseFloat(val, 32)
				check(err)
				tmpMap[transitLetter[i]] = conv
			}
			stateMap[tmpstate] = tmpMap
			continue

		} else if i == 4 {
			tmp := strings.Fields(strings.TrimSpace(scanner.Text()))
			tmpstate := tmp[0]
			tmp = tmp[1:]
			tmpMap := make(map[string]float64, 0)
			for i, val := range tmp {
				conv, err := strconv.ParseFloat(val, 64)
				check(err)
				tmpMap[emitLetter[i]] = conv
			}
			transitMap[tmpstate] = tmpMap
			continue

		}

	}
	return emitString, transitLetter
}

//utility
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func max(input []float64) (int, float64) {
	res := input[0]
	pos := 0
	for i, val := range input {
		if val > res {
			res = val
			pos = i 
		}
	}
	return pos, res
}

func transpose(slice [][]float64) [][]float64 {
    xl := len(slice[0])
    yl := len(slice)
    result := make([][]float64, xl)
    for i := range result {
        result[i] = make([]float64, yl)
    }
    for i := 0; i < xl; i++ {
        for j := 0; j < yl; j++ {
            result[i][j] = slice[j][i]
        }
    }
    return result
}
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}