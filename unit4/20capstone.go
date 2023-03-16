package main
import (
	"fmt"
	"math/rand"
	"time"
)

const (
	width = 80
	height = 15
)

type Universe [][]bool

func NewUniverse() Universe {
	toReturn := make([][]bool, height)
	for i := range toReturn {
		toReturn[i] = make([]bool, width)
	}
	return toReturn
}

func (u Universe) Show() {
	for i := 0; i < width + 2; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
	for row := range u {
		fmt.Printf("-")
		for col := range u[row] {
			if u[row][col] {
				fmt.Printf("*")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("-")
		fmt.Printf("\n")
	}
	for i := 0; i < width + 2; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("\n")
}

func (u Universe) Seed() {
	rand.Seed(time.Now().UnixNano())
	for row := range u {
		for col := range u[row] {
			randomFrac := rand.Float64()
			if randomFrac < 0.25 {
				u[row][col] = true
			} else {
				u[row][col] = false
			}
		}
	}
}

func (u Universe) Alive(row, col int) bool {
	row = row % width
	col = col % width
}

func (u Universe) Neighbors(x, y int) int {
	neighborCnt := 0
	if u.Alive(x + 1, y - 1) {
		neighborCnt++
	}
	if u.Alive(x + 1, y) {
		neighborCnt++
	}
	if u.Alive(x + 1, y + 1) {
		neighborCnt++
	}

	if u.Alive(x - 1, y + 1) {
		neighborCnt++
	}
	if u.Alive(x - 1, y) {
		neighborCnt++
	}
	if u.Alive(x - 1, y - 1) {
		neighborCnt++
	}

	if u.Alive(x, y + 1) {
		neighborCnt++
	}
	if u.Alive(x, y - 1) {
		neighborCnt++
	}
	return neighborCnt
}

func (u Universe) Next(x, y int) bool {
	neighbors := u.Neighbors(x, y)
	if u.Alive(x, y) {
		//the cell in question is alive originally
		if neighbors == 2 || neighbors == 3 {
			//cell stays alive
			return true
		} else {
			//cell dies
			return false
		}
	} else {
		//the cell in question is dead originally
		if neighbors == 3 {
			//revive the cell cause it has exactly 3 neighbors
			return true
		} else {
			//otherwise stays dead
			return false
		}
	}
}

func main() {
	tester := NewUniverse()
	tester.Show()
	tester.Seed()
	tester.Show()
}