import { LocationReference } from "./command.interface";

export interface Providence {
    id: string;
    name: string;
    supplyCenter: SupplyCenter | null;
    unit: Unit | null;
    coordinates: Coordinates;
    type: "army" | "fleet" | "all";
    routes: string[];
    coastalRoutes: {
        [key: string]: string[]
    }
}

export interface SupplyCenter {
    controlledBy: string | null;
    coordinates: Coordinates;
}

export interface Unit {
    id: string;
    controlledBy: string;
    type: "army" | "fleet";
    location: LocationReference;
}

export interface Coordinates {
    x: number;
    y: number;
}
