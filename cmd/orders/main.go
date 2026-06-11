package main

import (
	"diplomacy-api/internal/lambdas/orders"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(orders.Handler)
}
