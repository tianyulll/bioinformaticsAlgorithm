package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func graph(kmers []string) map[string][]string {
	res :=make(map[string][]string)
	for _, kmer := range kmers {
		res[getPrefix(kmer)] = append(res[getPrefix(kmer)], getSuffix(kmer))
	}
	return res
}


func getPrefix(kmer string) string {
	char := strings.Split(kmer, "")
	tmp := char[0:len(char)-1]
	return strings.Join(tmp, "")
}

func getSuffix(kmer string) string{
	char := strings.Split(kmer, "")
	tmp := char[1:]
	return strings.Join(tmp, "")
}

func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var sequence []string
	for scanner.Scan() {
		sequence = append(sequence, scanner.Text())
	}
	adjList := graph(sequence)
	fmt.Println(adjList)

	for key, val := range adjList {
		fmt.Println(key+" -> "+ strings.Join(val, ","))
	}
	//print to file
	f, err := os.Create("output.txt")
	check(err)
	defer f.Close()
	for key, val := range adjList {
		fmt.Fprintln(f, key+" -> "+ strings.Join(val, ","))
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
