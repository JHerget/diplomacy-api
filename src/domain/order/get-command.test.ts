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
                    location: "rom",
                    destination: "rom"
                },
                {
                    playerName: "Italy",
                    unitType: "fleet",
                    location: "tun",
                    destination: "tun"
                }
            ],
            move: [
                {
                    playerName: "Italy",
                    unitType: "army",
                    location: "pie",
                    destination: "mar"
                }
            ],
            retreat: [
                {
                    playerName: "Germany",
                    unitType: "army",
                    location: "ber",
                    destination: "sil"
                }
            ],
            support: [
                {
                    playerName: "Italy",
                    unitType: "fleet",
                    location: "gol",
                    move: {
                        playerName: "Italy",
                        unitType: "army",
                        location: "pie",
                        destination: "mar"
                    }
                }
            ],
            convoy: [],
            reinforce: [
                {
                    playerName: "Germany",
                    unitType: "fleet",
                    location: "pru",
                },
                {
                    playerName: "Germany",
                    unitType: "army",
                    location: "bal",
                }
            ],
            disband: [
                {
                    playerName: "Germany",
                    unitType: "army",
                    location: "mun"
                }
            ]
        };
        const actualOutput = OrderDomain.getCommands(mockOrders);

        expect(actualOutput).toEqual(expectedOutput);
    });
});
