import type {
    APIGatewayProxyEventV2,
    APIGatewayProxyResultV2,
} from "aws-lambda";
import { ErrorHandler, InvalidGameId, InvalidTurnId } from "@errors";
import { GameRepository } from "@repositories";
import { BoardDomain } from "@domain";
import { Constants } from "@constants";
import { S3Api } from "@apis";
import { buffer } from "node:stream/consumers";

const lambda = async (
    event: APIGatewayProxyEventV2,
): Promise<APIGatewayProxyResultV2> => {
    const gameId = event.pathParameters?.gid || "";
    const turnId = event.pathParameters?.tid || "";
    const order = event.body ? JSON.parse(event.body) : {};

    const game = await GameRepository.get(gameId);
    if (!game) throw new InvalidGameId(gameId);

    const turn = game.turns.find((t) => t.id === turnId);
    if (!turn) throw new InvalidTurnId(turnId);

    turn.orders = [order];

    const newBoard = BoardDomain.applyTurn(game.board, turn);
    const map = await S3Api.get(Constants.mapsBucket, game.map.filename);
    const mapBuf = await buffer(map);
    const boardBuf = await BoardDomain.drawBoard(
        newBoard,
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

export const orders = ErrorHandler(lambda);
