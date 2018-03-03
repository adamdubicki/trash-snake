package main

import (
	"math"
)

func toStringPointer(str string) *string {
	return &str
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func pointInSet(p Point, s []Point) bool {
	for i := 0; i < len(s); i++ {
		if p.X == s[i].X && p.Y == s[i].Y {
			return true
		}
	}

	return false
}

func distance(p1 Point, p2 Point) int {
	return abs(p1.X-p2.X) + abs(p1.Y-p2.Y)
}

func reconstructPath(current Point, pathMap map[Point]Point) []Point {
	path := make([]Point, 0)
	path = append(path, current)

	_, exists := pathMap[current]

	for ; exists; _, exists = pathMap[current] {
		current = pathMap[current]
		path = append(path, current)
	}

	return reverseList(path)
}

func projectSnakeAlongPath(path []Point, snake Snake) []Point {
	if len(path) < len(snake.Body) {
		p := make([]Point, 0)
		p = append(p, path[:len(path)]...)
		p = append(p, snake.Body[:(len(snake.Body)-len(path))+1]...)
		return p
	} else if len(path) > len(snake.Body) {
		return path[:len(snake.Body)]
	}

	return path
}

func pathIsSafe(path []Point, ourSnake Snake, b *Board) bool {
	path = reverseList(path)
	if len(path) < 2 {
		return false
	}

	copy := b.copy()
	for _, v := range ourSnake.Body {
		copy.insert(v, empty())
	}

	projected := projectSnakeAlongPath(path, ourSnake)
	for _, p := range projected {
		copy.insert(p, obstacle())
	}
	fakeHead := projected[0]
	fakeTail := projected[len(projected)-1]
	copy.insert(fakeHead, snakeHead())
	copy.insert(fakeTail, empty())

	pathToTail := shortestPath(fakeHead, fakeTail, copy)
	if len(pathToTail) > 2 {
		return true
	}

	return false
}

func reverseList(lst []Point) []Point {
	for i := 0; i < len(lst)/2; i++ {
		j := len(lst) - i - 1
		lst[i], lst[j] = lst[j], lst[i]
	}
	return lst
}

func getDirection(from Point, to Point) string {
	vertical := to.Y - from.Y
	horizontal := to.X - from.X
	if vertical == 0 {
		if horizontal > 0 {
			return RIGHT
		}
		return LEFT
	}
	if vertical < 0 {
		return UP
	}
	return DOWN
}

func pairIsValidExtension(p1 Point, p2 Point, board Board, path []Point) bool {
	return pointIsValidExtension(p1, board, path) && pointIsValidExtension(p2, board, path)
}

func pointIsValidExtension(p Point, board Board, path []Point) bool {
	return !board.getTile(p).Dangerous && !pointInSet(p, path)
}

func extendPath(path []Point, board Board, limit int) []Point {
	extended := make([]Point, 0)
	extended = append(extended, path...)
	for i := 0; i < len(extended)-1; i++ {
		current := extended[i]
		next := extended[i+1]
		direction := getDirection(current, next)
		if direction == RIGHT || direction == LEFT {
			currentUp := Point{current.X, current.Y - 1}
			currentDown := Point{current.X, current.Y + 1}
			nextUp := Point{next.X, next.Y - 1}
			nextDown := Point{next.X, next.Y + 1}
			if pairIsValidExtension(currentUp, nextUp, board, extended) {
				extended = append(extended[0:i+1], append([]Point{currentUp, nextUp}, extended[i+1:]...)...)
			} else if pairIsValidExtension(currentDown, nextDown, board, extended) {
				extended = append(extended[0:i+1], append([]Point{currentDown, nextDown}, extended[i+1:]...)...)
			}
		} else if direction == UP || direction == DOWN {
			currentLeft := Point{current.X - 1, current.Y}
			currentRight := Point{current.X + 1, current.Y}
			nextLeft := Point{next.X - 1, next.Y}
			nextRight := Point{next.X + 1, next.Y}
			if pairIsValidExtension(currentLeft, nextLeft, board, extended) {
				extended = append(extended[0:i+1], append([]Point{currentLeft, nextLeft}, extended[i+1:]...)...)
			} else if pairIsValidExtension(currentRight, nextRight, board, extended) {
				extended = append(extended[0:i+1], append([]Point{currentRight, nextRight}, extended[i+1:]...)...)
			}
		}
		if i == len(extended)-1 || len(extended) > limit {
			continue
		}
	}
	return extended
}

// Find the shortest path from start -> goal
func shortestPath(start Point, goal Point, board *Board) []Point {
	closedSet := make([]Point, 0)    // Tiles already explored
	openSet := make([]Point, 0)      // Tiles to explore
	openSet = append(openSet, start) // Start exploring from start tile

	gScore := make(map[Point]float32) // Shortest path distance
	fScore := make(map[Point]float32) // Manhatten distance heuristic
	cameFrom := make(map[Point]Point)
	for i := 0; i < board.Width; i++ {
		for j := 0; j < board.Height; j++ {
			gScore[Point{i, j}] = 1000.0
			fScore[Point{i, j}] = 1000.0
		}
	}
	gScore[start] = 0
	fScore[start] = float32(distance(start, goal))

	// While there are still tiles to explore
	for len(openSet) > 0 {
		// Pick the current closest based on the heuristic
		min := openSet[0]
		minIndex := 0
		for i := 0; i < len(openSet); i++ {
			if fScore[openSet[i]] < fScore[min] {
				min = openSet[i]
				minIndex = i
			}
		}
		if min.X == goal.X && min.Y == goal.Y {
			// fmt.Println("got here")
			return reconstructPath(goal, cameFrom)
		}

		// Remove the minimum from the open set, add to closed set
		openSet[minIndex] = openSet[len(openSet)-1]
		openSet = openSet[:len(openSet)-1] // << maybe here?
		closedSet = append(closedSet, min)
		neighbours := board.getValidTiles(min)

		// Explore the neighbours
		for _, n := range neighbours {
			if pointInSet(n, closedSet) {
				continue
			}

			tentativeGScore := gScore[min] + float32(distance(min, n))

			if !pointInSet(n, openSet) {
				openSet = append(openSet, n)
			} else if tentativeGScore >= gScore[n] {
				continue
			}

			cameFrom[n] = min
			gScore[n] = tentativeGScore

			var bonus float32
			if board.getTile(n).EntityType == EMPTY {
				bonus = -0.1
			} else {
				bonus = 0.0
			}

			fScore[n] = tentativeGScore + float32(distance(n, min)) + bonus
		}
	}

	return nil
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func getKillIncentive(direction string, head Point) []Point {
	switch direction {
	case UP:
		return []Point{
			Point{head.X - 1, head.Y - 1},
			Point{head.X, head.Y - 1},
			Point{head.X + 1, head.Y - 1},
		}
	case LEFT:
		return []Point{
			Point{head.X - 1, head.Y - 1},
			Point{head.X - 1, head.Y},
			Point{head.X - 1, head.Y + 1},
		}
	case DOWN:
		return []Point{
			Point{head.X - 1, head.Y + 1},
			Point{head.X, head.Y + 1},
			Point{head.X + 1, head.Y + 1},
		}
	case RIGHT:
		return []Point{
			Point{head.X + 1, head.Y - 1},
			Point{head.X + 1, head.Y},
			Point{head.X + 1, head.Y + 1},
		}
	default:
		return []Point{
			Point{head.X - 1, head.Y - 1},
			Point{head.X, head.Y - 1},
			Point{head.X + 1, head.Y - 1},
		}
	}
}
