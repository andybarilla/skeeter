<script lang="ts">
  import { repos, activeRepoName, refreshRepos } from '../lib/stores/repos';
  import { refreshBoard } from '../lib/stores/board';
  import { notify, notifyError } from '../lib/stores/notifications';
  import { SwitchRepo, RemoveRepo } from '../../wailsjs/go/main/App';

  export let onAddRepo: () => void;
  export let sidebarOpen: boolean;

  async function switchRepo(name: string) {
    try {
      await SwitchRepo(name);
      await refreshRepos();
      await refreshBoard();
      notify('info', `Switched to ${name}`);
    } catch (e) {
      notifyError(e);
    }
  }

  async function removeRepo(name: string) {
    try {
      await RemoveRepo(name);
      await refreshRepos();
      await refreshBoard();
      notify('info', `Removed ${name}`);
    } catch (e) {
      notifyError(e);
    }
  }
</script>

{#if sidebarOpen}
  <aside class="sidebar">
    <div class="sidebar-header">
      <h3>Repos</h3>
    </div>
    <div class="repo-list">
      {#each $repos as repo (repo.name)}
        <div
          class="repo-item"
          class:active={repo.name === $activeRepoName}
          on:click={() => switchRepo(repo.name)}
          on:keydown={(e) => e.key === 'Enter' && switchRepo(repo.name)}
          tabindex="0"
          role="button"
        >
          <span class="repo-icon">{repo.remote ? '\u2601' : '\u{1F4C1}'}</span>
          <span class="repo-name">{repo.name}</span>
          <button
            class="remove-btn"
            on:click|stopPropagation={() => removeRepo(repo.name)}
            title="Remove repo"
          >
            &times;
          </button>
        </div>
      {:else}
        <div class="empty">No repos added</div>
      {/each}
    </div>
    <button class="add-btn" on:click={onAddRepo}>+ Add Repo</button>
  </aside>
{/if}

<style>
  .sidebar {
    width: 200px;
    min-width: 200px;
    background: var(--bg-secondary);
    border-right: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .sidebar-header {
    padding: 12px 14px;
    border-bottom: 1px solid var(--border);
  }

  .sidebar-header h3 {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .repo-list {
    flex: 1;
    overflow-y: auto;
    padding: 4px;
  }

  .repo-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 10px;
    border-radius: var(--radius);
    cursor: pointer;
    transition: background 0.1s;
  }

  .repo-item:hover {
    background: var(--bg-hover);
  }

  .repo-item.active {
    background: var(--accent);
    color: var(--accent-text);
  }

  .repo-icon {
    font-size: 14px;
    flex-shrink: 0;
  }

  .repo-name {
    font-size: 13px;
    font-weight: 500;
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .remove-btn {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 16px;
    padding: 0 2px;
    opacity: 0;
    transition: opacity 0.1s;
    line-height: 1;
  }

  .repo-item:hover .remove-btn {
    opacity: 1;
  }

  .remove-btn:hover {
    color: var(--error);
  }

  .empty {
    padding: 16px;
    text-align: center;
    color: var(--text-muted);
    font-size: 13px;
  }

  .add-btn {
    margin: 8px;
    padding: 8px;
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    font-size: 13px;
    font-weight: 500;
  }

  .add-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }
</style>
