import { Providence } from "./providence.interface";
import { Player } from "./player.interface";

export interface Map {
    id: string;
    filename: string;
    name: string;
    players: Player[];
    providences: Providence[];
}
