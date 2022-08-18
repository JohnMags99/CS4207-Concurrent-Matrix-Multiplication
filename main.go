package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

var matrixA, matrixB [][]int

// PopulateMatrix == populates matrix with random ints
func PopulateMatrix(size int) [][]int {
	rand.Seed(time.Now().UnixNano())
	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrix[i] = append(matrix[i], rand.Intn(10))
		}
	}

	return matrix
}

// PrintMatrix == prints values stored iin matrix
func PrintMatrix(matrix [][]int, id string) {
	fmt.Println(id)

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {

			fmt.Printf("%d ", matrix[i][j])
		}
		fmt.Println("\n")
		time.Sleep(time.Millisecond * 50)
	}
}

// CreateEmpty == Fill Matrix with 0
func CreateEmpty(size int) [][]int {
	matrix := make([][]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrix[i] = append(matrix[i], 0)
		}
	}
	return matrix
}

func FoxAlgorithm(ch chan [][]int) {
	var q int //q == sq root of total matrix size/processes

	//bCastA == function which broadcasts aValues each stage to be multiplied with bValues
	bCastA := func(stage int, ch chan [][]int) {
		bCast := CreateEmpty(len(matrixA))

		for i := 0; i < len(matrixA); i++ {
			kBar := (i + stage) % q
			for j := 0; j < len(matrixA[i]); j++ {
				bCast[i][j] = matrixA[i][kBar]
			}
		}

		ch <- bCast
	}

	//bShift == function to shift up rows in matrix b after each stage
	bShift := func(ch chan [][]int) {
		ch <- matrixB
		shiftB := CreateEmpty(len(matrixB))
		keys := make([]int, len(matrixB))

		for i := 0; i < len(matrixB); i++ {
			for j := 0; j < len(matrixB[i]); j++ {
				shiftB[i][j] = matrixB[i][j]
			}
		}

		for i := len(matrixB) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(keys); i++ {
			matrixB[i] = shiftB[keys[i]]
		}
	}

	//c == Matrix Multiplication function which returns matrix c
	c := func(ch chan [][]int) {
		matrixC := CreateEmpty(len(matrixA))

		//channels to store bValues & broadcasted aValues
		bs := make(chan [][]int)
		pr := make(chan [][]int)

		q = int(math.Sqrt(float64(len(matrixA) * len(matrixA[0]))))

		for stage := 0; stage < q; stage++ {
			go bCastA(stage, pr)
			go bShift(bs)

			time.Sleep(time.Millisecond * 500)
			go func() {
				a := <-pr
				b := <-bs
				for i := 0; i < len(matrixA); i++ {
					for j := 0; j < len(matrixA[i]); j++ {
						matrixC[i][j] = matrixC[i][j] + (a[i][j] * b[i][j])
					}
				}
			}()
		}
		ch <- matrixC
	}

	c(ch)
}

func main() {
	matrixA = PopulateMatrix(6)
	matrixB = PopulateMatrix(6)

	var wg sync.WaitGroup
	fm := make(chan [][]int)

	PrintMatrix(matrixA, "A")
	PrintMatrix(matrixB, "B")

	wg.Add(1)
	go FoxAlgorithm(fm)
	wg.Done()

	PrintMatrix(<-fm, "C(Fox)")
}
