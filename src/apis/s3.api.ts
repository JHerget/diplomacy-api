import { S3Client, GetObjectCommand } from "@aws-sdk/client-s3";
import { Readable } from "stream";
import { Constants } from "@constants";

const client = new S3Client({ region: Constants.awsRegion });

export const S3Api = {
    get: async (bucket: string, path: string): Promise<Readable> => {
        const response = await client.send(
            new GetObjectCommand({
                Bucket: bucket,
                Key: path,
            }),
        );

        return response.Body as Readable;
    },
};
