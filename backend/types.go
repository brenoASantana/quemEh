package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type GameState string

const (
	StateLobby    GameState = "LOBBY"
	StateAnswer   GameState = "ANSWERING"
	StateVoting   GameState = "VOTING"
	StateResults  GameState = "RESULTS"
	StateGameOver GameState = "GAME_OVER"
)

type Player struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	IsHost bool            `json:"isHost"`
	Score  int             `json:"score"`
	Conn   *websocket.Conn `json:"-"`
	Room   *Room           `json:"-"`
}

type Room struct {
	ID              string
	Players         map[string]*Player
	State           GameState
	Question        string
	Answers         map[string]string
	VotedGuesses    map[string]map[int]string
	ShuffledAnswers []AnswerWithIndex
	Broadcast       chan Message
	Mutex           sync.RWMutex
	QuestionIndex   int
	QuestionDeck    []string
}

type AnswerWithIndex struct {
	Text      string `json:"text"`
	RealID    string `json:"realId"`
	DisplayID int    `json:"id"`
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type PublicGameState struct {
	Players         []*Player                 `json:"players"`
	State           GameState                 `json:"state"`
	Question        string                    `json:"question"`
	AnswersCount    int                       `json:"answersCount"`
	Answers         []AnswerWithIndex         `json:"answers,omitempty"`
	VotedGuesses    map[string]map[int]string `json:"votedGuesses,omitempty"`
	LastRoundResult map[string]int            `json:"lastRoundResult,omitempty"`
}
