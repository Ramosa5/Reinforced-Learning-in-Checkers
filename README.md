# Reinforced-Learning-in-Checkers
Q-learning algorithm implemented for a checkers game in go language

## Checkers bot
The script provides a q-learning checkers bot as well as a q-table for the best checkers moves which can be used in other checkers bot projects.
The q-table is updated after every game the bot plays whether it is against another bot (in this case an alpha-beta algorithm opponent) or user.

## How it works
The value of a move is calculated using Bellman equation and value of Q is updated.

```go
func updateQValue(state State, action Move, reward float64, newState State, QTable map[State]QValueMap) {
  currentQValue := GetQValue(state, action)
	maxFutureQ := getMaxFutureQ(newState, QTable)
	newQValue := (1-LearningRate)*currentQValue + LearningRate*(reward+DiscountFactor*maxFutureQ)

	SetQValue(state, action, newQValue)
}
```

A move is chosen based on ϵ-greedy algorithm allowing for a better move exploration.

```go
func chooseAction(state State, QTable map[State]QValueMap) Move {
	if rand.Float64() < ExplorationRate {
		return chooseRandomAction(state)
	} else {
		return chooseBestAction(state, QTable)
	}
}

```

The bot plays several (the number can be set in code) games against alpha-beta algorithm and the results are saved to a .json file.

```go
	InitializeQTable("qtable.json")

	for i := 0; i < 10; i++ {
		playQLearningGame()
	}

	err := ExportQTableToJSON("qtable.json")
	if err != nil {
		fmt.Println("Error exporting QTable:", err)
	}
```
The pieces are represented by numbers on an 8x8 board.

```go
	board = State{
		{0, 2, 0, 2, 0, 2, 0, 2},
		{2, 0, 2, 0, 2, 0, 2, 0},
		{0, 2, 0, 2, 0, 2, 0, 2},
		{3, 0, 3, 0, 3, 0, 3, 0},
		{0, 3, 0, 3, 0, 3, 0, 3},
		{1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 1, 0, 1, 0, 1, 0}}
```

## Usage
In order to use the program you need to launch the code on your computer and simulate several games needed for the q-table to update.
A sample representation of the table is attached to this project. Note that it is just an example and will not give satisfying results as it is an output of a few games.

```json
  "0,2,0,2,0,2,0,2,2,0,2,0,2,0,2,0,0,2,0,2,0,3,0,3,3,0,3,0,3,0,3,0,0,3,0,1,0,3,0,3,1,0,1,0,3,0,3,0,0,1,0,2,0,2,0,1,1,0,1,0,1,0,1,0,": {
    "5,0,4,1": 0,
    "7,6,5,4": -0.1
  },
```
A key represented by a string of numbers represents current state of the board and the positions inside the object are moves made by the bot in this state with their outcome values.
The higher value the better move was made in that state.
