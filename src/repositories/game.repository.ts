import { MongoDbApi } from "@apis";
import { InvalidGameId } from "@errors";
import { Game } from "@interfaces";
import { ObjectId } from "mongodb";

export const GameRepository = {
    get: async (id: string): Promise<Game | null> => {
        const db = await MongoDbApi.getDb();
        const games = db.collection<Game>("games");

        if (!ObjectId.isValid(id)) {
            throw new InvalidGameId(id);
        }

        return await games.findOne<Game>({ _id: new ObjectId(id) });
    },
    getAll: async (): Promise<Game[]> => {
        const db = await MongoDbApi.getDb();
        const games = db.collection<Game>("games");

        return await games.find().toArray();
    },
};
