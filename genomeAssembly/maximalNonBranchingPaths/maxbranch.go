package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
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

*/

func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	//read through input and generate egde map
	scanner := bufio.NewScanner(inputFile)

	nodemap := make(map[string][]string, 0) //node with edges

	leftlist := make([]string, 0) //all nodes
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " -> ")
		left := line[0]
		leftlist = append(leftlist, left)

		rightstring := strings.Split(line[1], ",")
		for _, val := range rightstring {
			nodemap[left] = append(nodemap[left], val)
		}
	}

	//nodeList access node by position
	//nodeLocation access node by value
	nodeList := make([]*node, 0)
	nodeLocation := make(map[string]*node, 0)
	for _, val := range leftlist {
		newNode := node{
			value: val,
		}
		nodeList = append(nodeList, &newNode)
		nodeLocation[val] = &newNode
	}
	//initate graph and adds nodes
	inputGraph := graph{
		Nodes: nodeList,
	}

	//completes the graph by adding egdes
	for _, val := range inputGraph.Nodes {
		for _, next := range nodemap[val.value] {
			if nodeLocation[next] != nil {
				inputGraph.addEdge(val, nodeLocation[next])
			} else {
				finNode := node{
					value: next,
				}
				inputGraph.Nodes = append(inputGraph.Nodes, &finNode)
				inputGraph.addEdge(val, &finNode)
			}
		}
	}

	//inputGraph.printGraph()
	res := nonbranchPath(inputGraph)
	//fmt.Println(res)
	printSequence(res)
}

func printSequence(paths []string)  {
	res := make([]string, 0)
	for _, path := range paths {
		res = append(res, assembleSingle(path))
	}
	fmt.Println(res)

	f, err := os.Create("output.txt")
	check(err)
	defer f.Close()
	for _, val := range res {
		fmt.Fprintln(f, val)
	}
}

func assembleSingle(seq string) string {
	subseq := strings.Split(seq, "->")
	sts := subseq[0]
	for i:=1; i< len(subseq);i++ {
		tmp := strings.Split(subseq[i], "")
		sts = sts+ tmp[len(tmp)-1]
	}
	return sts
}

