import axios from 'axios';

// While developing, we point to localhost:8080
// In production (embedded), this will just be "/"
const API_URL = import.meta.env.DEV ? 'http://localhost:8080/api' : '/api';

export const api = axios.create({
    baseURL: API_URL,
    withCredentials: true, // IMPORTANT: This sends the cookies (Auth)
});

// Helper types matching our Go structs
export interface FileInfo {
    name: string;
    size: number;
    is_dir: boolean;
    mod_time: string;
    type: string;
}

export interface TrashInfo {
    originalPath: string;
    deletedAt: string;
    filename: string;
}