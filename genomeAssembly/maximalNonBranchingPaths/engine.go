package main

import "fmt"

type node struct {
	value string
	prev  []*node
	next  []*node
}

type graph struct {
	Nodes []*node
}

//edge directing from start to end
func (this *graph) addEdge(start, end *node) {
	start.next = append(start.next, end)
	end.prev = append(end.prev, start)
}

func (this *node) inDegree() int {
	return len(this.prev)
}

func (this *node) outDegree() int {
	return len(this.next)
}

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
//finding all non-branching paths
//still needs the isolated cycles
func nonbranchPath(input graph) []string {
	emptyPath := make([]string, 0)
	for _, thisNode := range input.Nodes {
		if thisNode.inDegree() != 1 ||
			thisNode.outDegree() != 1 {
			if thisNode.outDegree() > 0 {
				for _, nextNode := range thisNode.next {
					localPath := thisNode.value + "->"+nextNode.value
					for nextNode.inDegree() == 1 && nextNode.outDegree() == 1 {
						localPath = localPath + "->"+nextNode.next[0].value
						nextNode = nextNode.next[0]
					}
					emptyPath = append(emptyPath, localPath)
				}
			}
		}
	}
	return emptyPath
}

var globalPath []*node

//supposed to extend the path
func singlePath(v *node, list []*node) []*node {
	for _, edge := range v.next {
		list = append(list, edge)
		return singlePath(edge, list)
	}
	return list
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (this *graph) printGraph() {
	for _, nodes := range this.Nodes {
		fmt.Println("node", nodes.value)
		fmt.Println("indegre", nodes.inDegree())
		for _, edge := range nodes.next {
			fmt.Print("edge", edge.value, " ")
		}
		fmt.Println()
	}
}
