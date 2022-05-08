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

func main() {

	s1 := []byte("I've always wanted to ride my bicycle, and can't wait to ride it up on the highway.")
	s2 := []byte("On a dark desert highway. I wanted to ride my bicycle.")

	fmt.Println(len(s1), len(s2))

	printSubstrArr(getCommonSubstrings(s1, s2))

	// lcsm, i := LCSubstr(s1, s2)
	// printLCSMatrix(lcsm, i)

	// ssa := calcCommonSubstrings(lcsm, s1)
	// printSubstrArr(ssa)
}

func getCommonSubstrings(control []byte, foreign []byte) []Substring {

	lcsm, _ := LCSubstr(control, foreign)
	return calcCommonSubstrings(lcsm, control)
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

/* Prints a prettified version of the Substring.
 * The byte strings are printed as character strings.
 */
func (s Substring) PrettyPrint() {

	fmt.Print("{ s: ", s.start, " e: ", s.end, " l: ", s.len,
		" '", string(s.value), "' of '", string(s.fullString), "' }\n")
}
