import { writable } from 'svelte/store';
import type { RepoEntry } from '../types';
import { GetRepos, GetActiveRepoName } from '../../../wailsjs/go/main/App';
import { notifyError } from './notifications';

export const repos = writable<RepoEntry[]>([]);
export const activeRepoName = writable('');

export async function refreshRepos() {
  try {
    const list = await GetRepos();
    repos.set(list || []);
    const name = await GetActiveRepoName();
    activeRepoName.set(name);
  } catch (e) {
    notifyError(e);
  }
}
