import {
    Commands,
    MoveCommand,
    SupportCommand,
    ConvoyCommand,
    ReinforceCommand,
    DisbandCommand,
    Providence,
    LocationReference,
} from "@interfaces";

const boardMap = new Map<string, Providence>();

export const validateCommands = (
    commands: Commands,
    board: Providence[],
): Commands => {
    boardMap.clear();
    board.forEach((p) => boardMap.set(p.id, p));

    return {
        hold: commands.hold.filter((h) => isValidHold(h)),
        move: commands.move.filter((m) => isValidMove(m)),
        retreat: commands.retreat.filter((r) => isValidRetreat(r)),
        support: commands.support.filter((s) => isValidSupport(s)),
        convoy: commands.convoy.filter((c) => isValidConvoy(c)),
        reinforce: commands.reinforce.filter((r) => isValidReinforce(r)),
        disband: commands.disband.filter((d) => isValidDisband(d)),
    };
};

const isValidHold = (hold: MoveCommand): boolean => {
    const loc = boardMap.get(hold.location.id)!;

    if (!isValidLocation(hold.location)) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== hold.unitType) return false;
    if (loc.unit.controlledBy !== hold.playerName) return false;
    if (hold.location.id !== hold.destination.id) return false;
    if (hold.location.coast !== hold.destination.coast) return false;

    return true;
};

const isValidMove = (
    move: MoveCommand,
    validatePlayer: boolean = true,
): boolean => {
    const loc = boardMap.get(move.location.id)!;
    const dest = boardMap.get(move.destination.id)!;

    if (!isValidLocation(move.location)) return false;
    if (!isValidLocation(move.destination)) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== move.unitType) return false;
    if (validatePlayer && loc.unit.controlledBy !== move.playerName)
        return false;

    if (!isValidRoute(loc, dest, move.destination)) return false;

    return true;
};

const isValidRetreat = (retreat: MoveCommand): boolean => {
    const loc = boardMap.get(retreat.location.id)!;
    const dest = boardMap.get(retreat.destination.id)!;

    if (!isValidLocation(retreat.location)) return false;
    if (!isValidLocation(retreat.destination)) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== retreat.unitType) return false;
    if (loc.unit.controlledBy !== retreat.playerName) return false;

    if (!isValidRoute(loc, dest, retreat.destination)) return false;

    return true;
};

const isValidSupport = (support: SupportCommand): boolean => {
    const loc = boardMap.get(support.location.id)!;
    const dest = boardMap.get(support.move.destination.id)!;

    if (!isValidMove(support.move, false)) return false;
    if (!isValidLocation(support.location)) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== support.unitType) return false;
    if (loc.unit.controlledBy !== support.playerName) return false;

    if (!isValidRoute(loc, dest, support.move.destination)) return false;

    return true;
};

const isValidConvoy = (convoy: ConvoyCommand): boolean => {
    const loc = boardMap.get(convoy.location.id)!;
    const dest = boardMap.get(convoy.move.destination.id)!;

    if (!isValidMove(convoy.move, false)) return false;
    if (!isValidLocation(convoy.location)) return false;

    if (!loc) return false;
    if (!dest) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== convoy.unitType) return false;
    if (loc.unit.controlledBy !== convoy.playerName) return false;

    return true;
};

const isValidReinforce = (reinforce: ReinforceCommand): boolean => {
    const loc = boardMap.get(reinforce.location.id)!;

    if (!isValidLocation(reinforce.location)) return false;

    if (loc.unit) return false;
    if (!loc.supplyCenter) return false;
    if (loc.supplyCenter.controlledBy !== reinforce.playerName) return false;
    if (reinforce.unitType === "fleet" && loc.type === "inland") return false;

    return true;
};

const isValidDisband = (disband: DisbandCommand): boolean => {
    const loc = boardMap.get(disband.location.id)!;

    if (!isValidLocation(disband.location)) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== disband.unitType) return false;
    if (loc.unit.controlledBy !== disband.playerName) return false;

    return true;
};

const isValidLocation = (ref: LocationReference): boolean => {
    const loc = boardMap.get(ref.id);

    if (!loc) return false;
    if (!ref.coast && loc.coastalRoutes) return false;
    if (ref.coast && !loc.coastalRoutes) return false;
    if (ref.coast && loc.coastalRoutes && !loc.coastalRoutes[ref.coast])
        return false;

    return true;
};

const isValidRoute = (
    location: Providence,
    destination: Providence,
    destRef: LocationReference,
): boolean => {
    if (!location.unit) return false;

    if (location.type == "ocean" && destination.type == "inland") return false;
    if (location.type == "inland" && destination.type == "ocean") return false;
    if (destination.type == "inland" && location.unit.type == "fleet")
        return false;
    if (destination.type == "ocean" && location.unit.type == "army")
        return false;

    if (destRef.coast) {
        if (!destination.coastalRoutes)
            return destination.routes.includes(location.id);

        const routes = destination.coastalRoutes[destRef.coast];
        if (!routes) return false;

        return routes.includes(location.id);
    }
    if (destination.coastalRoutes) return false;

    return destination.routes.includes(location.id);
};
