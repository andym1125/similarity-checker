package main

import (
	"fmt"
	"math"
)

const COMPARE_THRESH = 3 //Exclusive

/* VISITED: 	Used in a substring
				 CANNOT be used in another substring
 * MARKED:		Perhaps under threshold, cannot initiate substring;
 				 CAN be used in another substring
 * UNVISITED:	Can initiate substring;
 				 CAN be used in another substring */
const UNVISITED = 0
const MARKED = 1
const VISITED = 2

//func main() {

/* Test Percentage/Prune methods
s1 := []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")
s2 := []byte("On a dark desert highway. I wanted to ride my bicycle.")
ssa := getCommonSubstrings(s1, s2)

fmt.Println(getPercentage(s1, ssa))/**/

/* Test Substring SortByLength
s := []Substring{
	{start: 74, end: 83, len: 9, value: []byte(" highway."), fullString: []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")},
	{start: 7, end: 10, len: 3, value: []byte("way"), fullString: []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")},
	{start: 11, end: 37, len: 26, value: []byte(" wanted to ride my bicycle"), fullString: []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")},
	{start: 48, end: 51, len: 3, value: []byte(" wa"), fullString: []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")},
	{start: 53, end: 62, len: 9, value: []byte(" to ride "), fullString: []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")},
}
SortByLength(&s)
printSubstrArr(s)
/**/

/* Test LCS procedures
s1 := []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")
s2 := []byte("On a dark desert highway. I wanted to ride my bicycle.")

fmt.Println(len(s1), len(s2))

printSubstrArr(getCommonSubstrings(s1, s2))

//Intermediate checking
// lcsm, i := LCSubstr(s1, s2)
// printLCSMatrix(lcsm, i)

// ssa := calcCommonSubstrings(lcsm, s1)
// printSubstrArr(ssa)
/**/
//}

func getCommonSubstrings(control []byte, foreign []byte) []Substring {

	lcsm, _ := LCSubstr(control, foreign)
	return calcCommonSubstrings(lcsm, control)
}

func getPercentage(control []byte, overlap []Substring) float64 {

	s := prune(overlap)

	length := 0
	for _, a := range s {
		length += a.len
	}

	debug := false
	if debug {

		fmt.Println("DEBUG getPercentage():")

		for _, substr := range s {
			substr.PrettyPrint()
		}
		fmt.Println("Similar len:", length, "Control len:", len(control))
		fmt.Println("RETURN:", float64(length)/float64(len(control)))

		fmt.Println()
	}

	return float64(length) / float64(len(control))
}

/* Ensures that none of the Substrings overlap each other.
 * Does not ensure the Substrings actually belong to the control string,
 * only that all fullString attributes have the same length.
 */
func prune(overlap []Substring) []Substring {

	if len(overlap) == 0 {
		return []Substring{} //Return empty array if passed empty array
	}

	//Want to visit larger substrings before smaller
	SortByLength(&overlap)

	//Visit each, skipping if overlapping another substring
	length := len(overlap[0].fullString)
	ret := []Substring{}
	visit := CreateVisitorMap(length)
	for _, s := range overlap {

		if len(s.fullString) != length {
			panic(fmt.Sprintf("prune(): Substring had a length of %d when expected length of %d",
				len(s.fullString), length))
		}

		if !visit.IsInRange(s.start, s.end, VISITED) {
			visit.SetRange(s.start, s.end, VISITED)
			ret = append(ret, s)
		}
	}

	return ret
}

/* Prints the intermediate Least Common Suffix matrix.
 * param: a The LCS matrix
 * param: i The size of the largest LCS
 */
func printLCSMatrix(a [][]int, i int) {

	fmt.Println("PrintLCSMatrix: size", i)
	for _, x := range a {
		for _, y := range x {

			if y == 0 {
				fmt.Print("- ")
			} else {
				fmt.Print(y, " ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

/* Prints the intermediate Substring array using PrettyPrint()
 * param: a The array to print
 */
func printSubstrArr(a []Substring) {

	fmt.Println("PrintSubstrArr:")
	for _, x := range a {
		x.PrettyPrint()
	}
	fmt.Println()
}

/* Calculates the common substrings given a LCS matrix.
 * This function prefers larger substrings over smaller substrings.
 * Substrings smaller than length COMPARE_THRESH are disregarded.
 * param: lcsm The LCS matrix to use.
 * param: fullString The bytestring of the control string used in the creation of the LCS matrix.
 * return: An array of substrings calculated from the LCS matrix.
 */
func calcCommonSubstrings(lcsm [][]int, fullString []byte) []Substring {

	//Reduce matrix by turning the array of arrays into array of maxes
	var reduceSm []int
	for _, x := range lcsm {

		max, _ := max(x...)
		reduceSm = append(reduceSm, max)
	}

	//Preprocess so that vals under threshold are marked but not visited
	visited := make([]int, len(reduceSm))
	for i, _ := range visited {

		if reduceSm[i] < COMPARE_THRESH {
			visited[i] = MARKED
		}
	}

	//Until there are no more empty spaces, find comparisions (above certain threshold)
	var substrings []Substring
	for contains(UNVISITED, visited) != -1 {

		//Find curent max val
		max := 0
		maxI := -1
		for i, a := range reduceSm {
			if a > max && visited[i] == UNVISITED {
				max = a
				maxI = i
			}
		}

		//Create a substr using that max value
		substr := Substring{
			value:      fullString[maxI-max : maxI], //reduceSm[maxI-max:maxI], //TODO figure out how to elegantly get this info
			start:      maxI - max,
			end:        maxI,
			len:        max,
			fullString: fullString, //reduceSm,
		}

		//Confirm it is a valid substr that doesn't overlap
		checkVisit := visited[substr.start : substr.end+1]
		if contains(VISITED, checkVisit) != -1 {

			//If invalid, mark all substring as visited
			for i := substr.start; i < substr.end+1; i++ {
				visited[i] = int(math.Max(float64(visited[i]), float64(MARKED))) //Don't override if VISITED
			}
			continue
		}

		substrings = append(substrings, substr)

		//Visit indices used in this substr
		for i := substr.start; i < substr.end+1; i++ {
			visited[i] = VISITED
		}
	}

	return substrings
}

/* Determines the size of the largest common substring.
 * Note: if using for comparison, use x as the primary string, and y as the foreign string
 * param: x The control string
 * param: y The string to compare against
 * return: the intermediate Least Common Suffix matrix
 * return: the size of the largest common substring
 */
func LCSubstr(x []byte, y []byte) ([][]int, int) {
	m := len(x)
	n := len(y)

	lcsm := make([][]int, 0)
	for i := 0; i < m+1; i++ {
		lcsm = append(lcsm, make([]int, n+1))
	}
	result := 0 //Store len of longest common substr

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {

			//1st row/col have no meaning
			if i == 0 || j == 0 {
				lcsm[i][j] = 0

			} else if x[i-1] == y[j-1] {

				lcsm[i][j] = lcsm[i-1][j-1] + 1
				result = int(math.Max(float64(result), float64(lcsm[i][j])))

			} else {

				lcsm[i][j] = 0
			}
		}
	}

	return lcsm, result
}

/* ********** The Substring type ********** */

type Substring struct {
	value      []byte
	start      int //inclusive
	end        int //exclusive
	len        int
	fullString []byte
}

func SortByLength(list *[]Substring) {
	l := *list

	if len(l) <= 1 {
		return
	}

	//Selection sort :(
	for cursor := 0; cursor < len(l); cursor++ {

		//Find max len
		max := l[cursor].len
		idx := cursor
		for i := cursor; i < len(*list); i++ {

			if max < l[i].len {
				max = l[i].len
				idx = i
			}
		}

		//Swap w cursor
		temp := l[cursor]
		l[cursor] = l[idx]
		l[idx] = temp
	}
}

/* Prints a prettified version of the Substring.
 * The byte strings are printed as character strings.
 */
func (s Substring) PrettyPrint() {

	fmt.Print("{ s: ", s.start, " e: ", s.end, " l: ", s.len,
		" '", string(s.value), "' of '", string(s.fullString), "' }\n")
}

/* ********** VisitorMap type ********** */

//TODO: update methods before to use Visitor Map
type VisitorMap struct {
	data []int
}

func CreateVisitorMap(size int) *VisitorMap {
	return &VisitorMap{data: make([]int, size)}
}

func (m *VisitorMap) Is(idx int, value int) bool {
	return m.data[idx] == value
}

func (m *VisitorMap) IsInRange(start int, end int, value int) bool {
	ret := false
	for i := start; i < end; i++ {

		if m.data[i] == value {
			ret = true
		}
	}
	return ret
}

func (m *VisitorMap) Set(idx int, value int) int {
	ret := m.data[idx]
	m.data[idx] = value
	return ret
}

func (m *VisitorMap) SetRange(start int, end int, value int) []int {
	ret := m.data[start:end]
	for i := start; i < end; i++ {
		m.data[i] = value
	}
	return ret
}

func (m *VisitorMap) Print() {

	fmt.Println("VisitorMap: ", m.data)
}
