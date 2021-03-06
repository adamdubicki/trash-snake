package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respond(res http.ResponseWriter, obj interface{}) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(obj)
}

func handlePing(res http.ResponseWriter, req *http.Request) {
	respond(res, EmptyResponse{})
}

func handleStart(res http.ResponseWriter, req *http.Request) {
	respond(res, StartResponse{
		Color:          "orange",
		Name:           "trash-snake",
		Taunt:          "ITS MY CHARACTER! I'M THE TRASH SNAKE!",
		HeadType:       HEAD_SAND_WORM,
		TailType:       TAIL_FRECKLED,
		SecondaryColor: "pink",
		HeadURL:        "https://pbs.twimg.com/profile_images/535222646963572736/KZItD1f-_400x400.png",
	})
}

func handleMove(res http.ResponseWriter, req *http.Request) {
	currentMove := "down"
	data, err := NewMoveRequest(req)
	if err != nil {
		fmt.Println("ERROR: ", err)
		respond(res, MoveResponse{
			Move: "up",
		})
		return
	}
	bm := initializeBoard(data)

	foodChannel := make(chan Coord)
	tailChannel := make(chan Coord)
	optimalChannel := make(chan string)

	var foodResult Coord
	var tailResult Coord
	var optimalResult string

	// start := time.Now()
	// bm.GameBoard.show()
	go func() {
		foodMove := Coord{-1, -1}
		foodResult := bm.findBestFood()
		if foodResult.Food.X != -1 {
			foodPath := shortestPath(bm.OurHead, foodResult.Food, bm.GameBoard)

			var pathIsSafeCheck bool
			if len(bm.Req.You.Body) <= 2 {
				pathIsSafeCheck = true
			} else {
				pathIsSafeCheck = pathIsSafe(foodPath, bm.Req.You, bm.GameBoard)
				foodPath = reverseList(foodPath)
			}

			if len(foodPath) >= 2 && pathIsSafeCheck { // here
				foodMove = foodPath[1]
			}
		}
		foodChannel <- foodMove
	}()

	go func() {
		tailMove := Coord{-1, -1}
		copy := bm.GameBoard.copy()
		copy.insert(bm.Req.You.Tail(), food())
		tailPath := shortestPath(bm.OurHead, bm.Req.You.Tail(), copy)
		if len(tailPath) >= 2 && len(bm.Req.You.Body) > 2 {
			if len(tailPath) == 2 && bm.Req.You.Health > 99 {
				tailMove = Coord{-1, -1}
			} else {
				if bm.Req.Turn > 5 {
					extendPath := extendPath(tailPath, *bm.GameBoard, 50)
					tailMove = extendPath[1]
				} else {
					tailMove = tailPath[1]
				}
			}
		}

		tailChannel <- tailMove
	}()

	go func() {
		root := createRoot(data)
		currentMove := NO_MOVE
		if len(bm.Req.You.Body) >= 2 {
			currentMove = root.getOptimalMove()
		}
		optimalChannel <- currentMove
	}()

	// fmt.Println("---------------------")
	for i := 0; i < 3; i++ {
		select {
		case foodResult = <-foodChannel:
			// fmt.Println("Food Result:", getDirection(bm.Req.You.Head(), foodResult), time.Since(start))
			continue
		case tailResult = <-tailChannel:
			// fmt.Println("Tail Result:", getDirection(bm.Req.You.Head(), tailResult), time.Since(start))
			continue
		case optimalResult = <-optimalChannel:
			// fmt.Println("Optimal Result:", optimalResult, time.Since(start))
			continue
		}
	}
	// fmt.Println("---------------------")

	// fmt.Println(foodResult)
	if bm.GameBoard.tileInBounds(foodResult) {
		currentMove = getDirection(bm.Req.You.Head(), foodResult)
	} else if bm.GameBoard.tileInBounds(tailResult) && bm.Req.Turn > 3 {
		currentMove = getDirection(bm.Req.You.Head(), tailResult)
	} else if optimalResult != NO_MOVE {
		currentMove = optimalResult
	} else {
		neighbours := bm.GameBoard.getValidTiles(bm.Req.You.Head())
		if len(neighbours) > 0 {
			// fmt.Println(neighbours)
			currentMove = getDirection(bm.Req.You.Head(), neighbours[0])
		} else {
			currentMove = UP
		}
	}

	respond(res, MoveResponse{
		Move: currentMove,
	})
}
