import { writable } from 'svelte/store';
import type { Task } from '../types';

export const selectedTask = writable<Task | null>(null);
export const detailOpen = writable(false);

export function openDetail(task: Task) {
  selectedTask.set(task);
  detailOpen.set(true);
}

export function closeDetail() {
  detailOpen.set(false);
  selectedTask.set(null);
}
