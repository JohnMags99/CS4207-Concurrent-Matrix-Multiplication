package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var matrixA, matrixB [][]int

// PopulateMatrix == populates matrix with random ints
func PopulateMatrix(size int) [][]int {
	matrix := make([][]int, size)
	rand.Seed(time.Now().UnixNano())
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

func FoxAlgorithm() [][]int {
	var q int //q == sq root of total matrix size/processes

	//bCastA == function which broadcasts aValues each stage to be multiplied with bValues
	bCastA := func(A [][]int, stage int) [][]int {
		bCast := CreateEmpty(len(A))

		for i := 0; i < len(A); i++ {
			kBar := (i + stage) % q
			for j := 0; j < len(A[i]); j++ {
				bCast[i][j] = A[i][kBar]
			}
		}
		return bCast
	}

	//bShift == function to shift up rows in matrix b after each stage
	bShift := func() {
		shiftB := CreateEmpty(len(matrixB))
		for i := 0; i < len(matrixB); i++ {
			for j := 0; j < len(matrixB[i]); j++ {
				shiftB[i][j] = matrixB[i][j]
			}
		}

		keys := make([]int, len(matrixB))
		for i := len(matrixB) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(keys); i++ {
			matrixB[i] = shiftB[keys[i]]
		}
	}

	//c == Matrix Multiplication function which returns matrix c
	c := func() [][]int {
		matrixC := CreateEmpty(len(matrixA))
		q = int(math.Sqrt(float64(len(matrixA) * len(matrixA[0]))))

		for stage := 0; stage < q; stage++ {
			process := bCastA(matrixA, stage)
			for i := 0; i < len(matrixA); i++ {
				for j := 0; j < len(matrixA[i]); j++ {
					matrixC[i][j] = matrixC[i][j] + (process[i][j] * matrixB[i][j])
				}
			}
			bShift()
		}
		return matrixC
	}

	return c()
}

func main() {
	//var wg sync.WaitGroup

	matrixA = PopulateMatrix(4)
	matrixB = PopulateMatrix(4)

	PrintMatrix(matrixA, "A")
	PrintMatrix(matrixB, "B")

	PrintMatrix(FoxAlgorithm(), "C(Fox)")
}
