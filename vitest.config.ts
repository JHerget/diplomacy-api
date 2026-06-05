import { defineConfig } from "vitest/config";

export default defineConfig({
    resolve: {
        alias: {
            "@domain": "/src/domain/domain",
            "@interfaces": "/src/interfaces/interfaces",
            "@mocks": "/src/mocks/mocks",
        },
    },
    test: {
        environment: "node",
        globals: true,
        include: ["src/**/*.test.ts"],
        coverage: {
            provider: "v8",
            reportsDirectory: "./coverage",
            include: ["src/**/*.ts"],
        },
    },
});
