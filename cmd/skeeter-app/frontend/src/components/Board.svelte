<script lang="ts">
  import { board, loading } from '../lib/stores/board';
  import Column from './Column.svelte';

  $: columns = $board.columns || [];
</script>

<div class="board" class:loading={$loading}>
  {#if columns.length > 0}
    {#each columns as col (col.status)}
      <Column status={col.status} tasks={col.tasks || []} />
    {/each}
  {:else}
    <div class="empty-board">
      <p>No board data</p>
      <p class="hint">Add a repo from the sidebar to get started</p>
    </div>
  {/if}
</div>

<style>
  .board {
    display: flex;
    gap: 10px;
    padding: 12px;
    height: 100%;
    overflow-x: auto;
    transition: opacity 0.2s;
  }

  .board.loading {
    opacity: 0.6;
    pointer-events: none;
  }

  .empty-board {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    color: var(--text-muted);
    gap: 8px;
  }

  .empty-board p {
    font-size: 16px;
  }

  .hint {
    font-size: 13px !important;
    color: var(--text-muted);
  }
</style>
