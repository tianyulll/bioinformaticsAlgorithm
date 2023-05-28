package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
BWMatching(LastColumn, Pattern, LastToFirst)
    top ← 0
    bottom ← |LastColumn| − 1
    while top ≤ bottom
        if Pattern is nonempty
            symbol ← last letter in Pattern
            remove last letter from Pattern
            if positions from top to bottom in LastColumn contain an occurrence of symbol
                topIndex ← first position of symbol among positions from top to bottom in LastColumn
                bottomIndex ← last position of symbol among positions from top to bottom in LastColumn
                top ← LastToFirst(topIndex)
                bottom ← LastToFirst(bottomIndex)
            else
                return 0
        else
            return bottom − top + 1
*/

func BWmatching(input, pattern string) {
	input = input+"$"
	patterns := strings.Split(pattern, " ")
	suff := suffixArray(input)

	_, lc := burrowMatrix(input)
	lcSplit, fc := labeled(lc) //list of fc and lc columns
	ltfMap := lastToFirstMap(fc, lcSplit)

	f, err := os.Create("out.txt")
	check(err)
	defer f.Close()
	
	for _, pat := range patterns {
		b, t := match(strings.Split(pat, ""), ltfMap, lcSplit)
		fmt.Fprint(f, pat, ":")
		if b == 0 && t == 0 {
			fmt.Fprintln(f)
			continue
		}
		for i := t; i <= b; i++ {
			fmt.Fprint(f, suff[i], " ")
		}
		fmt.Fprintln(f)
	}
}


func match(pattern []string, ltfMap map[string]int, lc []string) [][]int {
	top := 0
	bottom := len(lc)-1

	res := make([][]int, 0)
	for top <= bottom {
		if len(pattern) != 0 {
			symbol := pattern[len(pattern)-1]
			pattern = pattern[:len(pattern)-1]
			
			first := true
			dne := false
			tmpTop, tmpBot := 0,0
			for i := top; i <= bottom; i++  {
				if  strings.Split(lc[i], "")[0] == symbol {
					dne = true
					if first {
						tmpTop = ltfMap[lc[i]] 
						first = false
					} 
					tmpBot = ltfMap[lc[i]]
				}
			}
			if ! dne {res = append(res, []int{0,0})}
			top = tmpTop
			bottom = tmpBot
		} else {
			res = append(res, []int{top, bottom})
		}
	}

	return res
} 

func lastToFirstMap(fc, lc []string) map[string]int {
	fm := make(map[string]int, 0)
	lm := make(map[string]int, 0)
	for i, val := range fc {
		fm[val] = i
	}
	for i, val := range lc {
		lm[val] = i
	}
	return fm
	
}


func main() {
	inputFileName := os.Args[1]
	inputFile, err := os.Open(inputFileName)
	check(err)
	defer inputFile.Close()

	var inputSeq string
	var pattern string
	i := 0

	scanner := bufio.NewScanner(inputFile)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		if i == 0 {inputSeq = scanner.Text();i++} else {pattern = scanner.Text()}
	}

	//fmt.Println(burrowMatrix(inputSeq)) 
	//fmt.Println(inverseBurrow(inputSeq))
	fmt.Println(len(pattern))
	BWmatching(inputSeq, pattern)
	//multiMatch(inputSeq, pattern)

}

func multiMatch(input string, pattern string) {
	_, lc := burrowMatrix(input)
	fmt.Println("last col", lc)
	//suffixarray := suffixArray(input)
	firstO := firstOccurence(strings.Split(lc, ""))
	patterns := strings.Split(pattern, " ")

	for _, pat := range patterns {
		t, b := bettermatch(firstO, strings.Split(lc, ""), strings.Split(pat, ""))
		fmt.Println(t,b)
	}
}


/*
func multiMatch(input string, pattern string) {
	_, lc := burrowMatrix(input)
	fmt.Println("last col", lc)
	//suffixarray := suffixArray(input)
	firstO := firstOccurence(strings.Split(lc, ""))
	patterns := strings.Split(pattern, " ")

	for _, pat := range patterns {
		t, b := bettermatch(firstO, strings.Split(lc, ""), strings.Split(pat, ""))
		fmt.Println(t,b)
	}
}

*/

func bettermatch(firsto map[string]int, lc, pattern []string) (int, int) {
	top := 0
	bottom := len(lc)-1

	for top <= bottom {
		if len(pattern) != 0 {
			symbol := pattern[len(pattern)-1]
			pattern = pattern[:len(pattern)-1]
			dne := true
			tmpTop, tmpBot := 0,0

			for i := top; i <= bottom; i++  {
				if  strings.Split(lc[i], "")[0] == symbol {
					dne = false
				}
			}

			if !dne {
				countT := countLastCol(lc, symbol, top)
				countB := countLastCol(lc, symbol, bottom+1)
				tmpTop = countT + firsto[symbol]
				tmpBot = countB + firsto[symbol]-1
			} else {
				return 0, 0
			}
			top = tmpTop
			bottom = tmpBot
		} else {
			return bottom, top
		}
	}
	return 0, 0
} 

//produce the list of suffix array
func suffixArray(input string) []int {
	suffix := make([]string, 0)
	pos := make(map[string]int, 0)
	sinput := strings.Split(input, "")
	for i:=0; i <len(sinput); i++ {
		aSuffix := strings.Join(sinput[i:], "")
		suffix = append(suffix, aSuffix)
		pos[aSuffix] = i
	}
	sort.Strings(suffix)

	sfArray := make([]int, 0)
	for _, val := range suffix {
		sfArray = append(sfArray, pos[val])
	}
	return sfArray
}

//first occurence list of the last column
func firstOccurence(lc []string) map[string]int {
	tmpMap := make(map[string]int, 0)
	for _, val := range lc {
		tmpMap[val] = -1
	}
	for i, val := range lc {
		if tmpMap[val] == -1 {
			tmpMap[val] = i
		}
	}
	fmt.Println("first occurence", tmpMap)
	return tmpMap
}

//count word n for l-th rows in last column
func countLastCol(lc []string, n string, l int) int {
	count := 0
	for i := 0; i <= l; i++ {
		if lc[i] == n {
			count++
		}
	}
	return count
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}


//easy question
func burrowMatrix(s string) (string, string) {
	mat := make([][]string, 0)
	ss := strings.Split(s, "")
	mat = append(mat, ss)
	for i := len(s)-2; i >=0; i-- {
		left := ss[:i+1]
		right := ss[i+1:]
		tmp := append(right, left...)
		mat = append(mat, tmp)
	}
	toSort := make([]string, 0)
	for i := range mat {
		toSort = append(toSort, strings.Join(mat[i], ""))
	}
	sort.Strings(toSort)
	lastC := make([]string, 0)
	ln := len(ss)-1
	for _, val := range toSort {
		line := strings.Split(val, "")
		lastC = append(lastC, line[ln])
	}
	lc := strings.Join(lastC, "")

	fc := make([]string, 0)
	for _, val := range toSort {
		fc = append(fc, strings.Split(val, "")[0])
	}
	return strings.Join(fc, ""), lc
} 

func inverseBurrow(lastc string) string {
	lastc = lastc
	labelLast, fcol := labeled(lastc)

	searchMap := make(map[string]string, 0)
	for i, key := range labelLast {
		searchMap[key] = fcol[i]
	}

	out := make([]string, 0)
	ini := fcol[0]
	//fmt.Println(fcol)
	//fmt.Println(searchMap)
	for len(out) < len(lastc) {
		tmp := strings.Split(searchMap[ini], "")
		out = append(out, tmp[0])
		ini = searchMap[ini]
	}
	return strings.Join(out, "")
}


//output labeled last, first column
func labeled(s string) ([]string, []string) {
	ss := strings.Split(s, "")
	letterMap := make(map[string]int, 0)
	outSeq := make([]string, 0)
	for _, v := range ss {
		letterMap[v]++
		tmp := strconv.Itoa(letterMap[v])
		outSeq = append(outSeq, v+tmp)
	}
	fcol := sortedlist(letterMap)
	return outSeq, fcol
}

func sortedlist(n map[string]int) []string {
	keys := make([]string, 0)
	for key := range n {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]string, 0)
	for _, key := range keys {
		for i := 1; i <= n[key];i++ {
			tmp := strconv.Itoa(i)
			out = append(out, key+tmp)
		}
	}
	return out
}
