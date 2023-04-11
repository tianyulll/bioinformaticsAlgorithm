package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//numberT = dna number
func gibbsSampler(dna []string, lengthK, numberT, timeN int) []string {
	bestmotifs := make([]string, numberT)
	for i, single := range dna {
		bestmotifs[i] = randomMotif(single, lengthK, numberT)
	}
	for n := 0; n < timeN; n++ {
		motifMatrix, i := randomSelectMotif(bestmotifs)
		popMotif := bestmotifs[i]
		popSequence := dna[i]
		//i stores the position of the neglected sequence

		profile := buildProfile(motifMatrix, lengthK)
		popMotGrade := gradeMotif(popMotif, profile)
		//generate new random motif
		newMotif, newMotifGrade := randomKmerofSeq(popSequence, lengthK, profile)

		//need to compare the random motif and original motif
		if newMotifGrade > popMotGrade {
			bestmotifs[i] = newMotif
		}

	}
	return bestmotifs
}

func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i, lenghK, numberT, timeN := 0, 0, 0, 0 //k == length, n == number
	var sequence []string
	for scanner.Scan() {
		if i == 0 {
			tmp := strings.Split(scanner.Text(), " ")
			lenghK, err = strconv.Atoi(tmp[0])
			numberT, err = strconv.Atoi(tmp[1])
			timeN, err = strconv.Atoi(tmp[2])
			i++
			continue
		}
		sequence = append(sequence, scanner.Text())
	}

	bestMotifs := gibbsSampler(sequence, lenghK, numberT, timeN)

	fmt.Println(bestMotifs)
	f, err := os.Create("output.txt")
	check(err)
	defer f.Close()
	for _, line := range bestMotifs {
		fmt.Fprintln(f, line)
	}
}
