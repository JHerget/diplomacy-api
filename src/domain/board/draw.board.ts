import * as pi from "pureimage";
import { Providence, Player, Coordinates } from "@interfaces";
import { Readable, Writable } from "stream";

const armyRadius = 8;
const supplyCenterRadius = 4;

export const drawBoard = async (
    board: Providence[],
    players: Player[],
    mapBuf: Buffer,
): Promise<Buffer> => {
    const map = await pi.decodePNGFromStream(Readable.from(mapBuf));
    const canvas = pi.make(map.width, map.height);
    const ctx = canvas.getContext("2d");
    const colors = new Map(players.map((p) => [p.name, p.color]));

    ctx.drawImage(map, 0, 0);

    for (const p of board) {
        if (p.supplyCenter && p.supplyCenter.controlledBy) {
            drawCircle(
                ctx,
                p.supplyCenter.coordinates,
                supplyCenterRadius,
                colors.get(p.supplyCenter.controlledBy)!,
            );
        }

        if (p.unit) {
            const color = colors.get(p.unit.controlledBy)!;

            if (p.unit.type == "army") {
                drawCircle(ctx, p.coordinates, armyRadius, color);
            }

            if (p.unit.type == "fleet") {
                drawTriangle(ctx, p.coordinates, color);
            }
        }
    }

    const chunks: Buffer[] = [];
    const writable = new Writable({
        write(chunk, _enc, cb) {
            chunks.push(chunk);
            cb();
        },
    });

    await pi.encodePNGToStream(canvas, writable);
    writable.end();

    return Buffer.concat(chunks);
};

const drawCircle = (
    ctx: pi.Context,
    coords: Coordinates,
    radius: number,
    color: string,
): void => {
    ctx.beginPath();
    ctx.arc(coords.x, coords.y, radius, 0, Math.PI * 2);
    ctx.fillStyle = color;
    ctx.fill();
    ctx.strokeStyle = "black";
    ctx.stroke();
};

const drawTriangle = (
    ctx: pi.Context,
    coords: Coordinates,
    color: string,
): void => {
    const topX = coords.x;
    const topY = coords.y - 5;
    const bottomY = coords.y + 5;
    const leftX = coords.x - 12;
    const rightX = coords.x + 12;

    ctx.beginPath();
    ctx.moveTo(leftX, bottomY);
    ctx.lineTo(topX, topY);
    ctx.lineTo(rightX, bottomY);
    ctx.lineTo(leftX, bottomY);
    ctx.fillStyle = color;
    ctx.fill();
    ctx.strokeStyle = "black";
    ctx.stroke();
};
