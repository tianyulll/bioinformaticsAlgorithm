package main

import (
	"fmt"
	"math"
	"strings"
)

func roundMatrix(matrix [][]float64) [][]float64 {
	for i := range matrix {
		for j, val := range matrix[i] {
			matrix[i][j] = math.Round(val*100) /100
		}
	}
	return matrix

}

func newProfile(prob [][]float64, motifs [][]string, len int) [][]float64 {
	result := make([][]float64, 4)
	for i := range result {
		result[i] = make([]float64, len)
	}
	for i := range motifs {
		for j := range motifs[i] {
			chars := strings.Split(motifs[i][j], "")
			for k, char := range chars {
				result[position(char)][k] += prob[i][j]
			}
		}
	}
	for i := range result {
		for j, val := range result[i] {
			result[i][j] = val / float64(len)
		}
	}
	return result
}

//ATGC
func position(seq string) int {
	if seq == "A" {
		return 0
	}
	if seq == "T" {
		return 1
	}
	if seq == "G" {
		return 2
	}
	if seq == "C" {
		return 3
	}
	return 0
}

func main() {
	sequence := []string{"TACATCCTG", "GGACAGCGT", "GACTTTGGC", "AACAGTGAC"}
	motifs := []string{"CCTG", "ACAG", "TTGG", "CAGT"}
	tmpmotifs := motifs[1:] //neglect line 1
	profile := buildProfile(tmpmotifs, 4)
	//fmt.Println(profile)
	chunks := make([][]string, 0)
	for _, seq := range sequence {
		chunks = append(chunks, chop(seq, 4))
	}

	grades := make([][]float64, 0)
	for _, chunk := range chunks {
		lineGrade := make([]float64, 0)
		for _, mot := range chunk {
			lineGrade = append(lineGrade, gradeMotif(mot, profile))
		}
		grades = append(grades, lineGrade)
	}

	//fmt.Println(grades)
	for i, seqGrade := range grades {
		total := 0.0
		for _, grade := range seqGrade {
			total += grade
		}
		//fmt.Println("total is", total)
		for j, grade := range seqGrade {
			grades[i][j] = grade / total //math.Round(grade/total*100) /100
		}
	}
	newpro := newProfile(grades, chunks, 4)
	fmt.Println("chunk")
	fmt.Println(chunks)
	fmt.Println("grade")
	fmt.Println(roundMatrix(grades))
	fmt.Println("Newprofile")
	fmt.Println(roundMatrix(newpro))

}

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
