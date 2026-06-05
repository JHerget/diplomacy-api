import {
    Providence,
    Turn,
    MoveCommand,
    SupportCommand,
    ConvoyCommand,
    ReinforceCommand,
    DisbandCommand,
} from "@interfaces";
import { OrderDomain } from "@domain";

interface Attack {
    unit: string;
    location: string;
    support: string[];
}

interface Node {
    providence: Providence;
    attacks: Attack[];
    unresolvedCommand: MoveCommand | null;
}

const boardMap = new Map<string, Node>();

export const applyTurn = (board: Providence[], turn: Turn): Providence[] => {
    const orders = turn.orders.filter((o) => o.phaseId === turn.phaseId);
    const allCommands = OrderDomain.getCommands(orders);
    const validCommands = OrderDomain.validateCommands(allCommands, board);

    boardMap.clear();
    board.forEach((p) => {
        boardMap.set(p.id, {
            providence: p,
            attacks: [],
            unresolvedCommand: null,
        });
    });

    validCommands.move.forEach((m) => applyMove(m));
    validCommands.retreat.forEach((r) => applyRetreat(r));
    validCommands.support.forEach((s) => applySupport(s));
    validCommands.convoy.forEach((c) => applyConvoy(c));

    if (turn.turnNumber % 2 == 0) {
        validCommands.reinforce.forEach((r) => applyReinforce(r));
        validCommands.disband.forEach((d) => applyDisband(d));
    }

    return finalizeState();
};

const applyMove = (move: MoveCommand): void => {
    const loc = boardMap.get(move.location.id)!;
    const dest = boardMap.get(move.destination.id)!;

    loc.unresolvedCommand = move;
    dest.attacks.push({
        unit: move.unitType,
        location: move.location.id,
        support: [],
    });
};

const applyRetreat = (retreat: MoveCommand): void => {};

const applySupport = (support: SupportCommand): void => {
    const dest = boardMap.get(support.move.destination.id)!;
    const attack = dest.attacks.find(
        (u) => u.location === support.move.location.id,
    );

    if (!attack) return;

    attack.support.push(support.location.id);
};

const applyConvoy = (convoy: ConvoyCommand): void => {};
const applyReinforce = (reinforce: ReinforceCommand): void => {};
const applyDisband = (disband: DisbandCommand): void => {};

const finalizeState = (): Providence[] => {
    const board = [];

    for (const node of boardMap.values()) {
        if (!node.unresolvedCommand) {
            board.push(resolve(node));
            continue;
        }

        const stack = [node];

        while (stack.length) {
            const top = stack.pop()!;
            const dest = boardMap.get(top.unresolvedCommand!.destination.id)!;

            if (dest.unresolvedCommand) {
                stack.push(top, dest);
                continue;
            }

            resolve(dest);
            board.push(resolve(top));
        }
    }

    return board;
};

const resolve = (node: Node): Providence => {
    const minStrength = node.providence.unit ? 1 : 0;
    const maxStrength = Math.max(
        ...node.attacks.map((a) => a.support.length + 1),
    );
    const strongestAttacks = node.attacks.filter(
        (a) => a.support.length + 1 >= maxStrength,
    );

    node.unresolvedCommand = null;
    node.attacks = [];

    if (strongestAttacks.length == 1 && maxStrength > minStrength) {
        const attack = strongestAttacks.pop();
        const loc = boardMap.get(attack!.location)!;

        if (node.providence.id !== loc.providence.id) {
            node.providence.unit = loc.providence.unit;
            loc.providence.unit = null;
        }
    }

    return node.providence;
};
