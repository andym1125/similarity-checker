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

/* Returns a list of substrings that are common between the control and foreign string.
 * This list will prioritize longer substrings, and the substrings in the list will not overlap
 * with each other. Thus, the return is considered pruned.
 * param: control The control string.
 * param: foreign The foreign string to compare against.
 * return: A list of common substrings
 */
func getCommonSubstrings(control []byte, foreign []byte) []Substring {

	lcsm, _ := LCSubstr(control, foreign)
	return calcCommonSubstrings(lcsm, control)
}

/* Calculates the percentage of the control string that is comprised of overlapping substrings.
 * param: control The control string.
 * param: overlap A list of Substrings that overlap the control string.
 * return: A float [0, 1] representing the percent overlap.
 */
func getPercentage(control []byte, overlap []Substring) float64 {

	s := prune(overlap)

	length := 0
	for _, a := range s {
		length += a.Len
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
	length := len(overlap[0].FullString)
	ret := []Substring{}
	visit := CreateVisitorMap(length)
	for _, s := range overlap {

		if len(s.FullString) != length {
			panic(fmt.Sprintf("prune(): Substring had a length of %d when expected length of %d",
				len(s.FullString), length))
		}

		if !visit.IsInRange(s.Start, s.End, VISITED) {
			visit.SetRange(s.Start, s.End, VISITED)
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
			Value:      fullString[maxI-max : maxI], //reduceSm[maxI-max:maxI], //TODO figure out how to elegantly get this info
			Start:      maxI - max,
			End:        maxI,
			Len:        max,
			FullString: fullString, //reduceSm,
		}

		//Confirm it is a valid substr that doesn't overlap
		checkVisit := visited[substr.Start : substr.End+1]
		if contains(VISITED, checkVisit) != -1 {

			//If invalid, mark all substring as visited
			for i := substr.Start; i < substr.End+1; i++ {
				visited[i] = int(math.Max(float64(visited[i]), float64(MARKED))) //Don't override if VISITED
			}
			continue
		}

		substrings = append(substrings, substr)

		//Visit indices used in this substr
		for i := substr.Start; i < substr.End+1; i++ {
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
	Value      []byte `json:"value"`
	Start      int    `json:"start"` //inclusive
	End        int    `json:"end"`   //exclusive
	Len        int    `json:"len"`
	FullString []byte `json:"fullString"`
}

func SortByLength(list *[]Substring) {
	l := *list

	if len(l) <= 1 {
		return
	}

	//Selection sort :(
	for cursor := 0; cursor < len(l); cursor++ {

		//Find max len
		max := l[cursor].Len
		idx := cursor
		for i := cursor; i < len(*list); i++ {

			if max < l[i].Len {
				max = l[i].Len
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

	fmt.Print("{ s: ", s.Start, " e: ", s.End, " l: ", s.Len,
		" '", string(s.Value), "' of '", string(s.FullString), "' }\n")
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
