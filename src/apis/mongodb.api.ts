import { MongoClient, Db } from "mongodb";
import { SecretsManagerApi } from "@apis";

type OptionalMongoClient = MongoClient | null;
type OptionalConnectionPromise = Promise<MongoClient> | null;

let client: OptionalMongoClient = null;
let connPromise: OptionalConnectionPromise = null;

export const MongoDbApi = {
    getDb: async (): Promise<Db> => {
        if (!client) {
            const secrets = await SecretsManagerApi.getCredentials();
            client = new MongoClient(secrets["mongodb:connection-string"]);
        }

        if (!connPromise) {
            connPromise = client.connect();
        }

        const conn = await connPromise;
        return conn.db("diplomacy");
    },
};
