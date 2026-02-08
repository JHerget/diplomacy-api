import { Player } from "./player.interface";
import { Turn } from "./turn.interface";
import { Providence } from "./providence.interface";

export interface Game {
    id: string;
    ownerId: string;
    map: {
        id: string;
        filename: string;
    };
    mapId: string;
    board: Providence[];
    players: Player[];
    turns: Turn[];
    daysPerTurn: number;
    turnStartHour: number;
    timezone: number;
    startDate: number;
    endDate: number;
    inProgress: boolean;
    isDeleted: boolean;
}
