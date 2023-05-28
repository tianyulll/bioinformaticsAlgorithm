package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*
compute LimbLength(j) by finding the minimum value of (D_{i,j} + D_{j,k} - D_{i,k})/2
over all pairs of leaves i and k.
*/
func limbLen(matrix [][]int, j int, dim int) int {
	candi := make([]int, 0)
	for i := 0; i < dim; i++ {
		for k := 0; k < dim; k++ {
			if i == k || i == j || k == j {
				continue
			}
			candi = append(candi, (matrix[i][j]+matrix[j][k]-matrix[i][k])/2)
		}
	}
	return listMin(candi)
}

/*

UPGMA(D, n)
    Clusters ← n single-element clusters labeled 1, ... , n
    construct a graph T with n isolated nodes labeled by single elements 1, ... , n
    for every node v in T
        Age(v) ← 0
    while there is more than one cluster
        find the two closest clusters Ci and Cj

		merge Ci and Cj into a new cluster Cnew with |Ci| + |Cj| elements
        add a new node labeled by cluster Cnew to T
        connect node Cnew to Ci and Cj by directed edges
        Age(Cnew) ← DCi, Cj / 2

		remove the rows and columns of D corresponding to Ci and Cj
        remove Ci and Cj from Clusters
        add a row/column to D for Cnew by computing D(Cnew, C) for each C in Clusters
        add Cnew to Clusters
    root ← the node in T corresponding to the remaining cluster
    for each edge (v, w) in T
        length of (v, w) ← Age(v) - Age(w)
    return T
*/

type clust struct {
	node int
	leaf int
}

//engine for UPGMA
func upgma(dm [][]float64, dim int) {
	//initialize cluster, graph, age
	cluster := make([]clust, dim) //list of map, key = node, val = leaf
	for i := range cluster {
		cluster[i].node = i
		cluster[i].leaf = 1
	}
	graph := make(map[int][]int, 0)
	age := make([]float64, dim)

	i := dim //stores c_new
	//iterate and merge the distance matrix
	for len(dm) >= 2 {
		ci, cj := closeClust(dm)
		//merge into new cluster
		tmpclus := clust{node: i, leaf: cluster[ci].leaf + cluster[cj].leaf}
		//add edge, bidirectional
		graph[i] = append(graph[i], cluster[ci].node)
		graph[i] = append(graph[i], cluster[cj].node)
		graph[cluster[ci].node] = append(graph[cluster[ci].node], i)
		graph[cluster[cj].node] = append(graph[cluster[cj].node], i)
		//add age
		age = append(age, dm[ci][cj]/2)
		//modify the matrix
		dm = addRowCol(dm, ci, cj, cluster)
		//modify the cluster
		if ci > cj {
			cluster = append(cluster[:ci], cluster[ci+1:]...)
			cluster = append(cluster[:cj], cluster[cj+1:]...)
		} else {
			cluster = append(cluster[:cj], cluster[cj+1:]...)
			cluster = append(cluster[:ci], cluster[ci+1:]...)
		}
		cluster = append(cluster, tmpclus)
		fmt.Println(cluster)
		i++ //increment on new node
	}
	//fmt.Println("age", age)
	//fmt.Println("graph", graph)
	parseAdjacency(graph, age) //print out the result
}

//decrease the size of the distance matrix and 
//calculate new distance matrix based on (DCi,Cm ·|Ci|+DCj,Cm ·|Cj|) / (|Ci|+|Cj|)
func addRowCol(d [][]float64, i, j int, cluster []clust) [][]float64 {
	var dnew [][]float64
	dc := copyM(d) //have to deep copy heres
	if i > j {
		dtmp := deleteRowCol(d, i)
		dnew = deleteRowCol(dtmp, j)
	} else {
		dtmp := deleteRowCol(d, j)
		dnew = deleteRowCol(dtmp, i)
	}
	tmpCol := make([]float64, 0)
	//append row
	for row := range dc {
		if i == row || j == row {
			continue
		}
		tmpVal := (dc[row][i]*float64(cluster[i].leaf) + dc[row][j]*float64(cluster[j].leaf)) / float64(cluster[i].leaf+cluster[j].leaf)
		tmpCol = append(tmpCol, tmpVal)
	}
	//append column
	for row := range dnew {
		dnew[row] = append(dnew[row], tmpCol[row])
	}
	tmpCol = append(tmpCol, 0)
	dnew = append(dnew, tmpCol)
	return dnew
}

//matrix utilities
//returns the closest 2 position in the distance matrix
func closeClust(dm [][]float64) (int, int) {
	minI, minJ := 0, 0
	minVal := math.MaxFloat32
	for i := range dm {
		for j := range dm {
			if i == j {
				continue
			}
			if minVal > dm[i][j] {
				minVal = dm[i][j]
				minI = i
				minJ = j
			}
		}
	}
	return minI, minJ
}

//delete c row and c column from the distance matrix
func deleteRowCol(dm [][]float64, c int) [][]float64 {
	res := make([][]float64, 0)
	for i := range dm {
		if i == c {
			continue
		}
		tmp := append(dm[i][:c], dm[i][c+1:]...)
		res = append(res, tmp)
	}
	return res
}

//form adjacency list based on the graph
//print out the adjacency list
func parseAdjacency(graph map[int][]int, age []float64) {
	//create edge list, potentially useful
	adj := make([][]edge, len(graph))
	for i := range adj {
		for _, val := range graph[i] {
			newEdge := edge{stt: i, end: val, lth: abs(age[i] - age[val])}
			newEdge.lth = math.Round(newEdge.lth*1000) / 1000
			adj[i] = append(adj[i], newEdge)
		}
	}

	f, err := os.Create("out.txt")
	check(err)
	defer f.Close()
	//why bother formatting...
	for _, nodes := range adj {
		for _, edges := range nodes {
			fmt.Fprint(f, edges.stt, "->", edges.end, ":", edges.lth)
			fmt.Fprintln(f)
		}
	}
}

func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i, dim := 0, 0
	matrix := make([]string, 0)
	for scanner.Scan() {
		if i == 0 {
			dim, err = strconv.Atoi(scanner.Text())
		} else {
			matrix = append(matrix, scanner.Text())
		}
		i++
	}

	matrixInt := parseMatrix(matrix)
	upgma(matrixInt, dim)

}

type edge struct {
	stt int
	end int
	lth float64
}

func parseMatrix(input []string) [][]float64 {
	matrix := make([][]float64, len(input))
	for i := range input {
		line := strings.Fields(input[i])
		for _, word := range line {
			val, err := strconv.Atoi(word)
			matrix[i] = append(matrix[i], float64(val))
			check(err)
		}
	}
	return matrix
}

//helpfer Functions
func listMin(list []int) int {
	min := math.MaxInt
	for _, val := range list {
		if val < min {
			min = val
		}
	}
	return min
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func abs(a float64) float64 {
	if a > 0 {
		return a
	} else {
		return -a
	}
}

func copyM(in [][]float64) [][]float64 {
	out := make([][]float64, len(in))
	for i, row := range in {
		for _, val := range row {
			out[i] = append(out[i], val)
		}
	}
	return out
}
