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
		Taunt:          "I eat garbage.",
		HeadType:       HEAD_SAND_WORM,
		TailType:       TAIL_FRECKLED,
		SecondaryColor: "pink",
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

	foodChannel := make(chan string)
	tailChannel := make(chan string)
	optimalChannel := make(chan string)

	var foodResult string
	var tailResult string
	var optimalResult string
	// fmt.Println(bm.Req.Food)

	go func() {
		foodMove := NO_MOVE
		food := bm.findBestFood()
		if food.X != -1 {
			foodPath := shortestPath(bm.OurHead, food, bm.GameBoard)
			if len(foodPath) >= 2 && pathIsSafe(foodPath, bm.Req.You, bm.GameBoard) {
				foodPath = reverseList(foodPath)
				foodMove = getDirection(foodPath[0], foodPath[1])
			}
		}
		foodChannel <- foodMove
	}()

	go func() {
		tailMove := NO_MOVE
		copy := bm.GameBoard.copy()
		copy.insert(bm.Req.You.Tail(), food())
		tailPath := shortestPath(bm.OurHead, bm.Req.You.Tail(), copy)
		if len(tailPath) >= 2 {
			if len(tailPath) == 2 && bm.Req.You.Health > 99 {
				tailMove = NO_MOVE
			} else {
				if len(bm.Req.You.Body) > len(tailPath) {
					extendPath := extendPath(tailPath, *bm.GameBoard)
					tailMove = getDirection(extendPath[0], extendPath[1])
				} else {
					tailMove = getDirection(tailPath[0], tailPath[1])
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
			// fmt.Println("Food Result:", foodResult, time.Since(start))
			continue
		case tailResult = <-tailChannel:
			// fmt.Println("Tail Result:", tailResult, time.Since(start))
			continue
		case optimalResult = <-optimalChannel:
			// fmt.Println("Optimal Result:", optimalResult, time.Since(start))
			continue
		}
	}
	// fmt.Println("---------------------")

	if foodResult != NO_MOVE {
		currentMove = foodResult
	} else if tailResult != NO_MOVE {
		currentMove = tailResult
	} else if optimalResult != NO_MOVE {
		currentMove = optimalResult
	}

	respond(res, MoveResponse{
		Move: currentMove,
	})
}
