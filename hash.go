package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

var uniset []byte
var hash_size int = 64
var chunk = 55 //md5 will pad 1 byte

func init() {

	var i uint8
	for i = 0; i < 128; i++ {
		uniset = append(uniset, i)
	}
}

func main() {

	_, err := ioutil.ReadFile("maury.txt")

	if err != nil {
		fmt.Println("ERROR READING")
	}

	smallD1 := []uint8{100, 20, 73, 20, 43, 120, 45, 20, 36}
	smallD2 := []uint8{100, 20, 73, 20, 120, 45, 20, 36}

	fmt.Println("original hash s1", hash(smallD1))
	fmt.Println("original hash s0", hash(smallD2))
	fmt.Println("New hash s1", fastforward(hash(smallD2), 1))
}

/* Fastforwards the hash so that 2 lsh's can be compared
 * param to indicates number of passes to complete in total
 * the first bytes of data indicates how many passes have already been completed
 */
func fastforward(data []uint8, to uint8) []uint8 {

	newData := data[1:]

	var i uint8 = 0
	for i = 0; i < to-data[0]; i++ {

		newData = append(newData, 0)
		newData = slim(newData)
	}

	return append([]uint8{to}, newData...)
}

/* Performs a locality sensitive hash
 * First byte indicates how many passes
 * Next 64 indicate the hash itself */
func hash(data []uint8) []uint8 {

	/* Padding */
	for i := len(data); i < hash_size; i++ {
		data = append(data, 0)
	}

	/* Initial hash of data */
	var newData []uint8
	for i := 0; i < len(data); i++ {

		t := []uint8{data[i]}
		sum := md5.Sum(t)
		newData = append(newData, sum[0])
	}

	/* Slim down hash to desired size */
	var passes uint8 = 0
	for len(newData) != hash_size {
		newData = slim(newData)
		passes++
	}

	/* Prepend # passes */
	newData = append([]uint8{passes}, newData...)

	return newData
}

/* Slims the given intermediate hash from n bytes to n-1 bytes */
func slim(data []uint8) []uint8 {

	var newData []uint8
	for i := 0; i < len(data)-1; i++ {

		d := []uint8{data[i], data[i+1]}
		sum := md5.Sum(d)
		newData = append(newData, sum[0]) //Append last byte of checksum as new value
	}

	return newData
}
