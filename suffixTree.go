package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*
ModifiedSuffixTrieConstruction(Text)
    Trie ← a graph consisting of a single node root
    for i ← 0 to |Text| - 1
        currentNode ← root
        for j ← i to |Text| - 1
            currentSymbol ← j-th symbol of Text
            if there is an outgoing edge from currentNode labeled by currentSymbol
                currentNode ← ending node of this edge
            else
                add a new node newNode to Trie
                add an edge newEdge connecting currentNode to newNode in Trie
                Symbol(newEdge) ← currentSymbol
                Position(newEdge) ← j
                currentNode ← newNode
        if currentNode is a leaf in Trie
            assign label i to this leaf
    return Trie

MaximalNonBranchingPaths(Graph)
    Paths ← empty list
    for each node v in Graph
        if v is not a 1-in-1-out node
            if out(v) > 0
                for each outgoing edge (v, w) from v
                    NonBranchingPath ← the path consisting of single edge (v, w)
                    while w is a 1-in-1-out node
                        extend NonBranchingPath by the edge (w, u)
                        w ← u
                    add NonBranchingPath to the set Paths
    for each isolated cycle Cycle in Graph
        add Cycle to Paths
    return Paths

ModifiedSuffixTreeConstruction(Text)
    Trie ← ModifiedSuffixTrieConstruction
    for each non-branching path Path in Trie
        substitute Path by a single edge e connecting the first and last nodes of Path
        Position(e) ← Position(first edge of Path)
        Length(e) ← number of edges of Path
    return Trie
*/

//Edge connect 2 nodes, position of symbol in TEXT.
type edge struct {
	st       *node
	end      *node
	symbol   string
	position int
}

//node is connected to next node by edges
type node struct {
	self int
	next []*edge
}

/*
Engine.
Based on the input string, first built a tree of all suffix
Find all max. nonbranching path. Print as an substring
*/
func suffixTrie(input string) ([]edge, []*node) {
	nodes := make([]*node, 0) //store pointer to node
	root := node{self: 0}
	nodes = append(nodes, &root)
	edges := make([]edge, 0)

	splitInput := strings.Split(input, "")
	for i := 0; i < len(splitInput); i++ {
		currentN := nodes[0]
		for j := i; j < len(splitInput); j++ {
			currentS := splitInput[j]
			found, nextNode := outEdgeSymbol(currentN, currentS)
			if found {
				currentN = nextNode
			} else {
				newNode := node{self: len(nodes)}
				nodes = append(nodes, &newNode)
				newEdge := edge{
					st:       currentN,
					end:      &newNode,
					symbol:   currentS,
					position: j,
				}
				edges = append(edges, newEdge)                  //append edge to graph
				currentN.next = append(currentN.next, &newEdge) //connect nodes
				currentN = &newNode
			}
		}
	}

	//path := nonBranchPath(edges, nodes)
	//printPath(path, splitInput)

	return edges, nodes
}

/*
TreeColoring(ColoredTree)
    while ColoredTree has ripe nodes
        for each ripe node v in ColoredTree
            if there exist differently colored children of v
                Color(v) ← "purple"
            else
                Color(v) ← color of all children of v
    return ColoredTree
*/


func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var inputSeq string
	for scanner.Scan() {
		inputSeq = scanner.Text() //input is in one line
	}

	//mergedTree(inputSeq)
	suffixArray(inputSeq)
	/*
	egdes, nodes := suffixTrie(inputSeq)
	path := nonBranchPath(egdes, nodes)
	printPath(path, strings.Split(inputSeq, ""))
	*/
	//countleaves(inputSeq)

}

//output all nonbranching path from graph
func nonBranchPath(graph []edge, nodes []*node) [][]int {
	path := make([][]int, 0)
	for _, v := range nodes {
		if len(v.next) > 1 {
			//initiate path if not 1-1
			for _, v_edge := range v.next {
				tmpPath := make([]int, 0)
				tmpPath = append(tmpPath, v_edge.position)
				//elongate the path
				w := v_edge.end
				for len(w.next) == 1 {
					u := w.next[0]
					tmpPath = append(tmpPath, u.position)
					w = u.end
				}
				path = append(path, tmpPath)
			}
		}
	}

	return path //path stores the position of substrings in TEXT
}

func mergedTree(input string) {
	input = input + "$"
	_, nodes := suffixTrie(input)

	spString := strings.Split(input, "")

	for _, v := range nodes {
		if len(v.next) > 1 {
			tmpEdges := make([]*edge, 0)
			for _, v_edge := range v.next {
				tmpPath := make([]int, 0)
				tmpPath = append(tmpPath, v_edge.position)
				w := v_edge.end
				for len(w.next) == 1 {
					u := w.next[0]
					tmpPath = append(tmpPath, u.position)
					w = u.end
				}

				//adding to new graph
				mSymbol := spString[tmpPath[0] : tmpPath[0]+len(tmpPath)]
				edgeIJ := edge{st: v, end: w, symbol: strings.Join(mSymbol, "")}
				tmpEdges = append(tmpEdges, &edgeIJ)
			}
			v.next = nil 
			v.next = tmpEdges
		}
	}


	//searchTree(nodes[0])
	maxWord = append(maxWord, "")
	longestRepeat(nodes[0], "")
	fmt.Println("max return is", maxWord[len(maxWord)-1])
}

//used to find longest repeat from one input
var maxWord []string
func longestRepeat(root *node, word string)  {
	lln := len(maxWord)-1
	if len(root.next) > 1 && len(word) > len(maxWord[lln]) {
		//fmt.Println("comparing", word,maxWord[lln] )
		maxWord = append(maxWord, word)
		//fmt.Println("in", maxWord)
	}
	for _, e := range root.next {
		longestRepeat(e.end, word+e.symbol)
	}

}


func countleaves(input string) {
	_, nodes :=  suffixTrie(input)
	c = 0
	count(nodes[0])
	fmt.Println(c)
}
var c int
func count(root *node)  {
	if len(root.next) == 0 {
		c++
	} else {
		for _, e := range root.next {
			count(e.end)
		}
	}
}

func outEdgeSymbol(n *node, s string) (bool, *node) {
	for _, oneedge := range n.next {
		if oneedge.symbol == s {
			return true, oneedge.end
		}
	}
	return false, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func searchTree(root *node) {
	if root.next == nil {
		fmt.Println(root.self, "is leaf")
	}
	fmt.Println(root.self)
	for _, e := range root.next {
		fmt.Println("searching", e.end.self)
		searchTree(e.end)
	}

}

//output all substrings
func printPath(path [][]int, text []string) {
	f, err := os.Create("out.txt")
	check(err)
	defer f.Close()

	for _, sPath := range path {
		syn := strings.Join(text[sPath[0]:sPath[0]+len(sPath)], "")
		fmt.Fprint(f, syn, " ")
	}

}

func suffixArray(input string) {
	suffix := make([]string, 0)
	pos := make(map[string]int, 0)
	sinput := strings.Split(input, "")
	for i:=0; i <len(sinput); i++ {
		aSuffix := strings.Join(sinput[i:], "")
		suffix = append(suffix, aSuffix)
		pos[aSuffix] = i
	}
	sort.Strings(suffix)

	f, err := os.Create("array.txt")
	check(err)
	defer f.Close()

	for _, val := range suffix {
		fmt.Fprint(f, pos[val], " ")
	}

}
