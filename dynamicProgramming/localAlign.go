package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readScoreMatrix(scorematrix string) map[string]map[string]int {
	inputFile, err := os.Open(scorematrix)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	j := 0
	aminoKeys := make([]string, 0)
	finalScore := make(map[string]map[string]int,0)
	for scanner.Scan() {
		//set up amino acid keys
		line := strings.Fields(scanner.Text())
		if j == 0 {
			for _, val := range line {
				aminoKeys = append(aminoKeys, val)
			}
			j++
			continue
		}
		aa := line[0]
		tmpMap := make(map[string]int)
		for i:=1;i<len(line);i++ {
			tmp, err := strconv.Atoi(line[i])
			check(err)
			tmpMap[aminoKeys[i-1]] = tmp 
		}
		finalScore[aa] = tmpMap
	}
	return finalScore
}

func buildScoreMatrix(str1, str2 string, grade map[string]map[string]int) [][]int {
	score := make([][]int, len(str2)+1)
	for i := 0; i < len(str2)+1; i++ {
		score[i] = make([]int, len(str1)+1)
	}
	seq1 := strings.Split(str1, "")
	seq2 := strings.Split(str2, "")
	
	for i := 1; i <= len(seq2); i++ {
		for j := 1; j <= len(seq1); j++ {
			left := seq2[i-1]
			right := seq1[j-1]	
			score[i][j] = max([]int{score[i-1][j-1] + grade[left][right], score[i-1][j]-5, score[i][j-1]-5, 0})
		}
	}
	//fmt.Println(score)
	return score
}

func buildAligin(str1, str2 string, i,j int, score[][]int, grade map[string]map[string]int) ([]string, []string) {
	seq1 := strings.Split(str1, "")
	seq2 := strings.Split(str2, "")
	aln1 := make([]string, 0)
	aln2 := make([]string, 0)

	for i!=0 && j!=0 {
		if score[i][j]==0 {
			break
		}
		if seq1[j-1] == seq2[i-1]  {
			aln1 = append(aln1, seq1[j-1])
			aln2 = append(aln2, seq2[i-1])
			i--
			j--
		} else if score[i-1][j-1]+grade[seq2[i-1]][seq1[j-1]] == score[i][j] {
			aln1 = append(aln1, seq1[j-1])
			aln2 = append(aln2, seq2[i-1])
			i--
			j--
		} else if score[i-1][j]-5 == score[i][j] {
			aln2 = append(aln2, seq2[i-1])
			aln1 = append(aln1, "-")
			i--
		} else if score[i][j-1]-5 == score[i][j] {
			aln1 = append(aln1, seq1[j-1])
			aln2 = append(aln2, "-")
			j--
		} 
	}
	return aln1, aln2
}

func findBestScore(score [][]int) (int, int, int) {
	besti := 0
	bestj := 0
	best := score[besti][bestj]
	for i:=0;i<len(score);i++ {
		for j:=0;j<len(score[i]);j++ {
			if score[i][j] > best {
				best = score[i][j]
				besti = i
				bestj = j
			}
		}
	}
	return best, besti, bestj
}

func main() {
	scorereference := "bscorematrix.txt"
	scorematrix := readScoreMatrix(scorereference)

	//read the input sequence
	inputfile := os.Args[1]
	inputFile, err := os.Open(inputfile)
	check(err)
	scanner := bufio.NewScanner(inputFile)
	i:=0
	var str1, str2 string 
	for scanner.Scan() {
		if i == 0 {
			str1 = scanner.Text()
			i++ 
		} else {
			str2 = scanner.Text()
		}
	}

	scores := buildScoreMatrix(str1, str2, scorematrix)
	best, besti, bestj := findBestScore(scores)
	fmt.Println("best score", best, besti, bestj)

	aln1, aln2 := buildAligin(str1, str2, besti, bestj, scores, scorematrix)
	printAlign(aln1, aln2)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func max(input []int) int {
	res := input[0]
	for _, val := range input {
		if val > res {
			res = val
		}
	}
	return res
}

func printAlign(aln1, aln2 []string) {
	for i, j := 0, len(aln1)-1; i < j; i, j = i+1, j-1 {
		aln1[i], aln1[j] = aln1[j], aln1[i]
	}
	fmt.Println(strings.Join(aln1, ""))
	for i, j := 0, len(aln2)-1; i < j; i, j = i+1, j-1 {
		aln2[i], aln2[j] = aln2[j], aln2[i]
	}
	fmt.Println(strings.Join(aln2, ""))

}
