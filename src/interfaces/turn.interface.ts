export interface Turn {
    id: string;
    phaseId: string;
    orders: Order[];
    turnNumber: number;
    startDate: number;
    endDate: number;
}

export interface Order {
    id: string;
    phaseId: string;
    playerName: string;
    createdDate: number;
    value: string;
}
