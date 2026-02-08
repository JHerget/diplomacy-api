import type {
    APIGatewayProxyEventV2,
    APIGatewayProxyResultV2,
} from "aws-lambda";
import { GameRepository } from "@repositories";
import { S3Api } from "@apis";
import { buffer } from "node:stream/consumers";
import { Constants } from "@constants";
import { ErrorHandler, InvalidGameId } from "@errors";
import { BoardDomain } from "@domain";

const lambda = async (
    event: APIGatewayProxyEventV2,
): Promise<APIGatewayProxyResultV2> => {
    const id = event.pathParameters?.gid || "";
    const game = await GameRepository.get(id);

    if (!game) throw new InvalidGameId(id);

    const map = await S3Api.get(Constants.mapsBucket, game.map.filename);
    const mapBuf = await buffer(map);
    const boardBuf = await BoardDomain.drawBoard(
        game.board,
        game.players,
        mapBuf,
    );

    return {
        statusCode: 200,
        headers: {
            "Content-Type": "image/png",
            "Cache-Control": "no-store",
        },
        isBase64Encoded: true,
        body: boardBuf.toString("base64"),
    };
};

export const board = ErrorHandler(lambda);
