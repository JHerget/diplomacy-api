export class InvalidTurnId extends Error {
    constructor(id: string) {
        super(`A turn with id '${id}' was not found.`);
        this.name = "InvalidTurnId";
    }
}
