package main

import (
	"diplomacy-api/internal/lambdas/test"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(test.Handler)
}
