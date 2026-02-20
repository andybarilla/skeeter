<script lang="ts">
  import type { Task } from '../lib/types';
  import PriorityBadge from './PriorityBadge.svelte';
  import { handleDragStart, handleDragEnd } from '../lib/dnd';
  import { openDetail } from '../lib/stores/taskDetail';

  export let task: Task;
</script>

<div
  class="card"
  draggable="true"
  on:dragstart={(e) => handleDragStart(e, task.id)}
  on:dragend={handleDragEnd}
  on:click={() => openDetail(task)}
  on:keydown={(e) => e.key === 'Enter' && openDetail(task)}
  tabindex="0"
  role="button"
>
  <div class="card-header">
    <span class="task-id">{task.id}</span>
    <PriorityBadge priority={task.priority} />
  </div>
  <div class="title">{task.title}</div>
  <div class="card-footer">
    {#if task.assignee}
      <span class="assignee">@{task.assignee}</span>
    {/if}
    {#if task.tags && task.tags.length > 0}
      <div class="tags">
        {#each task.tags as tag}
          <span class="tag">{tag}</span>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .card {
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 10px 12px;
    cursor: grab;
    transition: box-shadow 0.15s, border-color 0.15s;
    user-select: none;
  }

  .card:hover {
    border-color: var(--border-hover);
    box-shadow: var(--shadow);
  }

  .card:active {
    cursor: grabbing;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }

  .task-id {
    font-size: 12px;
    font-weight: 700;
    color: var(--text-muted);
    font-family: 'SF Mono', 'Fira Code', monospace;
  }

  .title {
    font-size: 13px;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 6px;
    line-height: 1.4;
  }

  .card-footer {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    align-items: center;
  }

  .assignee {
    font-size: 11px;
    color: var(--accent);
    font-weight: 500;
  }

  .tags {
    display: flex;
    gap: 3px;
    flex-wrap: wrap;
  }

  .tag {
    font-size: 10px;
    padding: 1px 6px;
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    border-radius: 8px;
  }
</style>
