package orders

import (
	"diplomacy-api/internal/game"
	"strings"
)

type rawCommand struct {
	PlayerName string
	Value      string
}

func GetCommands(orders []game.Order) game.Commands {
	rawCommands := process(orders)
	commands := game.Commands{
		Hold:      []game.MoveCommand{},
		Move:      []game.MoveCommand{},
		Retreat:   []game.MoveCommand{},
		Support:   []game.SupportCommand{},
		Convoy:    []game.SupportCommand{},
		Reinforce: []game.AdjustCommand{},
		Disband:   []game.AdjustCommand{},
	}

	for _, raw := range rawCommands {
		result := []string{}

		result = game.CommandsRegex.Hold.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Hold = append(commands.Hold, game.MoveCommand{
				PlayerName:  raw.PlayerName,
				UnitType:    parseUnit(result[1]),
				Location:    parseLocationReference(result[2]),
				Destination: parseLocationReference(result[2]),
			})
			continue
		}

		result = game.CommandsRegex.Move.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Move = append(commands.Move, game.MoveCommand{
				PlayerName:  raw.PlayerName,
				UnitType:    parseUnit(result[1]),
				Location:    parseLocationReference(result[2]),
				Destination: parseLocationReference(result[3]),
			})
			continue
		}

		result = game.CommandsRegex.Retreat.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Retreat = append(commands.Retreat, game.MoveCommand{
				PlayerName:  raw.PlayerName,
				UnitType:    parseUnit(result[1]),
				Location:    parseLocationReference(result[2]),
				Destination: parseLocationReference(result[3]),
			})
			continue
		}

		result = game.CommandsRegex.Support.FindStringSubmatch(raw.Value)
		if result != nil {
			moveDest := parseLocationReference(result[5])

			if result[5] == "hold" || result[5] == "h" {
				moveDest = parseLocationReference(result[4])
			}

			commands.Support = append(commands.Support, game.SupportCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
				Move: game.MoveCommand{
					PlayerName:  raw.PlayerName,
					UnitType:    parseUnit(result[3]),
					Location:    parseLocationReference(result[4]),
					Destination: moveDest,
				},
			})
			continue
		}

		result = game.CommandsRegex.Convoy.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Convoy = append(commands.Convoy, game.SupportCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
				Move: game.MoveCommand{
					PlayerName:  raw.PlayerName,
					UnitType:    parseUnit(result[3]),
					Location:    parseLocationReference(result[4]),
					Destination: parseLocationReference(result[5]),
				},
			})
			continue
		}

		result = game.CommandsRegex.Reinforce.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Reinforce = append(commands.Reinforce, game.AdjustCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
			})
			continue
		}

		result = game.CommandsRegex.Disband.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Disband = append(commands.Disband, game.AdjustCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
			})
			continue
		}
	}

	return commands
}

func process(orders []game.Order) []rawCommand {
	values := []rawCommand{}

	for _, o := range orders {
		commands := strings.Split(o.Value, ",")

		for _, c := range commands {
			values = append(values, rawCommand{
				PlayerName: o.PlayerName,
				Value:      strings.ToLower(strings.TrimSpace(c)),
			})
		}
	}

	return values
}

func parseUnit(rawUnit string) game.UnitType {
	if strings.ToLower(rawUnit) == "a" {
		return game.UnitArmy
	}

	return game.UnitFleet
}

func parseLocationReference(rawLocation string) game.LocationReference {
	locParts := strings.Split(rawLocation, "-")
	defaultRef := game.LocationReference{
		Id:    locParts[0],
		Coast: nil,
	}

	if len(locParts) < 2 {
		return defaultRef
	}

	id := locParts[0]
	coast := game.Coast(locParts[1])

	switch coast {
	case game.NorthCoast, game.SouthCoast, game.EastCoast, game.WestCoast:
		return game.LocationReference{
			Id:    id,
			Coast: &coast,
		}
	default:
		return defaultRef
	}
}
