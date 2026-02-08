import {
    SecretsManagerClient,
    GetSecretValueCommand,
} from "@aws-sdk/client-secrets-manager";
import { Constants } from "@constants";
import { Secrets } from "@interfaces";

const client = new SecretsManagerClient({ region: Constants.awsRegion });

export const SecretsManagerApi = {
    get: async (secretName: string): Promise<string | undefined> => {
        const response = await client.send(
            new GetSecretValueCommand({ SecretId: secretName }),
        );
        return response.SecretString;
    },
    getJson: async <T>(secretName: string): Promise<T | undefined> => {
        const secret = await SecretsManagerApi.get(secretName);
        return secret ? (JSON.parse(secret) as T) : undefined;
    },
    getCredentials: async (): Promise<Secrets> => {
        const secrets = await SecretsManagerApi.getJson<Secrets>(
            "diplomacy-credentials",
        );

        if (secrets) {
            return secrets;
        }

        return {
            "mongodb:connection-string": "",
        } as Secrets;
    },
};
