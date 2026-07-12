package orders

import (
	"diplomacy-api/internal/models"
	"strings"
)

type rawCommand struct {
	PlayerName string
	Value      string
}

func GetCommands(orders []models.Order) models.Commands {
	rawCommands := process(orders)
	commands := models.Commands{
		Hold:      []models.MoveCommand{},
		Move:      []models.MoveCommand{},
		Retreat:   []models.MoveCommand{},
		Support:   []models.SupportCommand{},
		Convoy:    []models.SupportCommand{},
		Reinforce: []models.AdjustCommand{},
		Disband:   []models.AdjustCommand{},
	}

	for _, raw := range rawCommands {
		result := []string{}

		result = models.CommandsRegex.Hold.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Hold = append(commands.Hold, models.MoveCommand{
				PlayerName:  raw.PlayerName,
				UnitType:    parseUnit(result[1]),
				Location:    parseLocationReference(result[2]),
				Destination: parseLocationReference(result[2]),
			})
			continue
		}

		result = models.CommandsRegex.Move.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Move = append(commands.Move, models.MoveCommand{
				PlayerName:  raw.PlayerName,
				UnitType:    parseUnit(result[1]),
				Location:    parseLocationReference(result[2]),
				Destination: parseLocationReference(result[3]),
			})
			continue
		}

		result = models.CommandsRegex.Retreat.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Retreat = append(commands.Retreat, models.MoveCommand{
				PlayerName:  raw.PlayerName,
				UnitType:    parseUnit(result[1]),
				Location:    parseLocationReference(result[2]),
				Destination: parseLocationReference(result[3]),
			})
			continue
		}

		result = models.CommandsRegex.Support.FindStringSubmatch(raw.Value)
		if result != nil {
			moveDest := parseLocationReference(result[5])

			if result[5] == "hold" || result[5] == "h" {
				moveDest = parseLocationReference(result[4])
			}

			commands.Support = append(commands.Support, models.SupportCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
				Move: models.MoveCommand{
					PlayerName:  raw.PlayerName,
					UnitType:    parseUnit(result[3]),
					Location:    parseLocationReference(result[4]),
					Destination: moveDest,
				},
			})
			continue
		}

		result = models.CommandsRegex.Convoy.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Convoy = append(commands.Convoy, models.SupportCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
				Move: models.MoveCommand{
					PlayerName:  raw.PlayerName,
					UnitType:    parseUnit(result[3]),
					Location:    parseLocationReference(result[4]),
					Destination: parseLocationReference(result[5]),
				},
			})
			continue
		}

		result = models.CommandsRegex.Reinforce.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Reinforce = append(commands.Reinforce, models.AdjustCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
			})
			continue
		}

		result = models.CommandsRegex.Disband.FindStringSubmatch(raw.Value)
		if result != nil {
			commands.Disband = append(commands.Disband, models.AdjustCommand{
				PlayerName: raw.PlayerName,
				UnitType:   parseUnit(result[1]),
				Location:   parseLocationReference(result[2]),
			})
			continue
		}
	}

	return commands
}

func process(orders []models.Order) []rawCommand {
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

func parseUnit(rawUnit string) models.UnitType {
	if strings.ToLower(rawUnit) == "a" {
		return models.UnitArmy
	}

	return models.UnitFleet
}

func parseLocationReference(rawLocation string) models.LocationReference {
	locParts := strings.Split(rawLocation, "-")
	defaultRef := models.LocationReference{
		ID:    locParts[0],
		Coast: nil,
	}

	if len(locParts) < 2 {
		return defaultRef
	}

	id := locParts[0]
	coast := models.Coast(locParts[1])

	switch coast {
	case models.NorthCoast, models.SouthCoast, models.EastCoast, models.WestCoast:
		return models.LocationReference{
			ID:    id,
			Coast: &coast,
		}
	default:
		return defaultRef
	}
}
