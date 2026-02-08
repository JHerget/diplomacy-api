type Unit = "army" | "fleet";
export type Coast = "nc" | "sc" | "ec" | "wc";

export interface LocationReference {
    name: string;
    coast: Coast | null;
}

export interface MoveCommand {
    playerName: string;
    unitType: Unit;
    location: LocationReference;
    destination: LocationReference;
}

export interface SupportCommand {
    playerName: string;
    unitType: Unit;
    location: LocationReference;
    move: MoveCommand;
}

export interface ConvoyCommand {
    playerName: string;
    unitType: Unit;
    location: LocationReference;
    move: MoveCommand;
}

export interface ReinforceCommand {
    playerName: string;
    unitType: Unit;
    location: LocationReference;
}

export interface DisbandCommand {
    playerName: string;
    unitType: Unit;
    location: LocationReference;
}

export interface Commands {
    hold: MoveCommand[];
    move: MoveCommand[];
    retreat: MoveCommand[];
    support: SupportCommand[];
    convoy: ConvoyCommand[];
    reinforce: ReinforceCommand[];
    disband: DisbandCommand[];
}

const UNIT = /([af])/;
const PROV = /([a-z]{3}|[a-z]{3}-[nsew]c)/;
const HOLD = /(?:hold|h)/;
const WS = /\s*/;
const DASH = /\s*-\s*/;

const hold = [/^/, UNIT, WS, PROV, DASH, HOLD, /$/];
const move = [/^/, UNIT, WS, PROV, DASH, PROV, /$/];
const retreat = [/^/, UNIT, WS, /r/, WS, PROV, DASH, PROV, /$/];
const support = [
    /^/,
    UNIT,
    WS,
    PROV,
    WS,
    /s/,
    WS,
    UNIT,
    WS,
    PROV,
    DASH,
    /(hold|h|[a-z]{3})/,
    /$/,
];
const convoy = [
    /^/,
    /f/,
    WS,
    PROV,
    WS,
    /c/,
    WS,
    /a/,
    WS,
    PROV,
    DASH,
    PROV,
    /$/,
];
const reinforce = [/^/, UNIT, WS, PROV, /$/];
const disband = [/^/, /d/, WS, UNIT, WS, PROV, /$/];

const build = (parts: RegExp[]) => RegExp(parts.map((p) => p.source).join(""));

export const CommandsRegex = {
    hold: build(hold),
    move: build(move),
    retreat: build(retreat),
    support: build(support),
    convoy: build(convoy),
    reinforce: build(reinforce),
    disband: build(disband),
};
