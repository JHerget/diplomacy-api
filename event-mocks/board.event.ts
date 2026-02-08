import { board } from "../src/lambdas/board.lambda";

(async () => {
    const response = await board();
    console.log(response);
})();
