package main

const (
	UP      = "up"
	DOWN    = "down"
	LEFT    = "left"
	RIGHT   = "right"
	NO_MOVE = "no_move"
)

// BoardManager for the board object
type BoardManager struct {
	GameBoard *Board
	Req       *MoveRequest
	OurHead   Point
}

// Fill the board in based on JSON from request
func initializeBoard(req *MoveRequest) *BoardManager {
	bm := new(BoardManager)
	bm.Req = req
	bm.GameBoard = createBoard(req.Width, req.Height)
	bm.addFood(req.Food)
	bm.OurHead = bm.addSnakes(req.Snakes, req.You.ID)

	return bm
}

// Add the food from the JSON
func (bm BoardManager) addFood(foodPoint []Point) {
	for _, element := range foodPoint {
		bm.GameBoard.insert(element, food())
	}
}

func (bm BoardManager) avgSnakeLength() float64 {
	avg := 0.0
	for _, snake := range bm.Req.Snakes {
		avg += float64(len(snake.Body))
	}
	return avg / float64(len(bm.Req.Snakes))
}

// Add our snake and the opposing snakes - with heuristic tiles
func (bm BoardManager) addSnakes(snakePoint []Snake, you string) Point {
	// Add each snake body segment
	for _, snake := range snakePoint {
		for _, snakeBody := range snake.Body {
			bm.GameBoard.insert(snakeBody, obstacle())
		}
	}

	ourHead := Point{}

	for _, snake := range snakePoint {
		if snake.ID == you {
			bm.GameBoard.insert(snake.Head(), snakeHead())
			ourHead = snake.Head()
		} else {
			if distance(snake.Head(), bm.Req.You.Head()) < 5 {
				potential := []Point{
					Point{snake.Head().X - 1, snake.Head().Y},
					Point{snake.Head().X + 1, snake.Head().Y},
					Point{snake.Head().X, snake.Head().Y - 1},
					Point{snake.Head().X, snake.Head().Y + 1},
					Point{snake.Head().X - 1, snake.Head().Y - 1},
					Point{snake.Head().X + 1, snake.Head().Y + 1},
					Point{snake.Head().X - 1, snake.Head().Y + 1},
					Point{snake.Head().X + 1, snake.Head().Y - 1},
				}
				for _, p := range potential {
					if (bm.GameBoard.tileInBounds(p)) && bm.GameBoard.getTile(p).EntityType == EMPTY {
						bm.GameBoard.insert(p, obstacle())
					}
				}
			}
		}

		if snake.Health != 100 && bm.Req.Turn > 5 {
			bm.GameBoard.insert(snake.Tail(), empty())
		}
	}

	return ourHead
}

// Find the best food, the one we are closest
// to compared to all other snakes
func (bm BoardManager) findBestFood() Point {

	best := make(map[Point]Point)

	for _, food := range bm.Req.Food {
		if distance(food, bm.OurHead) < bm.Req.You.Health {
			for _, snake := range bm.Req.Snakes {
				_, exists := best[food]
				if exists == true {
					if distance(best[food], food) > distance(snake.Head(), food) {
						best[food] = snake.Head()
					}
				} else {
					best[food] = snake.Head()
				}
			}
		}
	}

	bestFood := Point{-1, -1}
	for food := range best {
		if best[food] == bm.OurHead {
			if bestFood.X == -1 {
				bestFood = food
			} else {
				if distance(bestFood, bm.OurHead) > distance(food, bm.OurHead) {
					bestFood = food
				}
			}
		}
	}

	return bestFood
}
