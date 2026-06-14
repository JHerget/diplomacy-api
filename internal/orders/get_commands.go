package orders

import (
	"diplomacy-api/internal/game"
	"strings"
)

type rawCommand struct {
	PlayerName string
	Value string
}

func GetCommands(orders []game.Order) game.Commands {
	
}

func process(orders []game.Order) []rawCommand {
	values := []rawCommand{}

	for _, o := range orders {
		commands := strings.Split(o.Value, ",")
		for i := range commands {
			commands[i] = strings.TrimSpace(commands[i])
		}


	}

	return values
}
