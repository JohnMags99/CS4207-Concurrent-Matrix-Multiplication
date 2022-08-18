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

func CanonAlgorithm() [][]int {

	AllignMatrices := func() {

		//Allign A values
		for i := 1; i < len(matrixA); i++ {
			row := make([]int, len(matrixA[i]))
			for j := 0; j < len(matrixA[i]); j++ {
				row[j] = matrixA[i][j]
			}

			firstHalf := row[:i]
			secndHalf := row[i:]

			newRow := []int{}
			newRow = append(secndHalf, firstHalf...)

			for j := 0; j < len(matrixA[i]); j++ {
				matrixA[i][j] = newRow[j]
			}
		}

		//Allign B values
		for i := 1; i < len(matrixB); i++ {
			collumn := make([]int, len(matrixB))
			for j := 0; j < len(matrixB[i]); j++ {
				collumn[j] = matrixB[j][i]
			}

			firstHalf := collumn[:i]
			secndHalf := collumn[i:]

			newCollumn := []int{}
			newCollumn = append(secndHalf, firstHalf...)

			for j := 0; j < len(matrixB[i]); j++ {
				matrixB[j][i] = newCollumn[j]
			}
		}

		//For testing
		/*PrintMatrix(matrixA, "new A")
		PrintMatrix(matrixB, "new B")*/
	}

	//AShift == function to shift AValues left
	AShift := func() [][]int {
		shiftA := CreateEmpty(len(matrixA))
		keys := make([]int, len(matrixB))

		for i := len(matrixA) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(matrixA); i++ {
			for j := 0; j < len(matrixA[i]); j++ {
				shiftA[i][j] = matrixA[i][keys[j]]
			}
		}

		return shiftA
	}

	//BShift == function to shift BValues up
	BShift := func() [][]int {
		shiftB := CreateEmpty(len(matrixB))
		keys := make([]int, len(matrixB))

		for i := len(matrixB) - 2; i >= 0; i-- {
			keys[i] = i + 1
		}

		for i := 0; i < len(keys); i++ {
			shiftB[i] = matrixB[keys[i]]
		}

		return shiftB
	}

	//c == Matrix Multiplication function which returns matrix c
	c := func() [][]int {
		matrixC := CreateEmpty(len(matrixA))

		phases := int(math.Sqrt(float64(len(matrixA) * len(matrixA[0]))))

		AllignMatrices()
		for phase := 0; phase < phases; phase++ {
			matrixA = AShift()
			matrixB = BShift()
			for i := 0; i < len(matrixA); i++ {
				for j := 0; j < len(matrixA[i]); j++ {
					matrixC[i][j] = matrixC[i][j] + (matrixA[i][j] * matrixB[i][j])
				}
			}

		}
		return matrixC
	}

	return c()
}

func main() {
	var wg sync.WaitGroup

	fm := make(chan [][]int)
	//cm := make(chan [][]int)

	matrixA = PopulateMatrix(4)
	PrintMatrix(matrixA, "A")

	matrixB = PopulateMatrix(4)
	PrintMatrix(matrixB, "B")

	wg.Add(1)
	go FoxAlgorithm(fm)
	PrintMatrix(<-fm, "C(Fox)")
	wg.Done()

	PrintMatrix(CanonAlgorithm(), "C(Cannon)")
}
