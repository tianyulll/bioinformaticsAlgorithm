package main

import (
	"math/rand"
	"strings"
	"time"
)

func chop(s string, size int) []string {
	var chunks []string

	for i := range s {
		if i+size > len(s) {
			break
		}
		chunks = append(chunks, s[i:i+size])
	}
	return chunks
}

func gradeMotif(single string, profile [][]float64) float64 {
	seq := strings.Split(single, "")
	grade := 1.0
	for i, nucleotide := range seq {
		if nucleotide == "A" {
			grade *= profile[0][i]
		}
		if nucleotide == "C" {
			grade *= profile[1][i]
		}
		if nucleotide == "G" {
			grade *= profile[2][i]
		}
		if nucleotide == "T" {
			grade *= profile[3][i]
		}
	}
	return grade
}

//k = length, n = number
func randomMotif(seq string, k, n int) string {
	chunks := chop(seq, k)
	rand.Seed(time.Now().UnixNano())
	return chunks[rand.Intn(len(chunks))]
}

//randomly neglect one of the string
//returns motif for profile 
//and the popped out sequence number
func randomSelectMotif(motifs []string) ([]string, int) {
	motifMatrix := make([]string, len(motifs))

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(motifs))
	for j := 0; j < len(motifs); j++ {
		if j == i {
			continue
		}
		motifMatrix = append(motifMatrix, motifs[j])
	}
	
	return motifMatrix, i
}

//generate a random motif from sequence based on profile
func randomKmerofSeq(seq string, lengthK int, profile [][]float64) (string, float64) {
	chunks := chop(seq, lengthK)
	chunkGrade := make([]float64, len(chunks))
	totalGrade := 0.0
	for i, motif := range chunks {
		chunkGrade[i] = gradeMotif(motif, profile)
		totalGrade += chunkGrade[i]
	}
	for j, grade := range chunkGrade {
		chunkGrade[j] = grade / totalGrade
	}
	res := weightedRandomSel(chunkGrade, totalGrade)
	return chunks[res], chunkGrade[res] 
}

//need improvement
//weighted random selection based on probability of the motif collection
func weightedRandomSel(motifs []float64, total float64) int {
	rand.Seed(time.Now().UnixNano())
	//rand.Shuffle(len(motifs), func(i, j int)  { motifs[i], motifs[j] = motifs[j], motifs[i] })
	r := rand.Float64()*total
	for i, weight := range motifs {
		r-= weight 
		if r < 0 {
			return i
		}
	}
	return len(motifs) - 1
}

//build a profile based on motif
//dim - length k
func buildProfile(motif []string, dim int) [][]float64 {
	profile := make([][]float64, 4)
	for j := 0; j < 4; j++ {
		profile[j] = make([]float64, dim)
	}
	for j := 0; j < 4; j++ {
		for k := 0; k < dim; k++ {
			profile[j][k] = 1.0
		}
	}
	for _, s := range motif {
		seq := strings.Split(s, "")
		for i, n := range seq {
			if n == "A" {
				profile[0][i]++
			}
			if n == "C" {
				profile[1][i]++
			}
			if n == "G" {
				profile[2][i]++
			}
			if n == "T" {
				profile[3][i]++
			}
		}
	}
	colsum := 0.0
	for j := 0; j < 4; j++ {
		colsum += profile[j][0]
	}
	for j := 0; j < 4; j++ {
		for k := 0; k < dim; k++ {
			profile[j][k] = profile[j][k] / colsum
		}
	}
	return profile
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
