import { writable } from 'svelte/store';
import type { Notification } from '../types';

let nextId = 0;

export const notifications = writable<Notification[]>([]);

export function notify(type: Notification['type'], message: string) {
  const id = nextId++;
  notifications.update(n => [...n, { id, type, message }]);
  setTimeout(() => {
    notifications.update(n => n.filter(item => item.id !== id));
  }, 4000);
}

export function notifyError(err: unknown) {
  const message = err instanceof Error ? err.message : String(err);
  notify('error', message);
}
