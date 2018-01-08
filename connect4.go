// A little "Connect Four" game

package main

import (
	"bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
)

const (
	Width = 7
	Height = 6

	// Piece
	Empty = 0
	You = 1
	Me = 2

	// Endings
	Continue = 0
	Tie = 1
	Win = 2
)

type (
	Piece uint8
	Ending int
)

var (
	grid [Width*Height]Piece
	scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
)

func main() {
	for {
		draw()

		// Their move
		for {
			fmt.Fprintf(os.Stderr, "Enter your move column (0..6): ")
			move := readMove()
			if move < 0 {
				continue
			}
			if placeMove(move, You) {
				break
			}
			fmt.Fprintf(os.Stderr, "Couldn't make move at column %d\n", move)
		}
		end := getEnding(You)
		if end == Tie {
			fmt.Fprintf(os.Stderr, "Tie after your move\n")
			break
		} else if end == Win {
			fmt.Fprintf(os.Stderr, "You won!\n")
			break
		}

		// Our move
		move := makeMove()
		end = getEnding(Me)
		if end == Tie {
			fmt.Fprintf(os.Stderr, "Tie after my move\n")
			break
		} else if end == Win {
			fmt.Fprintf(os.Stderr, "I won!\n")
			break
		}
		fmt.Fprintf(os.Stderr, "My move: ")
		fmt.Printf("%d\n", move)
	}

	draw()
}

func put(x, y int, p Piece) {
	grid[y*Width+x] = p
}

func get(x, y int) Piece {
	return grid[y*Width+x]
}

func draw() {
	for y := 0; y < Height; y++ {
		fmt.Fprint(os.Stderr, "| ")
		for x := 0; x < Width; x++ {
			switch get(x, y) {
			case Empty:
				fmt.Fprint(os.Stderr, ". ")
			case Me:
				fmt.Fprint(os.Stderr, "M ")
			case You:
				fmt.Fprint(os.Stderr, "Y ")
			}
		}
		fmt.Fprint(os.Stderr, "|\n")
	}
	fmt.Fprint(os.Stderr, "+-", strings.Repeat("-", Width*2), "+\n")
	fmt.Fprint(os.Stderr, "| ")
	for x := 0; x < Width; x++ {
		fmt.Fprintf(os.Stderr, "%d ", x)
	}
	fmt.Fprint(os.Stderr, "|\n")
}

func readMove() int {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
		    fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		return -1
	}
	move, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return -1
	}
	if move < 0 || move >= Width {
		return -1
	}
	return move
}

func placeMove(x int, p Piece) bool {
	for y := 0; y < Height; y++ {
		if get(x, y) != Empty {
			if y == 0 {
				return false
			}
			put(x, y-1, p)
			return true
		}
	}
	put(x, Height-1, p)
	return true
}

func getEnding(p Piece) Ending {
	numEmpty := 0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if get(x, y) == Empty {
				numEmpty++
			}
			if get(x, y) != p {
				continue
			}
			if x <= Width-4 && get(x+1, y) == p && get(x+2, y) == p && get(x+3, y) == p {
				// Winning four to the right
				return Win
			}
			if x <= Width-4 && y <= Height-4 && get(x+1, y+1) == p && get(x+2, y+2) == p && get(x+3, y+3) == p {
				// Winning four down and to the right
				return Win
			}
			if y <= Height-4 && get(x, y+1) == p && get(x, y+2) == p && get(x, y+3) == p {
				// Winning four down
				return Win
			}
			if x >= 3 && y <= Height-4 && get(x-1, y+1) == p && get(x-2, y+2) == p && get(x-3, y+3) == p {
				// Winning four down and to the left
				return Win
			}
		}
	}
	if numEmpty == 0 {
		return Tie
	}
	return Continue
}

func makeMove() int {
	possibilities := make([]int, 0, Width)
	for x := 0; x < Width; x++ {
		if get(x, 0) == Empty {
			possibilities = append(possibilities, x)
		}
	}
	if len(possibilities) == 0 {
		return -1
	}
	index := rand.Intn(len(possibilities))
	move := possibilities[index]
	if !placeMove(move, Me) {
		panic(fmt.Sprintf("invalid makeMove() %d", move))
	}	
	return move
}
