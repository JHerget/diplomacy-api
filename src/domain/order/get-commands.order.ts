import { Order, Commands, CommandsRegex, Coast, LocationReference } from "@interfaces";

interface RawCommand {
    playerName: string;
    value: string;
}

export const getCommands = (orders: Order[]): Commands => {
    const rawCommands = process(orders);
    const commands = {
        hold: [],
        move: [],
        retreat: [],
        support: [],
        convoy: [],
        reinforce: [],
        disband: [],
    } as Commands;

    for (const raw of rawCommands) {
        let result;

        result = CommandsRegex.hold.exec(raw.value);
        if (result) {
            const [unit, location] = result.slice(1);

            if (unit && location) {
                commands.hold.push({
                    playerName: raw.playerName,
                    unitType: parseUnit(unit),
                    location: parseLocationReference(location),
                    destination: parseLocationReference(location),
                });
                continue;
            }
        }

        result = CommandsRegex.move.exec(raw.value);
        if (result) {
            const [unit, location, destination] = result.slice(1);

            if (unit && location && destination) {
                commands.move.push({
                    playerName: raw.playerName,
                    unitType: parseUnit(unit),
                    location: parseLocationReference(location),
                    destination: parseLocationReference(destination),
                });
                continue;
            }
        }

        result = CommandsRegex.retreat.exec(raw.value);
        if (result) {
            const [unit, location, destination] = result.slice(1);

            if (unit && location && destination) {
                commands.retreat.push({
                    playerName: raw.playerName,
                    unitType: parseUnit(unit),
                    location: parseLocationReference(location),
                    destination: parseLocationReference(destination),
                });
                continue;
            }
        }

        result = CommandsRegex.support.exec(raw.value);
        if (result) {
            const [unit, location, moveUnit, moveLocation, moveDestination] =
                result.slice(1);

            if (
                unit &&
                location &&
                moveUnit &&
                moveLocation &&
                moveDestination
            ) {
                const isHold = moveDestination === "hold" || moveDestination === "h";

                commands.support.push({
                    playerName: raw.playerName,
                    unitType: parseUnit(unit),
                    location: parseLocationReference(location),
                    move: {
                        playerName: raw.playerName,
                        unitType: parseUnit(moveUnit),
                        location: parseLocationReference(moveLocation),
                        destination:
                            isHold
                                ? parseLocationReference(moveLocation)
                                : parseLocationReference(moveDestination),
                    },
                });
                continue;
            }
        }

        result = CommandsRegex.convoy.exec(raw.value);
        if (result) {
            const [location, moveLocation, moveDestination] = result.slice(1);

            if (location && moveLocation && moveDestination) {
                commands.convoy.push({
                    playerName: raw.playerName,
                    unitType: "fleet",
                    location: parseLocationReference(location),
                    move: {
                        playerName: raw.playerName,
                        unitType: "army",
                        location: parseLocationReference(moveLocation),
                        destination: parseLocationReference(moveDestination),
                    },
                });
                continue;
            }
        }

        result = CommandsRegex.reinforce.exec(raw.value);
        if (result) {
            const [unit, location] = result.slice(1);

            if (unit && location) {
                commands.reinforce.push({
                    playerName: raw.playerName,
                    unitType: parseUnit(unit),
                    location: parseLocationReference(location),
                });
                continue;
            }
        }

        result = CommandsRegex.disband.exec(raw.value);
        if (result) {
            const [unit, location] = result.slice(1);

            if (unit && location) {
                commands.disband.push({
                    playerName: raw.playerName,
                    unitType: parseUnit(unit),
                    location: parseLocationReference(location),
                });
                continue;
            }
        }
    }

    return commands;
};

const process = (orders: Order[]): RawCommand[] =>
    orders
        .map((o) => {
            const commands = o.value.split(",").map((c) => c.trim());
            return commands.map((c) => ({
                playerName: o.playerName,
                value: c.toLowerCase(),
            }));
        })
        .flat();

type Unit = "army" | "fleet";
const parseUnit = (rawUnit: string): Unit =>
    rawUnit.toLowerCase() == "a" ? "army" : "fleet";

const parseLocationReference = (rawLocation: string): LocationReference => {
    const locParts = rawLocation.split("-");
    return {
        name: locParts[0]!,
        coast: locParts[1] as Coast || null
    }
}
