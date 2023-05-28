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
NeighborJoining(D)
    n ← number of rows in D
    if n = 2
        T ← tree consisting of a single edge of length D1,2
        return T
    D* ← neighbor-joining matrix constructed from the distance matrix D
    find elements i and j such that D*i,j is a minimum non-diagonal element of D*
    Δ ← (TotalDistanceD(i) - TotalDistanceD(j)) /(n - 2)
    limbLengthi ← (1/2)(Di,j + Δ)
    limbLengthj ← (1/2)(Di,j - Δ)
    add a new row/column m to D so that Dk,m = Dm,k = (1/2)(Dk,i + Dk,j - Di,j) for any k
    D ← D with rows i and j removed
    D ← D with columns i and j removed
    T ← NeighborJoining(D)
    add two new limbs (connecting node m with leaves i and j) to the tree T
    assign length limbLengthi to Limb(i)
    assign length limbLengthj to Limb(j)
    return T
*/

type clust struct {
	node int
	leaf int
}

//engine for NJ
func nj(d [][]float64, dim int) {
	cluster := make([]clust, dim) //list of clusters formed during merging
	for i := range cluster {
		cluster[i].node = i
		cluster[i].leaf = 1
	}
	graph := make(map[int][]edge, 0)

	i := dim //new nodes
	for len(d) > 2 {

		//calculation
		dStar := njMatrix(d, len(d))
		ci, cj := closeClust(dStar)
		delta := (totalDist(d, ci) - totalDist(d, cj)) / float64(len(d)-2)
		limbI := (d[ci][cj] + delta) / 2
		limbJ := (d[ci][cj] - delta) / 2
		
		//add edge, bidirectional
		graph[i] = append(graph[i], setEdge(cluster[ci].node, limbI))
		graph[i] = append(graph[i], setEdge(cluster[cj].node, limbJ))
		graph[cluster[ci].node] = append(graph[cluster[ci].node], setEdge(i, limbI))
		graph[cluster[cj].node] = append(graph[cluster[cj].node], setEdge(i, limbJ))
		//modify the matrix
		d = addRowCol(d, ci, cj, cluster)

		//merge into new cluster
		tmpclus := clust{node: i, leaf: cluster[ci].leaf + cluster[cj].leaf}
		if ci > cj {
			cluster = append(cluster[:ci], cluster[ci+1:]...)
			cluster = append(cluster[:cj], cluster[cj+1:]...)
		} else {
			cluster = append(cluster[:cj], cluster[cj+1:]...)
			cluster = append(cluster[:ci], cluster[ci+1:]...)
		}
		cluster = append(cluster, tmpclus)
		i++ //increment on new node
	}
	//handle the base case
	graph[cluster[0].node] = append(graph[cluster[0].node], setEdge(cluster[1].node, d[0][1]))
	graph[cluster[1].node] = append(graph[cluster[1].node], setEdge(cluster[0].node, d[0][1]))
	//output result
	parseAdjacency(graph)
}

func parseAdjacency(graph map[int][]edge) {
	adj := make([][]edge, len(graph))
	for i := range adj {
		for _, e := range graph[i] {
			adj[i] = append(adj[i], e)
		}
	}

	f, err := os.Create("out.txt")
	check(err)
	defer f.Close()
	for i := range adj {
		for _, e := range adj[i] {
			fmt.Fprint(f, i, "->", e.end, ":", e.lth)
			fmt.Fprintln(f)
		}
	}
}

//calculate new distance matrix
//Dk,m = Dm,k = (1/2)(Dk,i + Dk,j - Di,j)
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
		tmpVal := (dc[row][i] + dc[row][j] - dc[i][j]) / 2
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

type edge struct {
	end int
	lth float64
}

func setEdge(end int, lth float64) edge {
	newEdge := edge{
		end: end,
		lth: lth,
	}
	return newEdge
}

//return the neighbor joining matrix
//D*i,j = (n - 2) · Di,j - TotalDistanceD(i) - TotalDistanceD(j).
func njMatrix(mat [][]float64, dim int) [][]float64 {
	newM := make([][]float64, dim)
	for i := range newM {
		newM[i] = make([]float64, dim)
	}
	for i := range mat {
		for j := range mat[i] {
			if i == j {
				continue
			}
			newM[i][j] = float64(dim-2)*mat[i][j] - totalDist(mat, i) - totalDist(mat, j)
		}
	}
	return newM
}

//total distance given X in matrix
func totalDist(mat [][]float64, x int) float64 {
	return sum(mat[x])
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
	nj(matrixInt, dim)

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

func sum(list []float64) float64 {
	tot := 0.0
	for _, val := range list {
		tot += val
	}
	return tot
}
