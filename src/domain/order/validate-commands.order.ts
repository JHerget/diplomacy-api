import {
    Commands,
    MoveCommand,
    SupportCommand,
    ConvoyCommand,
    ReinforceCommand,
    DisbandCommand,
    Providence,
    LocationReference
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
    const loc = boardMap.get(hold.location.name);

    if (!loc) return false;
    if (hold.location.coast && !loc.coastalRoutes[hold.location.coast]) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== hold.unitType) return false;
    if (loc.unit.controlledBy !== hold.playerName) return false;
    if (hold.location != hold.destination) return false;

    return true;
};

const isValidMove = (
    move: MoveCommand,
    validatePlayer: boolean = true,
): boolean => {
    const loc = boardMap.get(move.location.name);
    const dest = boardMap.get(move.destination.name);

    if (!loc) return false;
    if (!dest) return false;
    if (move.location.coast && !loc.coastalRoutes[move.location.coast]) return false;
    if (move.destination.coast && !dest.coastalRoutes[move.destination.coast]) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== move.unitType) return false;
    if (validatePlayer && loc.unit.controlledBy !== move.playerName)
        return false;
    // if (!loc.routes.includes(dest.id)) return false;

    if (dest.type !== move.unitType && dest.type !== "all") return false;

    return true;
};

const isValidRetreat = (retreat: MoveCommand): boolean => {
    const loc = boardMap.get(retreat.location);
    const dest = boardMap.get(retreat.destination);

    if (!loc) return false;
    if (!dest) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== retreat.unitType) return false;
    if (loc.unit.controlledBy !== retreat.playerName) return false;
    if (!loc.routes.includes(dest.id)) return false;

    if (dest.type !== retreat.unitType && dest.type !== "all") return false;

    return true;
};

const isValidSupport = (support: SupportCommand): boolean => {
    const loc = boardMap.get(support.location);
    const dest = boardMap.get(support.move.destination);

    if (!isValidMove(support.move, false)) return false;

    if (!loc) return false;
    if (!dest) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== support.unitType) return false;
    if (loc.unit.controlledBy !== support.playerName) return false;
    if (!loc.routes.includes(dest.id)) return false;

    if (dest.type !== support.unitType && dest.type !== "all") return false;

    return true;
};

const isValidConvoy = (convoy: ConvoyCommand): boolean => {
    const loc = boardMap.get(convoy.location);
    const dest = boardMap.get(convoy.move.destination);

    if (!isValidMove(convoy.move, false)) return false;

    if (!loc) return false;
    if (!dest) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== convoy.unitType) return false;
    if (loc.unit.controlledBy !== convoy.playerName) return false;

    return true;
};

const isValidReinforce = (reinforce: ReinforceCommand): boolean => {
    const loc = boardMap.get(reinforce.location);

    if (!loc) return false;

    if (loc.unit) return false;
    if (!loc.supplyCenter) return false;
    if (loc.supplyCenter.controlledBy !== reinforce.playerName) return false;
    if (loc.type !== reinforce.unitType && loc.type !== "all") return false;

    return true;
};

const isValidDisband = (disband: DisbandCommand): boolean => {
    const loc = boardMap.get(disband.location);

    if (!loc) return false;

    if (!loc.unit) return false;
    if (loc.unit.type !== disband.unitType) return false;
    if (loc.unit.controlledBy !== disband.playerName) return false;

    return true;
};

const isValidLocation = (location: LocationReference): boolean => {
    const loc = boardMap.get(location.name);

    if (!loc) return false;
    if (location.coast && !loc.coastalRoutes[location.coast]) return false;

    return true;
}
