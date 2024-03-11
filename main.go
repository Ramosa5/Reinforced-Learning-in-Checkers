package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type State [8][8]int
type Move [2][2]int
type QValueMap map[Move]float64

var QTable map[State]QValueMap
var ExplorationRate float64 = 0.4 // epsilon
var LearningRate float64 = 0.1    // alpha
var DiscountFactor float64 = 0.9  // gamma
var board State

func NewBoard() {
	board = State{
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{3, 0, 3, 0, 3, 0, 3, 0},
		{0, 3, 0, 3, 0, 3, 0, 3},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0}}
}

func PrintState() {
	fmt.Printf("\n")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Printf("%v ", board[i][j])
		}
		fmt.Printf("\n")
	}
}

func OnBoard(vert int, horiz int) bool {
	if vert > 7 || vert < 0 || horiz > 7 || horiz < 0 {
		return false
	}
	return true
}

func Direction(player int) (int, int) {
	vert := -1
	if player == 2 {
		vert = 1
	}
	vertCap := -2
	if player == 2 {
		vertCap = 2
	}
	return vert, vertCap
}

func PossibleMove(player int, mov1 [2]int, mov2 [2]int, capOnly bool, noLog bool) bool {

	if !OnBoard(mov1[0], mov1[1]) || !OnBoard(mov2[0], mov2[1]) {
		return false
	}

	enemy := 2
	if player == 2 {
		enemy = 1
	}

	if board[mov1[0]][mov1[1]] != player && board[mov1[0]][mov1[1]] != player+6 {
		return false
	}
	if board[mov2[0]][mov2[1]] != 3 {
		return false
	}

	if CapCheck(player, mov1, mov2, capOnly, enemy, false) {
		return true
	} else if board[mov1[0]][mov1[1]] == player+6 && CapCheck(player, mov1, mov2, capOnly, enemy, true) {
		return true
	}

	return false
}

func CapCheck(player int, mov1 [2]int, mov2 [2]int, capOnly bool, enemy int, kingFlip bool) bool {
	horiz := mov1[1] - mov2[1]
	vert, vertCap := Direction(player)
	if kingFlip {
		vert = -vert
		vertCap = -vertCap
	}

	if mov2[0]-mov1[0] == vert && math.Abs(float64(horiz)) == 1 && !capOnly {
		return true
	}

	if mov2[0]-mov1[0] == vertCap {
		if horiz == -2 && OnBoard(mov1[0]+vert, mov1[1]+1) {
			if board[mov1[0]+vert][mov1[1]+1] == enemy || board[mov1[0]+vert][mov1[1]+1] == enemy+6 {
				return true
			}
		} else if horiz == 2 && OnBoard(mov1[0]+vert, mov1[1]-1) {
			if board[mov1[0]+vert][mov1[1]-1] == enemy || board[mov1[0]+vert][mov1[1]-1] == enemy+6 {
				return true
			}
		}
	}
	return false
}

func PossibleMoves(player int) [][2][2]int {
	var valids [][2][2]int
	vert, vertCap := Direction(player)

	for i, row := range board {
		for j, element := range row {
			if element != player && element != player+6 {
				continue
			}
			checkMoves := [4][2]int{
				{i + vert, j + 1},
				{i + vert, j - 1},
				{i + vertCap, j + 2},
				{i + vertCap, j - 2},
			}
			currentPos := [2]int{i, j}
			for _, move := range checkMoves {
				if PossibleMove(player, currentPos, move, false, true) {
					valids = append(valids, [2][2]int{currentPos, move})
				}
			}

			if element == player+6 {
				kingMoves := [4][2]int{
					{i - vert, j + 1},
					{i - vert, j - 1},
					{i - vertCap, j + 2},
					{i - vertCap, j - 2},
				}
				for _, move := range kingMoves {
					if PossibleMove(player, currentPos, move, false, true) {
						valids = append(valids, [2][2]int{currentPos, move})
					}
				}
			}
		}
	}
	return valids
}

func ABPlayer(player int, capOnly bool) {
	possibleMoves := PossibleMoves(player)
	if len(possibleMoves) == 0 {
		return
	}
	newBoard := board
	_, move := AlphaBeta(newBoard, 5, -1000, 1000, true, 2)
	//fmt.Printf("Max value %v \n", valu)
	board = ApplyMove(board, move, player, capOnly)

	//PrintState()
	//time.Sleep(4 * time.Second)
}

func EvaluateBoard(board State, player int) int {
	score := 0
	opponent := 3 - player

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] == player {
				score += 1
			} else if board[i][j] == player+6 {
				score += 4
			} else if board[i][j] == opponent {
				score -= 1
			} else if board[i][j] == opponent+6 {
				score -= 4
			}
		}
	}

	return score
}

func ApplyMove(board State, move Move, player int, capOnly bool) State {
	mov1 := move[0]
	mov2 := move[1]

	if PossibleMove(player, mov1, mov2, capOnly, false) {
		king := false
		if board[mov1[0]][mov1[1]] == player+6 {
			king = true
		} else if (player == 1 && mov2[0] == 0) || (player == 2 && mov2[0] == 7) {
			king = true
		}
		board[mov1[0]][mov1[1]] = 3
		if king {
			board[mov2[0]][mov2[1]] = player + 6
		} else {
			board[mov2[0]][mov2[1]] = player
		}

		if math.Abs(float64(mov2[0]-mov1[0])) == 2 { // Capture
			l := (mov1[0] + mov2[0]) / 2
			h := (mov1[1] + mov2[1]) / 2
			board[l][h] = 3
			_, vertCap := Direction(player)
			if PossibleMove(player, mov2, [2]int{mov2[0] + vertCap, mov2[1] + 2}, true, true) ||
				PossibleMove(player, mov2, [2]int{mov2[0] + vertCap, mov2[1] - 2}, true, true) {
				fmt.Println("Anoter capture possible, please go again")
				if player == 1 {
					chooseAction(board, QTable)
					board = ApplyMove(board, move, player, true)
				} else {
					ABPlayer(player, true)
				}
			} else if king && (PossibleMove(player, mov2, [2]int{mov2[0] - vertCap, mov2[1] + 2}, true, true) ||
				PossibleMove(player, mov2, [2]int{mov2[0] - vertCap, mov2[1] - 2}, true, true)) {
				fmt.Println("Anoter capture possible, please go again")
				if player == 1 {
					chooseAction(board, QTable)
					board = ApplyMove(board, move, player, true)
				} else {
					ABPlayer(player, true)
				}
			}
		}
	}

	return board
}

func AlphaBeta(board State, depth int, alpha int, beta int, maximizingPlayer bool, player int) (int, Move) {
	var bestMove Move

	if depth == 0 {
		return EvaluateBoard(board, 2), bestMove
	}

	if maximizingPlayer {
		maxEval := math.MinInt
		for _, move := range PossibleMoves(player) {
			newBoard := ApplyMove(board, move, player, false)
			eval, _ := AlphaBeta(newBoard, depth-1, alpha, beta, false, player)
			if eval >= maxEval {
				maxEval = eval
				bestMove = move
			}
			if alpha >= eval {
				alpha = eval
			}
			if beta <= alpha {
				break
			}
		}
		return maxEval, bestMove
	} else {
		minEval := math.MaxInt
		for _, move := range PossibleMoves(3 - player) {
			newBoard := ApplyMove(board, move, player, false)
			eval, _ := AlphaBeta(newBoard, depth-1, alpha, beta, true, player)
			if eval <= minEval {
				minEval = eval
				bestMove = move
			}
			if beta <= eval {
				beta = eval
			}
			if beta <= alpha {
				break
			}
		}
		return minEval, bestMove
	}
}

func updateQValue(state State, action Move, reward float64, newState State, QTable map[State]QValueMap) {
	currentQValue := GetQValue(state, action)
	maxFutureQ := getMaxFutureQ(newState, QTable)
	newQValue := (1-LearningRate)*currentQValue + LearningRate*(reward+DiscountFactor*maxFutureQ)

	SetQValue(state, action, newQValue)
}

func getMaxFutureQ(state State, QTable map[State]QValueMap) float64 {
	maxQ := math.Inf(-1)
	if _, exists := QTable[state]; !exists {
		return 0.0
	}

	for _, qValue := range QTable[state] {
		if qValue > maxQ {
			maxQ = qValue
		}
	}

	if maxQ == math.Inf(-1) {
		return 0.0 //new state
	}

	return maxQ
}

func chooseAction(state State, QTable map[State]QValueMap) Move {
	if rand.Float64() < ExplorationRate {
		return chooseRandomAction(state)
	} else {
		return chooseBestAction(state, QTable)
	}
}

func chooseRandomAction(state State) Move {
	PossibleMoves := PossibleMoves(1)
	return PossibleMoves[rand.Intn(len(PossibleMoves))]
}

func chooseBestAction(state State, QTable map[State]QValueMap) Move {
	bestMove := Move{}
	maxQValue := math.Inf(-1)

	for move, qValue := range QTable[state] {
		if qValue > maxQValue {
			maxQValue = qValue
			bestMove = move
		}
	}

	if maxQValue == math.Inf(-1) { //new state
		return chooseRandomAction(state)
	}

	return bestMove
}

func InitializeQTable(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonReadyTable := make(map[string]map[string]float64)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonReadyTable); err != nil {
		return err
	}

	QTable = make(map[State]QValueMap)

	for stateStr, qValuesStr := range jsonReadyTable { //from json map to types
		state, err := stringToState(stateStr)
		if err != nil {
			return err
		}
		QTable[state] = make(QValueMap)
		for moveStr, qValue := range qValuesStr {
			move, err := stringToMove(moveStr)
			if err != nil {
				return err
			}
			QTable[state][move] = qValue
		}
	}

	return nil
}

func stringToMove(s string) (Move, error) {
	var move Move
	parts := strings.Split(s, ",")

	if len(parts) != 4 {
		return move, fmt.Errorf("invalid move string: expected 4 values but got %d", len(parts))
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			value, err := strconv.Atoi(parts[i*2+j])
			if err != nil {
				return move, fmt.Errorf("invalid integer value in move string: %v", err)
			}
			move[i][j] = value
		}
	}

	return move, nil
}

func stringToState(s string) (State, error) {
	var state State
	parts := strings.Split(s, ",")

	if len(parts) != 65 { //64 + trailing empty string
		return state, fmt.Errorf("invalid state string: expected 64 values but got %d", len(parts)-1)
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			value, err := strconv.Atoi(parts[i*8+j])
			if err != nil {
				return state, fmt.Errorf("invalid integer value in state string: %v", err)
			}
			state[i][j] = value
		}
	}

	return state, nil
}

func GetQValue(state State, move Move) float64 {
	if _, ok := QTable[state]; !ok {
		QTable[state] = make(QValueMap)
	}
	if _, ok := QTable[state][move]; !ok {
		QTable[state][move] = 0.0 // default
	}
	return QTable[state][move]
}

func SetQValue(state State, move Move, value float64) {
	if _, ok := QTable[state]; !ok {
		QTable[state] = make(QValueMap)
	}
	QTable[state][move] = value
}

func evaluateReward(previousState State, currentState State, player int) float64 {
	reward := 0.0
	opponent := 3 - player

	previousScore, currentScore := EvaluateBoard(previousState, opponent), EvaluateBoard(currentState, opponent)

	scoreDif := currentScore - previousScore
	if scoreDif < 5 {
		reward += float64(scoreDif)
	} else {
		reward += 10
	}

	return reward
}

func ExportQTableToJSON(filename string) error {
	jsonReadyTable := make(map[string]map[string]float64)
	for state, qValues := range QTable {
		stateStr := stateToString(state)
		jsonReadyTable[stateStr] = make(map[string]float64)
		for move, qValue := range qValues {
			moveStr := moveToString(move)
			jsonReadyTable[stateStr][moveStr] = qValue
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonReadyTable); err != nil {
		return err
	}

	return nil
}

func playQLearningGame() {
	NewBoard()
	currentPlayer := 1
	moveCount := 0

	for moveCount < 300 {
		if len(PossibleMoves(currentPlayer)) == 0 {
			break
		}
		if currentPlayer == 2 {
			ABPlayer(currentPlayer, false)
		} else {
			currentState := board
			action := chooseAction(currentState, QTable)
			board = ApplyMove(currentState, action, currentPlayer, false)
			//PrintState()
			//time.Sleep(8 * time.Second)
			newState := board
			reward := evaluateReward(currentState, newState, 1)

			updateQValue(currentState, action, reward, newState, QTable)
		}

		currentPlayer = 3 - currentPlayer
		moveCount++
	}
}

func stateToString(state State) string {
	var builder strings.Builder
	for _, row := range state {
		for _, cell := range row {
			fmt.Fprintf(&builder, "%d,", cell)
		}
	}
	return builder.String()
}

func moveToString(move Move) string {
	return fmt.Sprintf("%d,%d,%d,%d", move[0][0], move[0][1], move[1][0], move[1][1])
}

func main() {
	QTable = make(map[State]QValueMap)

	InitializeQTable("qtable.json")

	for i := 0; i < 10; i++ {
		playQLearningGame()
	}

	err := ExportQTableToJSON("qtable.json")
	if err != nil {
		fmt.Println("Error exporting QTable:", err)
	}
}
