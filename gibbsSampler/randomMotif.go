package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type motif struct {
	sequence string
	grade    float64
}

/*
 RandomizedMotifSearch(Dna, k, t)
       randomly select k-mers Motifs = (Motif1, …, Motift) in each string from Dna
       BestMotifs ← Motifs
       while forever
           Profile ← Profile(Motifs)
           Motifs ← Motifs(Profile, Dna)
           if Score(Motifs) < Score(BestMotifs)
               BestMotifs ← Motifs
           else
               return BestMotifs
*/

func randomizedMotifSearch(dnas []string, lengthK, numberN int) ([]string, float64) {
	bestMotifs := make([]string, numberN)
	//randomly select motifs from DNA collection
	for i, dna := range dnas {
		bestMotifs[i] = randomMotif(dna, lengthK, numberN) 
	}
	//initial grade
	profile := buildProfile(bestMotifs, lengthK)
	bestGrade := 1.0
    for _, mot := range bestMotifs {
		bestGrade *= gradeMotif(mot, profile)
	}
	//start to substitute
	for {
		profile = buildProfile(bestMotifs, lengthK)
		for i, dna := range dnas {
			allMotif := chop(dna, lengthK)
			for _, singleMotif := range allMotif {
				if gradeMotif(singleMotif, profile) > gradeMotif(bestMotifs[i], profile) {
					bestMotifs[i] = singleMotif
				}
			}
		}
		tmpGrade := 1.0
		for _, mot := range bestMotifs {
			tmpGrade *= gradeMotif(mot, profile)
		}
		if tmpGrade == bestGrade {
			grade := 1.0
			for _, mot := range bestMotifs {
				grade *= gradeMotif(mot, profile)
			}
			return bestMotifs, grade
		} else {
			bestGrade = tmpGrade
		}

	}
}





func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i, k, n := 0, 0, 0 //k == length, n == number
	var sequence []string
	for scanner.Scan() {
		if i == 0 {
			tmp := strings.Split(scanner.Text(), " ")
			k, err = strconv.Atoi(tmp[0])
			n, err = strconv.Atoi(tmp[1])
			i++
			continue
		}
		sequence = append(sequence, scanner.Text())
	}

	var bestmotif []string
	grade := 0.0
	for i := 0; i < 3000; i ++ {
		tmpList, tmpGrade := randomizedMotifSearch(sequence, k, n)
		if tmpGrade > grade {
			bestmotif = tmpList
			grade = tmpGrade
		}
	}
	fmt.Println(bestmotif)
	f, err := os.Create("output.txt")
	check(err)
	defer f.Close()
	for _, line := range bestmotif {
		fmt.Fprintln(f, line)
	}
}

/*
func sampleresult(seq []string,k int) float64 {
	sampleOut := []string{"TCTCGGGG","CCAAGGTG","TACAGGCG","TTCAGGTG","TCCACGTG"}
	profile := buildProfile(sampleOut, k)
	grade := 1.0
	for _, s := range sampleOut {
		grade *= gradeMotif(s, profile)
	}
	return grade
}
*/

