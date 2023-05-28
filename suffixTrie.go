package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
TrieConstruction(Patterns)
    Trie ← a graph consisting of a single node root
    for each string Pattern in Patterns
        currentNode ← root
        for i ← 0 to |Pattern| - 1
            currentSymbol ← Pattern[i]
            if there is an outgoing edge from currentNode with label currentSymbol
                currentNode ← ending node of this edge
            else
                add a new node newNode to Trie
                add a new edge from currentNode to newNode with label currentSymbol
                currentNode ← newNode
    return Trie
*/
func triesConstruct(patterns []string) []edge {
	nodes := make([]int, 0)
	nodes = append(nodes, len(nodes))
	edges := make([]edge, 0)

	for _, pattern := range patterns {
		currentN := nodes[0]
		sPattern := strings.Split(pattern, "")
		for i := range sPattern {
			currentS := sPattern[i]
			tryEnd := findEdge(edges, currentN, currentS)
			if tryEnd != -1 {
				currentN = tryEnd
			} else {
				nodes = append(nodes, len(nodes))
				newEdge := edge {
					st: currentN,
					end: nodes[len(nodes)-1],
					symbol: currentS,
				}
				edges = append(edges, newEdge)
				currentN = newEdge.end
			}
		}
	}

	return edges
}

type edge struct {
	st int
	end int
	symbol string
}

func findEdge(edges []edge, start int, sym string) int {
	for _, oneedge := range edges {
		if oneedge.st == start {
			if oneedge.symbol == sym {
				return oneedge.end
			}
		}
	}
	return -1
}

func printGraph(edges []edge) {
	f, err := os.Create("out.txt")
	check(err)
	defer f.Close()
	
	for _, oneEdge := range edges {
		fmt.Fprintln(f, oneEdge.st, oneEdge.end, oneEdge.symbol)
	}
}

func main() {
	
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var inputSeq string
	for scanner.Scan() {
		inputSeq = scanner.Text()
	}
	seqs := strings.Split(inputSeq, " ")
	edges := triesConstruct(seqs)
	printGraph(edges)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
