// A little "Connect Four" game

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Width     = 7
	Height    = 6
	LookAhead = 6

	// Piece
	Empty = 0
	You   = 1
	Me    = 2

	// Endings
	Continue = 0
	Tie      = 1
	Win      = 2
)

type (
	Piece  uint8
	Ending int
)

var (
	grid    [Width * Height]Piece
	scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	quiet   bool           = true
)

func main() {
	flag.BoolVar(&quiet, "q", false, "quiet mode (don't show board on stderr)")
	flag.Parse()

	exitCode := 0
	for {
		draw()

		// Their move
		for {
			vprintf("Enter your move column (0..6): ")
			move := readMove()
			if move < 0 {
				continue
			}
			if placeMove(move, You) {
				break
			}
			vprintf("Couldn't make move at column %d\n", move)
		}
		end := getEnding(You)
		if end == Tie {
			vprintf("Tie after your move\n")
			exitCode = 3
			break
		} else if end == Win {
			vprintf("You won!\n")
			exitCode = 2
			break
		}

		// Our move
		move := makeMove()
		end = getEnding(Me)
		if end == Tie {
			vprintf("Tie after my move\n")
			exitCode = 3
			break
		} else if end == Win {
			vprintf("I won!\n")
			exitCode = 1
			break
		}
		vprintf("My move: ")
		fmt.Printf("%d\n", move)
	}

	draw()

	os.Exit(exitCode)
}

func vprintf(format string, args ...interface{}) {
	if !quiet {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func put(x, y int, p Piece) {
	grid[y*Width+x] = p
}

func get(x, y int) Piece {
	return grid[y*Width+x]
}

func draw() {
	if quiet {
		return
	}
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
		vprintf("%d ", x)
	}
	fmt.Fprint(os.Stderr, "|\n")
}

func readMove() int {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			vprintf("error reading standard input: %v\n", err)
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

func liftMove(x int) {
	for y := 0; y < Height; y++ {
		if get(x, y) != Empty {
			put(x, y, Empty)
			return
		}
	}
	panic(fmt.Sprintf("invalid liftMove() %d", x))
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
	move, _ := pickMove(Me, LookAhead)
	if move >= 0 {
		placeMove(move, Me)
	}
	return move
}

func pickMove(piece Piece, lookahead int) (move, score int) {
	scores := make([]int, Width)
	for i := range scores {
		if !placeMove(i, piece) {
			scores[i] = -20000
			continue
		}

		end := getEnding(piece)
		if end == Win {
			liftMove(i)
			return i, 10000
		} else if end == Tie {
			// scores[i] is 0 already
			liftMove(i)
			continue
		}

		if lookahead > 0 {
			other := Piece(You)
			if piece == You {
				other = Me
			}
			_, otherScore := pickMove(other, lookahead-1)
			scores[i] = -otherScore
		} else {
			scores[i] = -getScore(piece)
		}

		liftMove(i)
	}

	highest := -20000
	highestIndex := -1
	for i := range scores {
		if scores[i] > highest {
			highest = scores[i]
			highestIndex = i
		}
	}
	if highest == -20000 {
		return -1, highest
	}
	return highestIndex, highest
}

var runScores = map[int]int{
	1: 1,
	2: 10,
	3: 100,
	4: 1000,
}

func getScore(piece Piece) int {
	type pos struct{ x, y int }

	counted := make(map[pos]bool)
	score := 0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			p := get(x, y)
			if p == Empty {
				continue
			}
			if _, ok := counted[pos{x, y}]; ok {
				// Already counted
				continue
			}

			// Scan for run to the right
			run := 1
			for i := 1; i < 4 && x+i < Width; i++ {
				q := get(x+i, y)
				if q != p {
					if q == Empty {
						run++
					}
					break
				}
				counted[pos{x, y}] = true
				run++
			}
			score += runScores[run]

			// Scan for run down and to the right
			run = 1
			for i := 1; i < 4 && x+i < Width && y+i < Height; i++ {
				q := get(x+i, y+i)
				if q != p {
					if q == Empty {
						run++
					}
					break
				}
				counted[pos{x, y}] = true
				run++
			}
			score += runScores[run]

			// Scan for run down
			run = 1
			for i := 1; i < 4 && y+i < Height; i++ {
				q := get(x, y+i)
				if q != p {
					if q == Empty {
						run++
					}
					break
				}
				counted[pos{x, y}] = true
				run++
			}
			score += runScores[run]

			// Scan for run down and to the left
			run = 1
			for i := 1; i < 4 && x-i >= 0 && y+i < Height; i++ {
				q := get(x-i, y+i)
				if q != p {
					if q == Empty {
						run++
					}
					break
				}
				counted[pos{x, y}] = true
				run++
			}
			score += runScores[run]
		}
	}

	return score
}
