import { orders } from "../src/lambdas/orders.lambda";

const event = {
    version: "2.0",
    routeKey: "POST /example",
    rawPath: "/example",
    rawQueryString: "",
    headers: {
        "content-type": "application/json",
    },
    requestContext: {
        accountId: "123456789012",
        apiId: "example",
        domainName: "example.com",
        domainPrefix: "example",
        http: {
            method: "POST",
            path: "/example",
            protocol: "HTTP/1.1",
            sourceIp: "127.0.0.1",
            userAgent: "jest-test",
        },
        requestId: "id",
        routeKey: "POST /example",
        stage: "$default",
        time: "12/Mar/2023:19:03:58 +0000",
        timeEpoch: 1678657438000,
    },
    isBase64Encoded: false,
    body: JSON.stringify({
        id: "testOrderId",
        phaseId: "6956432f33c5739468982b5a",
        playerName: "Italy",
        createdDate: 0,
        value: "a ven-pie",
    }),
    pathParameters: {
        gid: "69564bb933c5739468982b67",
        tid: "0001",
    },
};

(async () => {
    const response = await orders(event);
    console.log(response);
})();
