# Diplomacy API

```mermaid
classDiagram
    class Game {
        id: Guid
        ownerId: Guid
        mapId: Guid
        players: User[]
        turns: Turn[]
        daysPerTurn: number
        turnStartHour: number
        timezone: string
        startDate: number
        endDate: number
        inProgress: bool
        isDeleted: bool
    }

    class User {
        id: Guid
        firstName: string
        lastName: string
        username: string
        password: string
        salt: Uint8Array
        createdDate: number
        isDeleted: bool
    }

    class Player {
        id: Guid
        userId: Guid
        powerId: Guid
        isPlaying: bool
    }

    class Turn {
        id: Guid
        phaseId: Guid
        orders: Order[]
        turnNumber: number
        startDate: number
        endDate: number
    }

    class Order {
        id: Guid
        phaseId: Guid
        playerId: Guid
        createdDate: number
        value: string
        isDeleted: bool
    }

    class Phase {
        id: Guid
        name: string
        description: string
        phaseOrder: number
    }

    class GreatPower {
        id: Guid
        name: string
        color: string
    }

    class Map {
        id: Guid
        name: string
        providences: Providence[]
        supplyCenters: SupplyCenter[]
    }

    class Providence {
        id: Guid
        name: string
        shortCode: string
        supplyCenterId: Guid | null
        unitX: number
        unitY: number
        type: army | fleet | all
        routes: Guid[]
        controlledBy: Guid | null
    }

    class SupplyCenter {
        providenceId: Guid
        unitX: number
        unitY: number
        controlledBy: Guid | null
    }
```
