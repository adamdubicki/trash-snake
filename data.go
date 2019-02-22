package main

import (
	"encoding/json"
	"net/http"
)

type StartRequest struct {
	GameID int `json:"game_id"`
}

func NewStartRequest(req *http.Request) (*StartRequest, error) {
	decoded := StartRequest{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	return &decoded, err
}

type StartResponse struct {
	Color          string   `json:"color,omitempty"`
	Name           string   `json:"name,omitempty"`
	HeadURL        string   `json:"head_url,omitempty"`
	Taunt          string   `json:"taunt,omitempty"`
	HeadType       HeadType `json:"head_type,omitempty"`
	TailType       TailType `json:"tail_type,omitempty"`
	SecondaryColor string   `json:"secondary_color,omitempty"`
}

type Game struct {
	ID string `json:"id"`
}

type EmptyResponse struct{}

type Board struct {
	Height int     `json:"height"`
	Width  int     `json:"width"`
	Food   []Coord `json:"food"`
	Snakes []Snake `json:"snakes"`
	Grid   [][]entity
}

type MoveRequest struct {
	Game  Game  `json:"game"`
	Turn  int   `json:"turn"`
	Board Board `json:"board"`
	You   Snake `json:"you"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Snake - Destructured Snake object for request
type Snake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int     `json:"health"`
	Body   []Coord `json:"body"`
}

// Board - Destructured Board object for reques

func NewMoveRequest(req *http.Request) (*MoveRequest, error) {
	decoded := MoveRequest{}
	err := json.NewDecoder(req.Body).Decode(&decoded)
	return &decoded, err
}

type MoveResponse struct {
	Move string `json:"move"`
}

func (snake Snake) Head() Coord { return snake.Body[0] }

type HeadType string

const (
	HEAD_BENDR     HeadType = "bendr"
	HEAD_DEAD               = "dead"
	HEAD_FANG               = "fang"
	HEAD_PIXEL              = "pixel"
	HEAD_REGULAR            = "regular"
	HEAD_SAFE               = "safe"
	HEAD_SAND_WORM          = "sand-worm"
	HEAD_SHADES             = "shades"
	HEAD_SMILE              = "smile"
	HEAD_TONGUE             = "tongue"
)

type TailType string

const (
	TAIL_BLOCK_BUM    = "block-bum"
	TAIL_CURLED       = "curled"
	TAIL_FAT_RATTLE   = "fat-rattle"
	TAIL_FRECKLED     = "freckled"
	TAIL_PIXEL        = "pixel"
	TAIL_REGULAR      = "regular"
	TAIL_ROUND_BUM    = "round-bum"
	TAIL_SKINNY       = "skinny"
	TAIL_SMALL_RATTLE = "small-rattle"
)

// Parses List<Coord> into []Coord
type PointList []Coord

func (list *PointList) UnmarshalJSON(data []byte) error {
	var obj struct {
		Data []Coord `json:"data"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*list = obj.Data
	return nil
}

// Parses List<Snake> into []Snake
type SnakeList []Snake

func (list *SnakeList) UnmarshalJSON(data []byte) error {
	var obj struct {
		Data []Snake `json:"data"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*list = obj.Data
	return nil
}

func (snake Snake) Tail() Coord { return snake.Body[len(snake.Body)-1] }
