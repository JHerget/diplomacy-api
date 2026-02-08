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
            location: "pie",
            destination: "pie"
        };
        const wrongUnitTypeHold: MoveCommand = {
            playerName: "France",
            unitType: "fleet",
            location: "pie",
            destination: "pie"
        };
        const wrongDestHold: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: "pie",
            destination: "ven"
        };
        const validMove: MoveCommand = {
            playerName: "Italy",
            unitType: "army",
            location: "ven",
            destination: "pie"
        };
        const wrongPlayerMove: MoveCommand = {
            playerName: "Italy",
            unitType: "army",
            location: "pie",
            destination: "ven"
        };
        const unknownLocMove: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: "abc",
            destination: "ven"
        };
        const unknownDestMove: MoveCommand = {
            playerName: "France",
            unitType: "army",
            location: "pie",
            destination: "abc"
        };

        const allCommands = {
            hold: [],
            move: [],
            retreat: [],
            support: [],
            convoy: [],
            reinforce: [],
            disband: []
        };
    });
});
