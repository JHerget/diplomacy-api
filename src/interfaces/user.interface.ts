export interface User {
    id: string;
    firstName: string;
    lastName: string;
    username: string;
    password: string;
    salt: Uint8Array;
    createdDate: number;
    isDeleted: boolean;
}
