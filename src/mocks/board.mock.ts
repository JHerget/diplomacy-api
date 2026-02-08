import { Providence } from "@interfaces";

export const board: Providence[] = [
    {
        id: "pie",
        name: "Piedmont",
        supplyCenter: null,
        unit: {
            id: "frenchUnitOne",
            controlledBy: "France",
            type: "army"
        },
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["ven", "tus", "gol"]
    },
    {
        id: "ven",
        name: "Venice",
        supplyCenter: {
            controlledBy: null,
            coordinates: {
                x: 0,
                y: 0
            }
        },
        unit: {
            id: "italianUnitOne",
            controlledBy: "Italy",
            type: "army"
        },
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["pie", "tus", "rom", "adr", "apu"]
    },
    {
        id: "tus",
        name: "Tuscany",
        supplyCenter: null,
        unit: {
            id: "italianUnitTwo",
            controlledBy: "Italy",
            type: "fleet"
        },
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["pie", "ven", "rom", "tyn", "gol"]
    },
    {
        id: "rom",
        name: "Rome",
        supplyCenter: {
            controlledBy: null,
            coordinates: {
                x: 0,
                y: 0
            }
        },
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["tus", "ven", "apu", "nap", "tyn"]
    },
    {
        id: "apu",
        name: "Apulia",
        supplyCenter: null,
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["ven", "adr", "ion", "nap", "rom"]
    },
    {
        id: "nap",
        name: "Naples",
        supplyCenter: {
            controlledBy: null,
            coordinates: {
                x: 0,
                y: 0
            }
        },
        unit: {
            id: "italianUnitThree",
            controlledBy: "Italy",
            type: "fleet"
        },
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["tyn", "rom", "apu", "ion"]
    },
    {
        id: "tun",
        name: "Tunisia",
        supplyCenter: {
            controlledBy: null,
            coordinates: {
                x: 0,
                y: 0
            }
        },
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["naf", "wes", "tyn", "ion"]
    },
    {
        id: "naf",
        name: "North Africa",
        supplyCenter: null,
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "all",
        routes: ["tun", "wes"]
    },
    {
        id: "adr",
        name: "Adriatic Sea",
        supplyCenter: null,
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "fleet",
        routes: ["ven", "apu", "ion"]
    },
    {
        id: "ion",
        name: "Ionian Sea",
        supplyCenter: null,
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "fleet",
        routes: ["adr", "nap", "apu", "tun", "tyn"]
    },
    {
        id: "tyn",
        name: "Tyrrhenian Sea",
        supplyCenter: null,
        unit: {
            id: "austriaHungaryUnitOne",
            controlledBy: "Austria-Hungary",
            type: "army"
        },
        coordinates: {
            x: 0,
            y: 0
        },
        type: "fleet",
        routes: ["tun", "wes", "gol", "tus", "rom", "nap", "ion"]
    },
    {
        id: "wes",
        name: "Western Mediterranean Sea",
        supplyCenter: null,
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "fleet",
        routes: ["naf", "tun", "gol"]
    },
    {
        id: "gol",
        name: "Gulf of Lyon",
        supplyCenter: null,
        unit: null,
        coordinates: {
            x: 0,
            y: 0
        },
        type: "fleet",
        routes: ["pie", "wes", "tyn"]
    }
]
