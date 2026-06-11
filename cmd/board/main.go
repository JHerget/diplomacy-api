package main

import (
	"diplomacy-api/internal/lambdas/board"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(board.Handler)
}
