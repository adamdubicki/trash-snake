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

	foodChannel := make(chan Point)
	tailChannel := make(chan Point)
	optimalChannel := make(chan string)

	var foodResult Point
	var tailResult Point
	var optimalResult string

	// start := time.Now()
	// bm.GameBoard.show()
	go func() {
		foodMove := Point{-1, -1}
		foodResult := bm.findBestFood()
		if foodResult.Food.X != -1 {
			foodPath := shortestPath(bm.OurHead, foodResult.Food, bm.GameBoard)
			if len(foodPath) >= 2 && pathIsSafe(foodPath, bm.Req.You, bm.GameBoard) {
				foodPath = reverseList(foodPath)
				foodMove = foodPath[1]
			}
		}
		foodChannel <- foodMove
	}()

	go func() {
		tailMove := Point{-1, -1}
		copy := bm.GameBoard.copy()
		copy.insert(bm.Req.You.Tail(), food())
		tailPath := shortestPath(bm.OurHead, bm.Req.You.Tail(), copy)
		if len(tailPath) >= 2 {
			if len(tailPath) == 2 && bm.Req.You.Health > 99 {
				tailMove = Point{-1, -1}
			} else {
				if bm.Req.Turn > 5 {
					extendPath := extendPath(tailPath, *bm.GameBoard, 15)
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
		optimalChannel <- root.getOptimalMove()
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
		// fmt.Println("WENT FOR OPTIMAL")
		currentMove = optimalResult
	} else {
		neighbours := bm.GameBoard.getValidTiles(bm.Req.You.Head())
		if len(neighbours) > 0 {
			currentMove = getDirection(bm.Req.You.Head(), neighbours[0])
		} else {
			currentMove = UP
		}
	}

	respond(res, MoveResponse{
		Move: currentMove,
	})
}
