package main

// Tree node in our min max tree
type MinMax struct {
	Width  int
	Height int
	Snakes map[string][]Point
	YouID  string
}

func createRoot(req *MoveRequest) MinMax {
	root := MinMax{0, 0, make(map[string][]Point), req.You.ID}
	for _, snake := range req.Snakes {
		newBody := append([]Point{}, snake.Body[:len(snake.Body)]...)
		root.Snakes[snake.ID] = newBody
	}
	root.Height = req.Height
	root.Width = req.Width
	return root
}

func (m MinMax) getOptimalMove() string {
	permutations := m.generatePermutations(m.YouID)
	move := UP
	if len(permutations) > 0 {
		score := permutations[0].recursiveScore(1)
		move = permutations[0].direction(m.YouID)
		for _, k := range permutations {
			recursiveScore := k.recursiveScore(1)
			if score < recursiveScore {
				score = recursiveScore
				// fmt.Println("BEST score", score
				move = k.direction(m.YouID)
			}
		}
	} else {
		move = UP
	}

	return move
}

func (m MinMax) score() int {
	board := createBoard(m.Width, m.Height)
	for _, snakes := range m.Snakes {
		for _, point := range snakes {
			board.insert(point, obstacle())
		}
	}
	score := len(board.getValidTiles(m.Snakes[m.YouID][0]))
	return score
}

func (m MinMax) recursiveScore(depth int) int {
	if depth > 3 {
		return m.score()
	} else if depth%2 == 0 {
		permutations := m.generatePermutations(m.YouID)
		if len(permutations) == 0 {
			return m.score()
		} else {
			score := getBestScore(permutations, depth)
			return score
		}
	} else {
		permutations := make([]MinMax, 0)
		for k := range m.Snakes {
			permutations = append(permutations, m.generatePermutations(k)...)
		}
		if len(permutations) == 0 {
			return m.score()
		} else {
			score := getBestScore(permutations, depth)
			return score
		}
	}
}

func getBestScore(p []MinMax, depth int) int {
	score := p[0].recursiveScore(depth + 1)
	for _, p := range p {
		if score < p.score() {
			score = p.recursiveScore(depth + 1)
		}
	}
	return score
}

func (m MinMax) direction(snakeID string) string {
	return getDirection(m.Snakes[snakeID][1], m.Snakes[snakeID][0])
}

func (m MinMax) generatePermutations(id string) []MinMax {
	board := createBoard(m.Width, m.Height)
	for _, snakes := range m.Snakes {
		for _, point := range snakes {
			board.insert(point, obstacle())
		}
	}

	neighbours := board.getValidTiles(m.Snakes[id][0])
	permuations := make([]MinMax, 0)
	for _, n := range neighbours {
		permutation := m.copy()
		permutation.Snakes[id] = append([]Point{n}, permutation.Snakes[id][:len(permutation.Snakes[id])-1]...)
		permuations = append([]MinMax{permutation}, permuations...)
	}

	return permuations
}

func (m MinMax) copy() MinMax {
	new := MinMax{m.Width, m.Height, make(map[string][]Point), m.YouID}
	for k, v := range m.Snakes {
		new.Snakes[k] = v
	}

	return new
}
