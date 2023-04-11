package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//each node in matrix should have three layers
type grade struct {
	low    int
	middle int
	high   int
}

//Return 3 layer score matrix based on the sequence input
func scoreMatrix(str1, str2 string, match, mismatch, gapOpen, gapExtend int) ([][]grade, [][]string) {
	scores := make([][]grade, len(str2)+1) //records 3 layer score
	for i := 0; i < len(str2)+1; i++ {
		scores[i] = make([]grade, len(str1)+1)
	}
	backtrack := make([][]string, len(scores)) //decisions for backtrack
	for i := 0; i < len(scores); i++ {
		backtrack[i] = make([]string, len(scores[0]))
	}

	seq1 := strings.Split(str1, "")
	seq2 := strings.Split(str2, "")

	//intialize middle
	for i:=0; i< len(scores);i++ {
		scores[i][0].middle = - i*gapExtend
		scores[i][0].high = -i*gapOpen
	} 
	for j := 0; j < len(scores[0]); j++ {
		scores[0][j].middle = -j * gapExtend
		scores[0][j].low = -j*gapOpen
	}
	scores[0][0].middle = 0

	//lower level: vertical, deletion
	//higher level: horizontal, insertion
	for i := 1; i <= len(seq2); i++ {
		for j := 1; j <= len(seq1); j++ {
			scores[i][j].low = max([]int{scores[i-1][j].low - gapExtend, scores[i-1][j].middle - gapOpen})
			scores[i][j].high = max([]int{scores[i][j-1].high - gapExtend, scores[i][j-1].middle - gapOpen})

			matchscore := matchScore(seq1[j-1], seq2[i-1], match, mismatch)
			scores[i][j].middle = max([]int{scores[i][j].low,
				scores[i-1][j-1].middle + matchscore,
				scores[i][j].high})

			if scores[i][j].middle == scores[i-1][j-1].middle+matchscore {
				backtrack[i][j] = "down"
			} else if scores[i][j].middle == scores[i][j].low {
				backtrack[i][j] = "del"
			} else if scores[i][j].middle == scores[i][j].high {
				backtrack[i][j] = "ins"
			}
		}
	}

	return scores, backtrack
}

func buildAlign(str1, str2 string, backtrack [][]string) ([]string, []string) {
	seq1 := strings.Split(str1, "")
	seq2 := strings.Split(str2, "")
	aln1 := make([]string, 0)
	aln2 := make([]string, 0)
	i := len(seq2) //used for aln2, seq2
	j := len(seq1) //used for aln1, seq1
	for i != 0 && j != 0 {
		if backtrack[i][j] == "down" {
			aln1 = append(aln1, seq1[j-1])
			aln2 = append(aln2, seq2[i-1])
			i--
			j--
		} else if backtrack[i][j] == "del" {
			aln2 = append(aln2, seq2[i-1])
			aln1 = append(aln1, "-")
			i--
		} else if backtrack[i][j] == "ins" {
			aln1 = append(aln1, seq1[j-1])
			aln2 = append(aln2, "-")
			j--
		}
	}
	return aln1, aln2
}

//utility for scoring match or mismatch
//inputs are both positive
func matchScore(a, b string, match, mismatch int) int {
	if a == b {
		return match
	} else {
		return -mismatch
	}
}

//utility to print the matrix
func printScoreMatrix(input [][]grade) {
	Middle := make([][]int, len(input))
	low := make([][]int, len(input))
	high := make([][]int, len(input))

	for i := 0; i < len(input); i++ {
		Middle[i] = make([]int, len(input[0]))
		low[i] = make([]int, len(input[0]))
		high[i] = make([]int, len(input[0]))
	}
	for i := 0; i < len(input); i++ {
		for j := 0; j < len(input[0]); j++ {
			Middle[i][j] = input[i][j].middle
			low[i][j] = input[i][j].low
			high[i][j] = input[i][j].high

		}
	}
	fmt.Println(Middle)
	fmt.Println(low)
	fmt.Println(high)
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

/*
Input: A match reward, a mismatch penalty, a gap opening penalty, a gap extension penalty, and two nucleotide strings.
Output: The maximum alignment score between v and w, followed by an alignment of v and w achieving this maximum score.
1 3 2 1
GA
GTTA
*/
func main() {
	inputfile := os.Args[1]
	inputFile, err := os.Open(inputfile)
	check(err)
	scanner := bufio.NewScanner(inputFile)
	i := 0
	var match, mismatch, gapOpen, gapExtend int
	var str1, str2 string
	for scanner.Scan() {
		if i == 0 {
			gdval := strings.Split(scanner.Text(), " ")
			match, err = strconv.Atoi(gdval[0])
			mismatch, err = strconv.Atoi(gdval[1])
			gapOpen, err = strconv.Atoi(gdval[2])
			gapExtend, err = strconv.Atoi(gdval[3])
			check(err)
			i++
		} else if i == 1 {
			str1 = scanner.Text()
			i++
		} else if i == 2 {
			str2 = scanner.Text()
		}
	}

	scores, backtrack := scoreMatrix(str1, str2, match, mismatch, gapOpen, gapExtend)
	aln1, aln2 := buildAlign(str1, str2, backtrack)
	//fmt.Println(scores[len(scores)-1])
	fmt.Println(scores[len(scores)-1][len(scores[0])-1].middle)
	printAlign(aln1, aln2)
}

//utility
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
