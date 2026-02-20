import { writable } from 'svelte/store';
import type { BoardFilter } from '../types';

export const filters = writable<BoardFilter>({
  priority: '',
  assignee: '',
  tag: '',
});

export function clearFilters() {
  filters.set({ priority: '', assignee: '', tag: '' });
}
