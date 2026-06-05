import { describe, it, expect } from "vitest";
import { Commands, MoveCommand, SupportCommand } from "@interfaces";
import { board } from "@mocks";
import { OrderDomain } from "@domain";

const boardMap = new Map(board.map((p) => [p.id, p]));

describe("Validate all command types.", () => {
    it("Returns a list of valid commands.", () => {
        const validHold: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: {
                id: "pie",
                coast: null,
            },
            destination: {
                id: "pie",
                coast: null,
            },
        };
        const wrongUnitTypeHold: MoveCommand = {
            playerName: "France",
            unitType: "fleet",
            location: {
                id: "pie",
                coast: null,
            },
            destination: {
                id: "pie",
                coast: null,
            },
        };
        const wrongDestHold: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: {
                id: "pie",
                coast: null,
            },
            destination: {
                id: "ven",
                coast: null,
            },
        };
        const validMove: MoveCommand = {
            playerName: "Italy",
            unitType: "army",
            location: {
                id: "ven",
                coast: null,
            },
            destination: {
                id: "pie",
                coast: null,
            },
        };
        const wrongPlayerMove: MoveCommand = {
            playerName: "Italy",
            unitType: "army",
            location: {
                id: "pie",
                coast: null,
            },
            destination: {
                id: "ven",
                coast: null,
            },
        };
        const unknownLocMove: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: {
                id: "abc",
                coast: null,
            },
            destination: {
                id: "ven",
                coast: null,
            },
        };
        const unknownDestMove: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: {
                id: "pie",
                coast: null,
            },
            destination: {
                id: "abc",
                coast: null,
            },
        };

        const allCommands = {
            hold: [],
            move: [],
            retreat: [],
            support: [],
            convoy: [],
            reinforce: [],
            disband: [],
        };
    });
});
