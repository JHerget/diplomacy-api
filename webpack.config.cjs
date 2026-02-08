const path = require("path");

const lambda_root = "./src/lambdas";

module.exports = {
    target: "node",
    mode: "none",
    // entry: {
    //     "board.lambda": `${lambda_root}/board.lambda.ts`,
    // },
    entry: "./src/lambdas/lambdas.ts",
    output: {
        filename: "lambdas.js",
        path: path.resolve(__dirname, "dist"),
        libraryTarget: "commonjs2",
    },
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: "ts-loader",
                exclude: /node_modules/,
            },
            { test: /\.node$/, use: "node-loader" },
        ],
    },
    resolve: {
        extensions: [".ts", ".js"],
        alias: {
            "@apis": path.resolve(__dirname, "./src/apis/apis.ts"),
            "@constants": path.resolve(__dirname, "./src/constants.ts"),
            "@domain": path.resolve(__dirname, "./src/domain/domain.ts"),
            "@errors": path.resolve(__dirname, "./src/errors/errors.ts"),
            "@interfaces": path.resolve(
                __dirname,
                "./src/interfaces/interfaces.ts",
            ),
            "@lambdas": path.resolve(__dirname, "./src/lambdas/lambdas.ts"),
            "@repositories": path.resolve(
                __dirname,
                "./src/repositories/repositories.ts",
            ),
        },
    },
};
