<script lang="ts">
  import { onMount } from 'svelte';
  import { board, refreshBoard } from './lib/stores/board';
  import { activeRepoName, refreshRepos } from './lib/stores/repos';
  import Board from './components/Board.svelte';
  import FilterBar from './components/FilterBar.svelte';
  import RepoSidebar from './components/RepoSidebar.svelte';
  import AddRepoDialog from './components/AddRepoDialog.svelte';
  import CreateDialog from './components/CreateDialog.svelte';
  import TaskDetail from './components/TaskDetail.svelte';
  import Toast from './components/Toast.svelte';

  let sidebarOpen = true;
  let addRepoOpen = false;
  let createOpen = false;

  onMount(async () => {
    await refreshRepos();
    await refreshBoard();
  });

  function toggleSidebar() {
    sidebarOpen = !sidebarOpen;
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.ctrlKey && e.key === 'n') {
      e.preventDefault();
      createOpen = true;
    }
  }

  function toggleTheme() {
    document.documentElement.classList.toggle('light');
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="app">
  <header class="header">
    <div class="header-left">
      <button class="icon-btn" on:click={toggleSidebar} title="Toggle sidebar">
        &#9776;
      </button>
      <h1 class="app-title">
        Skeeter
        {#if $activeRepoName}
          <span class="repo-label">— {$activeRepoName}</span>
        {/if}
      </h1>
    </div>
    <div class="header-right">
      <button class="icon-btn" on:click={toggleTheme} title="Toggle theme">
        &#9788;
      </button>
      <button class="new-task-btn" on:click={() => createOpen = true}>
        + New Task
      </button>
    </div>
  </header>

  <div class="content">
    <RepoSidebar
      {sidebarOpen}
      onAddRepo={() => addRepoOpen = true}
    />
    <div class="main">
      <FilterBar />
      <Board />
    </div>
  </div>
</div>

<TaskDetail />
<CreateDialog open={createOpen} onClose={() => createOpen = false} />
<AddRepoDialog open={addRepoOpen} onClose={() => addRepoOpen = false} />
<Toast />

<style>
  .app {
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border);
    min-height: 48px;
    -webkit-app-region: drag;
  }

  .header-left, .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
    -webkit-app-region: no-drag;
  }

  .app-title {
    font-size: 15px;
    font-weight: 700;
    color: var(--text-primary);
  }

  .repo-label {
    font-weight: 400;
    color: var(--text-secondary);
  }

  .icon-btn {
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: 18px;
    padding: 4px 8px;
    border-radius: var(--radius);
  }

  .icon-btn:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .new-task-btn {
    background: var(--accent);
    color: var(--accent-text);
    border: none;
    border-radius: var(--radius);
    padding: 6px 14px;
    font-size: 13px;
    font-weight: 500;
  }

  .new-task-btn:hover {
    background: var(--accent-hover);
  }

  .content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  .main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
</style>
