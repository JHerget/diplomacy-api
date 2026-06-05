import { describe, it, expect } from "vitest";
import { orders } from "@mocks";
import { OrderDomain } from "@domain";
import { Commands } from "@interfaces";

describe("Parse player orders into individual RawCommands.", () => {
    it("Outputs a list of RawCommands given a list of Orders.", () => {
        const mockOrders = orders.slice(0, 2);
        const expectedOutput: Commands = {
            hold: [
                {
                    playerName: "Italy",
                    unitType: "army",
                    location: {
                        id: "rom",
                        coast: null,
                    },
                    destination: {
                        id: "rom",
                        coast: null,
                    },
                },
                {
                    playerName: "Italy",
                    unitType: "fleet",
                    location: {
                        id: "tun",
                        coast: null,
                    },
                    destination: {
                        id: "tun",
                        coast: null,
                    },
                },
            ],
            move: [
                {
                    playerName: "Italy",
                    unitType: "army",
                    location: {
                        id: "pie",
                        coast: null,
                    },
                    destination: {
                        id: "mar",
                        coast: null,
                    },
                },
            ],
            retreat: [
                {
                    playerName: "Germany",
                    unitType: "army",
                    location: {
                        id: "ber",
                        coast: null,
                    },
                    destination: {
                        id: "sil",
                        coast: null,
                    },
                },
            ],
            support: [
                {
                    playerName: "Italy",
                    unitType: "fleet",
                    location: {
                        id: "gol",
                        coast: null,
                    },
                    move: {
                        playerName: "Italy",
                        unitType: "army",
                        location: {
                            id: "pie",
                            coast: null,
                        },
                        destination: {
                            id: "mar",
                            coast: null,
                        },
                    },
                },
            ],
            convoy: [],
            reinforce: [
                {
                    playerName: "Germany",
                    unitType: "fleet",
                    location: {
                        id: "pru",
                        coast: null,
                    },
                },
                {
                    playerName: "Germany",
                    unitType: "army",
                    location: {
                        id: "bal",
                        coast: null,
                    },
                },
            ],
            disband: [
                {
                    playerName: "Germany",
                    unitType: "army",
                    location: {
                        id: "mun",
                        coast: null,
                    },
                },
            ],
        };
        const actualOutput = OrderDomain.getCommands(mockOrders);

        expect(actualOutput).toEqual(expectedOutput);
    });
});
