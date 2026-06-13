# Converting Diplomacy API to Go

This guide describes a practical, step-by-step migration from the current TypeScript Lambda project to Go. It does not map every TypeScript file one-for-one, but it does call out the main project boundaries and gives detailed guidance for converting `src/domain/board/draw.board.ts`.

## 1. Establish the Go project layout

Start by creating a Go module at the repository root:

```sh
go mod init diplomacy-api
```

Organize the Go code by feature instead of by technical layer. Keep each feature's handler, models, persistence, and domain behavior close together, and reserve shared packages for infrastructure that is genuinely used across features:

```text
cmd/
  board/
    main.go
  orders/
    main.go
  test/
    main.go
internal/
  board/
    handler.go
    draw_board.go
    apply_turn.go
    model.go
  orders/
    handler.go
    commands.go
    validate.go
    model.go
  games/
    model.go
    repository.go
  platform/
    aws/
      s3.go
      secrets.go
    mongo/
      mongo.go
    http/
      errors.go
      handler.go
```

Keep Lambda entry points in `cmd/<lambda>/main.go`, and keep reusable implementation under `internal/`. The `cmd` files should stay thin: initialize shared clients, call the relevant feature handler, and start `lambda.Start`.

Feature-oriented organization means:

- Board-specific behavior goes in `internal/board`, including drawing, turn application, and board request handling.
- Order parsing and validation go in `internal/orders`.
- Game persistence and the shared game document shape go in `internal/games`.
- AWS clients, Mongo connection management, and generic HTTP/Lambda error handling go in `internal/platform`.
- Avoid a broad `model` package unless a type is truly shared by most features. Prefer `board.Providence`, `orders.MoveCommand`, and `games.Game` over one large catch-all model namespace.

## 2. Convert TypeScript interfaces to Go structs

Move the files under `src/interfaces/` into the feature package that owns the behavior using those types.

Suggested ownership:

- `Game`, `Player`, `Turn`, and map metadata: `internal/games/model.go`
- `Providence`, `SupplyCenter`, `Unit`, and `Coordinates`: `internal/board/model.go`
- `LocationReference`, command structs, and command regex helpers: `internal/orders/model.go`

If a type is needed across features, import it from its owning feature package. For example, `games.Game` can include `[]board.Providence`, `[]games.Player`, and `[]games.Turn`. This keeps ownership clear without flattening everything into a generic model package.

Important model conversions:

- TypeScript `string | null` should generally become `*string`.
- TypeScript object fields that may be absent or null should become pointers.
- JSON and BSON tags are both needed for MongoDB-backed structs.
- TypeScript string unions should become named Go string types with constants.

Example:

```go
package board

type UnitType string

const (
    UnitArmy  UnitType = "army"
    UnitFleet UnitType = "fleet"
)

type Coordinates struct {
    X float64 `json:"x" bson:"x"`
    Y float64 `json:"y" bson:"y"`
}

type Providence struct {
    ID            string              `json:"id" bson:"id"`
    Name          string              `json:"name" bson:"name"`
    SupplyCenter  *SupplyCenter       `json:"supplyCenter" bson:"supplyCenter"`
    Unit          *Unit               `json:"unit" bson:"unit"`
    Coordinates   Coordinates         `json:"coordinates" bson:"coordinates"`
    Type          string              `json:"type" bson:"type"`
    Routes        []string            `json:"routes" bson:"routes"`
    CoastalRoutes map[string][]string `json:"coastalRoutes" bson:"coastalRoutes"`
}

type SupplyCenter struct {
    ControlledBy *string     `json:"controlledBy" bson:"controlledBy"`
    Coordinates  Coordinates `json:"coordinates" bson:"coordinates"`
}

type Unit struct {
    ID           string                   `json:"id" bson:"id"`
    ControlledBy string                   `json:"controlledBy" bson:"controlledBy"`
    Type         UnitType                 `json:"type" bson:"type"`
    Location     orders.LocationReference `json:"location" bson:"location"`
}
```

MongoDB documents currently appear to be queried by `_id`, so the `Game` struct should include a BSON ObjectID field if the data stores Mongo IDs:

```go
ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
```

If API responses need the string ID shape currently used by TypeScript, add a response DTO or custom JSON handling rather than leaking Mongo-specific details everywhere.

## 3. Replace dependency APIs

Replace the TypeScript dependencies with these Go equivalents:

- AWS S3: `github.com/aws/aws-sdk-go-v2/service/s3`
- AWS Secrets Manager: `github.com/aws/aws-sdk-go-v2/service/secretsmanager`
- MongoDB: `go.mongodb.org/mongo-driver/v2/mongo`
- Lambda runtime: `github.com/aws/aws-lambda-go/lambda`
- API Gateway events: `github.com/aws/aws-lambda-go/events`
- Drawing: Go standard library `image`, `image/color`, `image/draw`, `image/png`

The current TypeScript code caches the Mongo client and connection promise in `src/apis/mongodb.api.ts`. Do the same in Go with package-level variables and `sync.Once` or a guarded lazy initializer so warm Lambda invocations reuse connections.

## 4. Convert Lambda handlers

Each Lambda should become a separate binary:

- `src/lambdas/board.lambda.ts` -> `cmd/board/main.go`
- `src/lambdas/orders.lambda.ts` -> `cmd/orders/main.go`
- `src/lambdas/test.lambda.ts` -> `cmd/test/main.go`

Use `events.APIGatewayV2HTTPRequest` and `events.APIGatewayV2HTTPResponse`.

Put request handling logic in the feature package, not in a generic handlers package. For example, `cmd/board/main.go` should call `board.Handle`, and `board.Handle` can use `games.Get`, `platform/aws.GetObject`, and `board.DrawBoard`.

The board Lambda flow should remain:

1. Read `gid` from `event.PathParameters`.
2. Fetch the game from Mongo.
3. Fetch the map PNG from S3.
4. Draw the board.
5. Return a base64 PNG response with `Content-Type: image/png` and `Cache-Control: no-store`.

The response body must still be base64 encoded:

```go
return events.APIGatewayV2HTTPResponse{
    StatusCode:      200,
    Headers:         map[string]string{"Content-Type": "image/png", "Cache-Control": "no-store"},
    IsBase64Encoded: true,
    Body:            base64.StdEncoding.EncodeToString(boardPNG),
}, nil
```

## 5. Convert feature persistence

Move `src/repositories/game.repository.ts` to `internal/games/repository.go`.

Preserve these behaviors:

- Validate Mongo ObjectID strings before querying.
- Query the `diplomacy` database and `games` collection.
- Return a typed `games.Game`.
- Return a domain-specific invalid game ID error for malformed IDs or missing games, depending on the current API behavior you want to preserve.

Keep persistence functions next to the feature data they load. If the board handler needs to fetch a game, it should call `games.Get` rather than owning a duplicate game repository.

## 6. Convert order parsing and validation

Move `src/domain/order/*` into `internal/orders/`.

Keep the existing order pipeline:

1. Parse raw order strings into command structs.
2. Validate parsed commands against the board state.
3. Return grouped command slices for hold, move, retreat, support, convoy, reinforce, and disband.

Use Go `regexp` for the current command patterns. Avoid overhauling the parser during the language migration. Once the Go version matches behavior, the parser can be improved separately.

## 7. Convert turn application

Move `src/domain/board/apply-turn.board.ts` into `internal/board/apply_turn.go`.

The current TypeScript implementation uses a package-level `boardMap`. In Go, prefer a local `map[string]*node` inside `ApplyTurn` and pass it to helper functions. That avoids cross-request state in warm Lambda containers.

Suggested shape:

```go
func ApplyTurn(providences []Providence, turn games.Turn) []Providence {
    boardMap := make(map[string]*node, len(providences))
    // populate boardMap, apply commands, finalize state
}
```

This is a good place to add tests before changing behavior. Port `src/domain/board/apply-turn.test.ts` first, then convert the implementation until the tests pass.

## 8. Detailed conversion for `draw.board.ts`

The current file does four things:

1. Decode the original map PNG from bytes using `pureimage`.
2. Create a same-size canvas and draw the map onto it.
3. Draw ownership and unit markers:
   - Supply centers: filled circle, radius `4`, black stroke.
   - Armies: filled circle, radius `8`, black stroke.
   - Fleets: filled triangle, top point `y - 5`, bottom points `x +/- 12, y + 5`, black stroke.
4. Encode the modified image back to PNG bytes.

In Go, implement this in `internal/board/draw_board.go` with standard image packages. The function signature should be:

```go
func DrawBoard(board []Providence, players []games.Player, mapPNG []byte) ([]byte, error)
```

Use `png.Decode` to read the input image:

```go
src, err := png.Decode(bytes.NewReader(mapPNG))
if err != nil {
    return nil, err
}
```

Create a mutable RGBA canvas and copy the decoded PNG into it:

```go
bounds := src.Bounds()
canvas := image.NewRGBA(bounds)
draw.Draw(canvas, bounds, src, bounds.Min, draw.Src)
```

Build the player color lookup from `players`:

```go
colors := map[string]color.RGBA{}
for _, p := range players {
    c, err := parseHexColor(p.Color)
    if err != nil {
        return nil, err
    }
    colors[p.Name] = c
}
```

The TypeScript implementation assumes every `controlledBy` value has a matching player color. Go should either preserve that assumption with a clear error or add a fallback. Prefer returning an error because missing colors indicate invalid game data.

### Color parsing

If player colors are stored as CSS hex strings such as `#ff0000`, parse them into `color.RGBA`:

```go
func parseHexColor(s string) (color.RGBA, error) {
    if len(s) != 7 || s[0] != '#' {
        return color.RGBA{}, fmt.Errorf("invalid color %q", s)
    }

    r, err := strconv.ParseUint(s[1:3], 16, 8)
    if err != nil {
        return color.RGBA{}, err
    }
    g, err := strconv.ParseUint(s[3:5], 16, 8)
    if err != nil {
        return color.RGBA{}, err
    }
    b, err := strconv.ParseUint(s[5:7], 16, 8)
    if err != nil {
        return color.RGBA{}, err
    }

    return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}
```

If current data can contain named CSS colors like `"red"` or `"black"`, either normalize them in the database during migration or extend `parseHexColor` with an explicit name map.

### Drawing circles

The TypeScript version uses canvas paths with fill and black stroke. In Go, draw filled pixels inside the radius, then draw a thin black outline near the edge:

```go
const (
    armyRadius         = 8
    supplyCenterRadius = 4
)

var black = color.RGBA{A: 255}

func drawCircle(img *image.RGBA, coords Coordinates, radius int, fill color.RGBA) {
    cx := int(math.Round(coords.X))
    cy := int(math.Round(coords.Y))
    r2 := radius * radius
    inner := (radius - 1) * (radius - 1)

    for y := cy - radius; y <= cy+radius; y++ {
        for x := cx - radius; x <= cx+radius; x++ {
            dx := x - cx
            dy := y - cy
            d2 := dx*dx + dy*dy
            if d2 > r2 {
                continue
            }
            if d2 >= inner {
                img.Set(x, y, black)
            } else {
                img.Set(x, y, fill)
            }
        }
    }
}
```

`image.RGBA.Set` safely ignores points outside the bounds, so this is acceptable for markers near the image edge.

### Drawing fleets

The TypeScript fleet marker is a filled triangle:

- Top: `(x, y - 5)`
- Left bottom: `(x - 12, y + 5)`
- Right bottom: `(x + 12, y + 5)`

Use a simple point-in-triangle test for the fill and draw three black outline lines:

```go
func drawTriangle(img *image.RGBA, coords Coordinates, fill color.RGBA) {
    top := image.Point{X: int(math.Round(coords.X)), Y: int(math.Round(coords.Y)) - 5}
    left := image.Point{X: int(math.Round(coords.X)) - 12, Y: int(math.Round(coords.Y)) + 5}
    right := image.Point{X: int(math.Round(coords.X)) + 12, Y: int(math.Round(coords.Y)) + 5}

    minX := min(left.X, top.X, right.X)
    maxX := max(left.X, top.X, right.X)
    minY := min(left.Y, top.Y, right.Y)
    maxY := max(left.Y, top.Y, right.Y)

    for y := minY; y <= maxY; y++ {
        for x := minX; x <= maxX; x++ {
            if pointInTriangle(image.Point{X: x, Y: y}, left, top, right) {
                img.Set(x, y, fill)
            }
        }
    }

    drawLine(img, left, top, black)
    drawLine(img, top, right, black)
    drawLine(img, right, left, black)
}
```

Implement `pointInTriangle` using barycentric signs:

```go
func pointInTriangle(p, a, b, c image.Point) bool {
    d1 := sign(p, a, b)
    d2 := sign(p, b, c)
    d3 := sign(p, c, a)

    hasNeg := d1 < 0 || d2 < 0 || d3 < 0
    hasPos := d1 > 0 || d2 > 0 || d3 > 0
    return !(hasNeg && hasPos)
}

func sign(p1, p2, p3 image.Point) int {
    return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}
```

Implement `drawLine` with Bresenham's algorithm so the triangle outline is deterministic and dependency-free:

```go
func drawLine(img *image.RGBA, from, to image.Point, c color.RGBA) {
    x0 := from.X
    y0 := from.Y
    x1 := to.X
    y1 := to.Y

    dx := abs(x1 - x0)
    dy := -abs(y1 - y0)

    sx := -1
    if x0 < x1 {
        sx = 1
    }

    sy := -1
    if y0 < y1 {
        sy = 1
    }

    err := dx + dy
    for {
        img.Set(x0, y0, c)
        if x0 == x1 && y0 == y1 {
            break
        }

        e2 := 2 * err
        if e2 >= dy {
            err += dy
            x0 += sx
        }
        if e2 <= dx {
            err += dx
            y0 += sy
        }
    }
}

func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
```

### Full draw flow

The translated `DrawBoard` function should look like this at a high level:

```go
func DrawBoard(board []Providence, players []games.Player, mapPNG []byte) ([]byte, error) {
    src, err := png.Decode(bytes.NewReader(mapPNG))
    if err != nil {
        return nil, err
    }

    bounds := src.Bounds()
    canvas := image.NewRGBA(bounds)
    draw.Draw(canvas, bounds, src, bounds.Min, draw.Src)

    colors := make(map[string]color.RGBA, len(players))
    for _, p := range players {
        c, err := parseHexColor(p.Color)
        if err != nil {
            return nil, err
        }
        colors[p.Name] = c
    }

    for _, p := range board {
        if p.SupplyCenter != nil && p.SupplyCenter.ControlledBy != nil {
            c, ok := colors[*p.SupplyCenter.ControlledBy]
            if !ok {
                return nil, fmt.Errorf("missing color for player %q", *p.SupplyCenter.ControlledBy)
            }
            drawCircle(canvas, p.SupplyCenter.Coordinates, supplyCenterRadius, c)
        }

        if p.Unit != nil {
            c, ok := colors[p.Unit.ControlledBy]
            if !ok {
                return nil, fmt.Errorf("missing color for player %q", p.Unit.ControlledBy)
            }

            switch p.Unit.Type {
            case UnitArmy:
                drawCircle(canvas, p.Coordinates, armyRadius, c)
            case UnitFleet:
                drawTriangle(canvas, p.Coordinates, c)
            default:
                return nil, fmt.Errorf("unknown unit type %q", p.Unit.Type)
            }
        }
    }

    var out bytes.Buffer
    if err := png.Encode(&out, canvas); err != nil {
        return nil, err
    }
    return out.Bytes(), nil
}
```

### Draw board tests

Add focused tests for `DrawBoard` before wiring it into Lambda:

1. Create a small blank PNG in memory.
2. Draw one supply center, one army, and one fleet.
3. Decode the returned PNG.
4. Assert known marker pixels changed to the expected player colors or black outline.
5. Assert missing player colors return an error.
6. Assert invalid PNG input returns an error.

This gives confidence without needing golden image snapshots. Golden PNG tests can be added later if visual regressions become hard to diagnose.

## 9. Update build and packaging

Replace the current webpack and zip build with Go Lambda builds.

Example build commands:

```sh
GOOS=linux GOARCH=arm64 go build -o dist/board/bootstrap ./cmd/board
cd dist/board && zip ../board.zip bootstrap

GOOS=linux GOARCH=arm64 go build -o dist/orders/bootstrap ./cmd/orders
cd dist/orders && zip ../orders.zip bootstrap
```

Use the Lambda `provided.al2023` runtime for custom Go binaries, or use a managed Go runtime if available in your AWS account and region. With `provided.al2023`, the executable inside the zip must be named `bootstrap`.

## 10. Update Terraform

Change `terraform/lambda.tf` so each function points to its own Go zip instead of one webpack bundle.

Current TypeScript settings:

```hcl
handler = "lambdas.board"
runtime = "nodejs24.x"
```

Go custom runtime settings:

```hcl
handler       = "bootstrap"
runtime       = "provided.al2023"
architectures = ["arm64"]
```

Because Go builds separate binaries, stop using the current shared `../dist/lambdas.zip` archive for every function. Each Lambda should have its own zip containing a binary named `bootstrap` at the zip root.

One straightforward option is to create the zips in the build script and reference them directly from Terraform:

```hcl
locals {
  lambdas = {
    board = {
      zip = "../dist/board.zip"
    }
    orders = {
      zip = "../dist/orders.zip"
    }
    test = {
      zip = "../dist/test.zip"
    }
  }
}

resource "aws_lambda_function" "fn" {
  for_each = local.lambdas

  filename         = each.value.zip
  source_code_hash = filebase64sha256(each.value.zip)

  function_name = "diplomacy-api-v1-${each.key}"
  role          = aws_iam_role.lambda_role.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  memory_size   = 1024
  timeout       = 120
}
```

With this approach, remove the current shared archive block:

```hcl
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "../dist"
  output_path = "../dist/lambdas.zip"
}
```

Alternatively, keep using the `archive` provider, but define one `archive_file` data source per Lambda directory:

```hcl
data "archive_file" "lambda_zip" {
  for_each = local.lambdas

  type        = "zip"
  source_dir  = "../dist/${each.key}"
  output_path = "../dist/${each.key}.zip"
}
```

Then point each function at `data.archive_file.lambda_zip[each.key].output_path` and `output_base64sha256`.

The existing IAM role and policy attachments can mostly stay as-is. Go still needs the same CloudWatch Logs, Secrets Manager, and S3 permissions.

`terraform/gateway.tf` and `terraform/openapi.json` can also mostly stay as-is if the Lambda keys remain `board`, `orders`, and `test`. The API Gateway integrations reference `aws_lambda_function.fn["board"].arn`, `orders`, and `test`, so the route wiring does not care whether the functions are Node or Go.

## 11. Replace local scripts

Replace the current `package.json` scripts with a small `Makefile` or shell scripts.

Suggested targets:

```make
test:
	go test ./...

build:
	rm -rf dist
	mkdir -p dist/board dist/orders dist/test
	GOOS=linux GOARCH=arm64 go build -o dist/board/bootstrap ./cmd/board
	GOOS=linux GOARCH=arm64 go build -o dist/orders/bootstrap ./cmd/orders
	GOOS=linux GOARCH=arm64 go build -o dist/test/bootstrap ./cmd/test
	cd dist/board && zip ../board.zip bootstrap
	cd dist/orders && zip ../orders.zip bootstrap
	cd dist/test && zip ../test.zip bootstrap

validate:
	cd terraform && terraform validate
```

Keep `yarn test` and `yarn build` available during the migration until the Go test suite and Lambda build are complete. Remove Node tooling only after the Go deployment is verified.

## 12. Migration order

Use this order to reduce risk:

1. Add Go module and feature package skeletons.
2. Port feature-owned structs into `internal/board`, `internal/orders`, and `internal/games`.
3. Port tests for order parsing, validation, and turn application.
4. Port `internal/orders` parsing and validation.
5. Port `board.ApplyTurn`.
6. Port `board.DrawBoard` and add image-level tests.
7. Port shared AWS wrappers for Secrets Manager and S3 into `internal/platform/aws`.
8. Port Mongo connection management into `internal/platform/mongo`.
9. Port `games` persistence.
10. Port feature Lambda handlers.
11. Build one Go Lambda locally and invoke it with the existing event mocks converted to JSON fixtures.
12. Update Terraform for a single Lambda, deploy to a non-production stage if available, and verify.
13. Convert the remaining Lambdas.
14. Remove webpack, TypeScript, Yarn, and Node dependencies after the Go deployment fully replaces them.

## 13. Verification checklist

Before removing the TypeScript implementation, verify:

- `go test ./...` passes.
- Board PNG responses are valid images and visually match the TypeScript output.
- API Gateway responses still set `isBase64Encoded` for PNG bodies.
- Mongo ObjectID handling matches the current API behavior.
- S3 map filenames still come from `game.map.filename`.
- Lambda memory and timeout settings are still appropriate after the Go build.
- Terraform deploys without changing unrelated infrastructure.
