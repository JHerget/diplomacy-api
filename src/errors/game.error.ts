export class InvalidGameId extends Error {
    constructor(id: string) {
        super(`A game with id '${id}' was not found.`);
        this.name = "InvalidGameId";
    }
}
