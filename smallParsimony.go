package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*

SmallParsimony(T, Character)
    for each node v in tree T
        Tag(v) ← 0
        if v is a leaf
            Tag(v) ← 1
            for each symbol k in the alphabet
                if Character(v) = k
                    sk(v) ← 0
                else
                    sk(v) ← ∞
    while there exist ripe nodes in T
        v ← a ripe node in T
        Tag(v) ← 1
        for each symbol k in the alphabet
            sk(v) ← minimumall symbols i {si(Daughter(v))+αi,k} + minimumall symbols j {sj(Son(v))+αj,k}
    return minimum over all symbols k {sk(v)}
*/

func smallParsimony(tree []*node, ind int) {
	tag := make(map[*node]int, 0)
	for _, n := range tree {
		if n.dgt == nil {
			tag[n] = 1
			n.score = dictN()
			for nx := range n.score {
				if nx == n.seq[ind] {
					n.score[nx] = 0
				}
			}
		}
	}
	for {
		ripe := findRipe(tree, tag)
		if len(ripe) == 0 {
			break
		}
		for _, v := range ripe {
			tag[v] = 1
			gradeInternal(v)
		}
	}

}

func gradeInternal(v *node) {
	v.score = dictN()
	for k := range v.score {
		gk := gradeLeftRight(v, k)
		v.score[k] = gk
	}
}

func gradeLeftRight(v *node, k string) int {
	tmpD := make([]int, 0)
	for d := range v.dgt.score {
		if d == k {
			tmpD = append(tmpD, v.dgt.score[d])
		} else {
			tmpD = append(tmpD, v.dgt.score[d]+1)
		}
	}
	tmpS := make([]int, 0)
	for s := range v.son.score {
		if s == k {
			tmpS = append(tmpS, v.son.score[s])
		} else {
			tmpS = append(tmpS, v.son.score[s]+1)
		}
	}
	return minInt(tmpD) + minInt(tmpS)
}

func addLabel(root *node) {
	minRl := minMap(root.score)
	root.seq = append(root.seq, minRl)
	q := make([]*node, 0)
	q = append(q, root)
	for len(q) != 0 {
		v := q[0]
		q = q[1:]
		if v.dgt.dgt != nil {
			q = append(q, v.dgt, v.son)
		} else {
			continue
		}
		gradeNext(v.dgt)
		v.dgt.seq = append(v.dgt.seq, minMap(v.dgt.score))
		gradeNext(v.son)
		v.son.seq = append(v.son.seq, minMap(v.son.score))
	}
}

func gradeNext(n *node) {
	parentN := n.pnt.seq[len(n.pnt.seq)-1]
	for key := range n.score {
		if key != parentN {
			n.score[key] += 1
		}
	}
}

//an internal node of T ripe if its tag is 0 but its children’s tags are both 1
func findRipe(t []*node, tag map[*node]int) []*node {
	ripe := make([]*node, 0)
	for _, n := range t {
		if n.dgt != nil && n.son != nil {
			if tag[n] == 0 && tag[n.dgt] == 1 && tag[n.son] == 1 {
				ripe = append(ripe, n)
			}
		}
	}
	return ripe
}

func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	i, tLeav := 0, 0
	tree := make([]string, 0)
	for scanner.Scan() {
		if i == 0 {
			tLeav, err = strconv.Atoi(scanner.Text())
		} else {
			tree = append(tree, scanner.Text())
		}
		i++
	}
	root, nodes, lth := parseTree(tree, tLeav)
	for i := 0; i < lth; i++ {
		smallParsimony(nodes, i)
		addLabel(root)
	}

	f, err := os.Create("out.txt")
	check(err)
	defer f.Close()

	totalDist(nodes, root)
	fmt.Fprintln(f, parsimony)
	printAdjacency(nodes, f)

}

var parsimony int

func totalDist(tree []*node, root *node) {
	for _, n := range tree {
		if n.dgt != nil {
			i1 := hammingDist(n.seq, n.dgt.seq)
			i2 := hammingDist(n.seq, n.son.seq)
			parsimony = parsimony + i1 + i2
		}
	}
	rootd := hammingDist(root.seq, root.dgt.seq)
	roots := hammingDist(root.seq, root.son.seq)
	ds := hammingDist(root.son.seq, root.dgt.seq)
	parsimony = parsimony - rootd - roots + ds
}

func printAdjacency(tree []*node, f *os.File) {
	for _, n := range tree {
		if n.pnt == nil {
			fmt.Fprint(f, strings.Join(n.son.seq, ""), "->", strings.Join(n.dgt.seq, ""), ":", hammingDist(n.son.seq, n.dgt.seq))
			fmt.Fprintln(f)
			fmt.Fprint(f, strings.Join(n.dgt.seq, ""), "->", strings.Join(n.son.seq, ""), ":", hammingDist(n.son.seq, n.dgt.seq))
			fmt.Fprintln(f)
			continue
		}
		if n.dgt != nil {
			fmt.Fprint(f, strings.Join(n.seq, ""), "->", strings.Join(n.dgt.seq, ""), ":", hammingDist(n.seq, n.dgt.seq))
			fmt.Fprintln(f)
			fmt.Fprint(f, strings.Join(n.dgt.seq, ""), "->", strings.Join(n.seq, ""), ":", hammingDist(n.seq, n.dgt.seq))
			fmt.Fprintln(f)
			fmt.Fprint(f, strings.Join(n.seq, ""), "->", strings.Join(n.son.seq, ""), ":", hammingDist(n.seq, n.son.seq))
			fmt.Fprintln(f)
			fmt.Fprint(f, strings.Join(n.son.seq, ""), "->", strings.Join(n.seq, ""), ":", hammingDist(n.seq, n.son.seq))
			fmt.Fprintln(f)
		}
	}
}

func hammingDist(str1, str2 []string) int {
	c := 0
	for i := range str1 {
		if str1[i] != str2[i] {
			c++
		}
	}
	return c
}

//gives an alphabet map
func dictN() map[string]int {
	n := make(map[string]int, 0)
	n["A"] = 100000
	n["C"] = 100000
	n["G"] = 100000
	n["T"] = 100000
	return n
}

type node struct {
	pnt   *node
	dgt   *node
	son   *node
	seq   []string
	score map[string]int
}

func (this *node) addChild(n *node) {
	if this.dgt == nil {
		this.dgt = n
	} else if this.son == nil {
		this.son = n
	} else {
		panic("incorrect add child")
	}
	n.pnt = this
}

func parseTree(tree []string, tLeav int) (*node, []*node, int) {
	nodeMap := make(map[int]*node, 0)
	treeMap := make(map[*node]int, 0)

	lth := 0
	for i, line := range tree {
		if i%2 == 0 {
			continue
		}
		sline := strings.Split(line, "->")
		a, b := sline[0], sline[1]

		left, err := strconv.Atoi(a)
		if i < 2*tLeav {
			cN := node{seq: strings.Split(b, "")}
			lth = len(cN.seq)
			if nodeMap[left] == nil {
				pN := node{}
				nodeMap[left] = &pN
				treeMap[&pN]++
			}
			nodeMap[left].addChild(&cN)
			treeMap[&cN]++
			check(err)
		} else {
			right, err := strconv.Atoi(b)
			check(err)
			if nodeMap[left] == nil {
				lN := node{}
				nodeMap[left] = &lN
				treeMap[&lN]++
			}
			if nodeMap[left].dgt == nil || nodeMap[left].son == nil {
				nodeMap[left].addChild(nodeMap[right])
			}
		}
	}

	root := node{}
	for k := range treeMap {
		if k.pnt == nil {
			root.addChild(k)
		}
	}
	allnodes := make([]*node, 0)
	for key := range treeMap {
		allnodes = append(allnodes, key)
	}
	allnodes = append(allnodes, &root)
	return &root, allnodes, lth
}

func minInt(l []int) int {
	r := l[0]
	for _, val := range l {
		if val < r {
			r = val
		}
	}
	return r
}

func minMap(m map[string]int) string {
	r := m["A"]
	rk := "A"
	for key, val := range m {
		if val < r {
			r = val
			rk = key
		}
	}
	return rk
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
