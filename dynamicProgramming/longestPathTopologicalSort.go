package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type edge struct {
	right  int
	weight int
}

/*
Based on topological algorithm,
add a longestPath function.
First, sort the graph topologically
Initialize dist[v] = infinity for all nodes.
Set dist[s] to 0.
Keep a queue of nodes and start with the roots.
Pop  the  front  of  the  queue,v,
and  iterate  its  neighbors,u.
If dist[v] +d[v,u]>dust[u],
then  setdist[u] =dist[v]+d[v,u].
*/

func longpath(input map[int][]edge, len int) [][]int {
	paths := make([]int, len+1)
	for i := 0; i <= len; i++ {
		paths[i] = math.MinInt
	}
	paths[0]=0
	pathMap := make([][]int, len+1)
	pathMap[0] = []int{0}
	for i:=0;i<=len;i++ {
		for _, val := range input[i] {
			length := val.weight
			next := val.right
			if (paths[i] + length) > paths[next] {
				paths[val.right] = paths[i] + val.weight
				pathMap[next] = addtoEnd(pathMap[i], next)
				fmt.Println(pathMap[next],pathMap[13], pathMap[18],next)
			}
		}
	}
	fmt.Println(paths)
	return pathMap
}

func addtoEnd(slice []int, num int) []int {
	res := make([]int, len(slice)+1)
	for i, val := range slice {
		res[i] = val
	}
	res[len(res)-1] = num
	return res
}


func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	//add all edges
	scanner := bufio.NewScanner(inputFile)
	i := 0
	edgemap := make(map[int][]edge)
	len := 0
	for scanner.Scan() {
		if i == 0 {
			ends := strings.Split(scanner.Text(), " ")
			len, err = strconv.Atoi(ends[1])
			edgemap[len] = make([]edge, 0)
			i++
		} else {
			addEgde(scanner.Text(), edgemap)
		}
	}
	//nodes := sortedKeys(edgemap)
	//printEdge(nodes, edgemap)
	path := longpath(edgemap, len)
	fmt.Println(path)

}

func printEdge(nodes []int,input map[int][]edge) {
	for _, val := range nodes {
		fmt.Println("node ", val)
		for _, this := range input[val] {
			fmt.Println("edge ", this.right, "weight ", this.weight)
		}
	}
}

func sortedKeys(input map[int][]edge) []int{

	keys := make([]int, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
    sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func addEgde(line string, edgemap map[int][]edge) {
	lines := strings.Split(line, " ")

	ints := make([]int, 0)
	for _, val := range lines {
		num, err := strconv.Atoi(val)
		check(err)
		ints = append(ints, num)
	}

	newedge := edge{
		right:  ints[1],
		weight: ints[2],
	}
	edgemap[ints[0]] = append(edgemap[ints[0]], newedge)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
