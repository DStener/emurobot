
// store.js
import { writable } from 'svelte/store';

// Хранилище для функции
export const actionFunction = writable(null);

// Опционально: хранилище для данных
export const sharedData = writable(null);