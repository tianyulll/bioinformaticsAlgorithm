package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type edge struct {
	x         int
	y         int
	state1    int
	state2    int
	state1Pos int
	state2Pos int
	emit      string
}

//state 1, 2, 3 correspond to insertion, match, deletion
func (this *edge) setCoordinate() {
	switch {
	case this.state1 == 0:
		this.x = 0 //changhe
	case this.state1 == 1:
		this.x = 1 + 3*this.state1Pos
	case this.state1 == 2:
		this.x = 2 + 3*(this.state1Pos-1)
	case this.state1 == 3:
		this.x = 3 + 3*(this.state1Pos-1)
	}
	switch {
	case this.state2 == 0:
		this.y = 0
	case this.state2 == 1:
		this.y = 1 + 3*this.state2Pos
	case this.state2 == 2:
		this.y = 2 + 3*(this.state2Pos-1)
	case this.state2 == 3:
		this.y = 3 + 3*(this.state2Pos-1)
	}
}

func hmm(align [][]string, cols []int) ([][]edge, int) {
	seqPath := make([][]edge, 0)
	dim := 3*(len(align[0])-len(cols)) + 3

	for i := 0; i < len(align); i++ {
		path := make([]edge, 0)

		//handle source->first column
		isIns := contains(0, cols)
		isDel, _ := delDetect(align[i][0], align[i][0])
		source := edge{state1: 0, state1Pos: 0}
		if isIns && !isDel {
			source.state2 = 1
			source.state2Pos = 0
		} else if isIns && isDel {
			source.state2 = 0
			source.state2Pos = 0
		} else if !isIns && isDel {
			source.state2 = 3
			source.state2Pos = 1
		} else {
			source.state2 = 2
			source.state2Pos = 1
		}
		source.setCoordinate()
		if !isDel {
			source.emit = align[i][0]
		}
		path = append(path, source)

		for j := 0; j < len(align[0])-1; j++ {
			nextIns := contains(j+1, cols)
			_, nextDel := delDetect(align[i][j], align[i][j+1])

			if nextIns && nextDel {
				continue
			}
			preState := path[len(path)-1].state2
			preState2Pos := path[len(path)-1].state2Pos
			newEdge := edge{state1: preState, state1Pos: preState2Pos}

			if nextIns && !nextDel {
				newEdge.state2 = 1
				newEdge.state2Pos = newEdge.state1Pos
			} else if !nextIns && nextDel {
				newEdge.state2 = 3
				newEdge.state2Pos = newEdge.state1Pos + 1
			} else {
				newEdge.state2 = 2
				newEdge.state2Pos = newEdge.state1Pos + 1
			}
			newEdge.setCoordinate()
			if !nextDel {
				newEdge.emit = align[i][j+1]
			}

			path = append(path, newEdge)
		}
		seqPath = append(seqPath, path)
	}
	return seqPath, dim
}

//normalize and print state matrix
func stateMatrix(list [][]edge, dim int) {
	m := make([][]float64, dim)
	for i := range m {
		m[i] = make([]float64, dim)
	}
	for _, line := range list {
		for i, val := range line {
			m[val.x][val.y] += 1
			if i == len(line)-1 {
				m[val.y][dim-1] += 1
			}
		}
	}
	mNormal := normalizeMatrix(m)
	fmt.Println("state Matrix")
	for _, line := range mNormal {
		fmt.Println(arrayToString(line, ","))
	}
}

func emitMatrix(list [][]edge, letter []string, dim int) {
	m := make([][]float64, dim)
	for i := range m {
		m[i] = make([]float64, len(letter))
	}
	for _, line := range list {
		for _, val := range line {
			emitpos := val.emitPos(letter)
			if emitpos == -1 {
				continue
			}
			m[val.y][emitpos] += 1
		}
	}
	mNormal := normalizeMatrix(m)
	fmt.Println("emission matrix")
	for _, line := range mNormal {
		fmt.Println(arrayToString(line, ","))
	}

}

func (this *edge) emitPos(letter []string) int {
	for i, l := range letter {
		if this.emit == l {
			return i
		}
	}
	return -1
}

//return index of columns that are above insertion threshold
func insertCol(align []string, thre float64) ([]int, [][]string) {
	mat := make([][]string, len(align))
	res := make([]int, 0)

	for i, val := range align {
		mat[i] = strings.Split(val, "")
	}
	matT := transpose(mat)
	for i, row := range matT {
		count := 0.0
		for _, val := range row {
			if val == "-" {
				count++
			}
		}
		if (count / float64(len(row))) > thre {
			res = append(res, i)
		}
	}
	return res, mat
}

func delDetect(str1, str2 string) (bool, bool) {
	this := false
	that := false
	if str1 == "-" {
		this = true
	}
	if str2 == "-" {
		that = true
	}
	return this, that
}

//read the input data
func parseData(inputfile string) ([]string, float64, []string) {
	inputFile, err := os.Open(inputfile)
	check(err)
	scanner := bufio.NewScanner(inputFile)

	var dict []string
	align := make([]string, 0)
	thre := 0.0

	i := 0
	for scanner.Scan() {
		if scanner.Text() == "--------" {
			i++
			continue
		}
		if i == 0 {
			thre, err = strconv.ParseFloat(scanner.Text(), 64)
		}
		if i == 1 {
			dict = strings.Fields(scanner.Text())
		}
		if i == 2 {
			align = append(align, scanner.Text())
		}
	}
	return align, thre, dict
}

func main() {
	inputfile := os.Args[1]
	align, thre, dict := parseData(inputfile)
	inCol, alignSeq := insertCol(align, thre)

	list, dim := hmm(alignSeq, inCol)
	stateMatrix(list, dim)
	emitMatrix(list, dict, dim)
}

//utility
func transpose(slice [][]string) [][]string {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]string, xl)
	for i := range result {
		result[i] = make([]string, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func contains(input int, list []int) bool {
	for _, val := range list {
		if input == val {
			return true
		}
	}
	return false
}

func sum(array []float64) float64 {
	result := 0.0
	for _, v := range array {
		result += v
	}
	return result
}

func printPath(list [][]edge, dim int) {
	m := make([][]int, dim)
	for i := range m {
		m[i] = make([]int, dim)
	}
	for _, line := range list {
		for i, val := range line {
			m[val.x][val.y] += 1
			if i == len(line)-1 {
				m[val.y][dim-1] += 1
			}
		}
	}
	for _, line := range m {
		fmt.Println(line)
	}
}

func normalizeMatrix(m [][]float64) [][]float64 {
	for i, line := range m {
		sumLine := sum(line)
		if sumLine == 0 {
			continue
		}
		for j, val := range line {
			tmp := val / float64(sumLine)
			m[i][j] = math.Round(tmp*1000) / 1000
		}
	}
	return m
}

func arrayToString(a []float64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}
