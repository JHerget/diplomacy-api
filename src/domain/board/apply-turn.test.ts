import { describe, it, expect } from "vitest";

describe("Apply all commands for the phase to the board.", () => {
    it("A unit still stands off another attack even if it is dislodged.", () => {});
    it("A unit still cuts support to an attack even if it is dislodged.", () => {});
});

describe("Apply MoveCommands to the board.", () => {
    it("Adds an Attack to the destination providence.", () => {});
    it("The providence's unresolvedCommand matches the MoveCommand.", () => {});
});

describe("Apply RetreatCommands to the board.", () => {
    it("Moves the unit to the destination.", () => {});
    it("Disbands the unit if retreating to the same destination as another unit.", () => {});
    it("Cancels the retreat if the destination ist occupied by a unit.", () => {});
});

describe("Apply SupportCommands to the board.", () => {
    it("Adds support to the destination if there is an attack there.", () => {});
    it("Adds an attack with support if the destination is holding.", () => {});
    it("Does not add suport when the location is being attacked.", () => {});
});

describe("Apply ConvoyCommands to the board.", () => {
    it("Adds convoy to the destination if there is an attack there.", () => {});
    it("Does not add convoy when the location is being attacked.", () => {});
});

describe("Apply ReinforceCommands to the board.", () => {
    it("Adds a unit to one of the player's starting providences that they still control.", () => {});
    it("Does not add a unit to a location that isn't one of the player's starting providences.", () => {});
    it("Does not add a unit to a location that the player does not controll.", () => {});
    it("Does not add a unit that has a type that conflicts with the providence's type.", () => {});
});

describe("Apply DisbandCommands to the board.", () => {
    it("Removes a unit from a location.", () => {});
});

describe("Finalize the state of the board after applying commands.", () => {
    it("Returns a board that represents the state after appling all commands.", () => {});
    it("Does not allow 2 units to swap places.", () => {});
    it("Allows 3 or more units to swap places.", () => {});
    it("Handles long chains of providences that depend on another providence.", () => {});
});

describe("Resolve the state of an individual providence.", () => {
    it("Leaves a providence unchanged if there are no attacks.", () => {});
    it("Dislodges the unit if there is an overwhelming attack.", () => {});
    it("Does not dislodge the unit if there is an attack equal in power to the defense.", () => {});
    it("Does not dislodge the unit if there are multiple attacks with equal power.", () => {});
});
