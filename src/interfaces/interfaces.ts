import {
    Coast,
    LocationReference,
    MoveCommand,
    SupportCommand,
    ConvoyCommand,
    ReinforceCommand,
    DisbandCommand,
    Commands,
    CommandsRegex,
} from "./command.interface";
import { Game } from "./game.interface";
import { Map } from "./map.interface";
import { Phase } from "./phase.interface";
import { Player } from "./player.interface";
import {
    ProvidenceType,
    Providence,
    SupplyCenter,
    Unit,
    Coordinates,
} from "./providence.interface";
import { Secrets } from "./secrets.interface";
import { Turn, Order } from "./turn.interface";
import { User } from "./user.interface";

export {
    type Coast,
    type LocationReference,
    type MoveCommand,
    type SupportCommand,
    type ConvoyCommand,
    type ReinforceCommand,
    type DisbandCommand,
    type Commands,
    type Game,
    type Map,
    type Phase,
    type Player,
    type ProvidenceType,
    type Providence,
    type SupplyCenter,
    type Unit,
    type Coordinates,
    type Secrets,
    type Turn,
    type Order,
    type User,
    CommandsRegex,
};
