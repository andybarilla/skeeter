<script lang="ts">
  import type { Task } from '../lib/types';
  import TaskCard from './TaskCard.svelte';
  import { handleDragOver, handleDragEnter, handleDragLeave, handleDrop } from '../lib/dnd';

  export let status: string;
  export let tasks: Task[];

  $: count = tasks ? tasks.length : 0;
</script>

<div
  class="column"
  on:dragover={handleDragOver}
  on:dragenter={handleDragEnter}
  on:dragleave={handleDragLeave}
  on:drop={(e) => handleDrop(e, status)}
  role="list"
  aria-label="{status} column"
>
  <div class="column-header">
    <h3 class="column-title">{status}</h3>
    <span class="count">{count}</span>
  </div>
  <div class="card-list">
    {#if tasks && tasks.length > 0}
      {#each tasks as task (task.id)}
        <TaskCard {task} />
      {/each}
    {:else}
      <div class="empty">No tasks</div>
    {/if}
  </div>
</div>

<style>
  .column {
    display: flex;
    flex-direction: column;
    min-width: 240px;
    flex: 1;
    background: var(--bg-secondary);
    border-radius: var(--radius-lg);
    border: 2px solid transparent;
    transition: border-color 0.15s;
  }

  :global(.column.drop-target) {
    border-color: var(--accent);
    border-style: dashed;
  }

  .column-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 14px 8px;
  }

  .column-title {
    font-size: 13px;
    font-weight: 600;
    text-transform: capitalize;
    color: var(--text-secondary);
  }

  .count {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
    background: var(--bg-tertiary);
    padding: 1px 8px;
    border-radius: 10px;
  }

  .card-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 4px 8px 12px;
    overflow-y: auto;
    flex: 1;
    min-height: 100px;
  }

  .empty {
    padding: 20px;
    text-align: center;
    color: var(--text-muted);
    font-size: 13px;
  }
</style>
