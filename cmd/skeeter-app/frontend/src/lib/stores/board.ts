import { writable, get } from 'svelte/store';
import type { BoardData, BoardFilter } from '../types';
import { GetBoard } from '../../../wailsjs/go/main/App';
import { filters } from './filters';
import { currentConfig } from './config';
import { activeRepoName } from './repos';
import { notifyError } from './notifications';

export const board = writable<BoardData>({ columns: [], config: null, repoName: '' });
export const loading = writable(false);

export async function refreshBoard() {
  loading.set(true);
  try {
    const f = get(filters);
    const data = await GetBoard(f);
    board.set(data);
    if (data.config) {
      currentConfig.set(data.config);
    }
    activeRepoName.set(data.repoName || '');
  } catch (e) {
    notifyError(e);
  } finally {
    loading.set(false);
  }
}
