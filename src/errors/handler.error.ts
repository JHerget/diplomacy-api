import type {
    APIGatewayProxyEventV2,
    APIGatewayProxyResultV2,
} from "aws-lambda";
import { InvalidGameId } from "./game.error";
import { InvalidTurnId } from "./turn.error";

type LambdaFunc = (
    event: APIGatewayProxyEventV2,
) => Promise<APIGatewayProxyResultV2>;

const res = (status: number, message: string): APIGatewayProxyResultV2 => ({
    statusCode: status,
    body: message,
});

export const ErrorHandler = (lambda: LambdaFunc): LambdaFunc => {
    return async (
        event: APIGatewayProxyEventV2,
    ): Promise<APIGatewayProxyResultV2> => {
        try {
            return await lambda(event);
        } catch (err) {
            if (err instanceof InvalidGameId) return res(400, err.message);
            if (err instanceof InvalidTurnId) return res(400, err.message);

            throw err;
        }
    };
};
