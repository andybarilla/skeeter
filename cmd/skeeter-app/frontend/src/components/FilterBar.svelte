<script lang="ts">
  import { filters, clearFilters } from '../lib/stores/filters';
  import { currentConfig } from '../lib/stores/config';
  import { refreshBoard } from '../lib/stores/board';

  $: config = $currentConfig;

  // Collect unique assignees/tags from the board would need another call,
  // but we keep it simple: priority is from config, assignee/tag are free-text.

  async function onFilterChange() {
    await refreshBoard();
  }
</script>

<div class="filter-bar">
  {#if config}
    <select
      bind:value={$filters.priority}
      on:change={onFilterChange}
      class="filter-select"
    >
      <option value="">All Priorities</option>
      {#each config.priorities as p}
        <option value={p}>{p}</option>
      {/each}
    </select>
  {/if}

  <input
    type="text"
    bind:value={$filters.assignee}
    placeholder="Assignee..."
    class="filter-input"
    on:change={onFilterChange}
  />

  <input
    type="text"
    bind:value={$filters.tag}
    placeholder="Tag..."
    class="filter-input"
    on:change={onFilterChange}
  />

  {#if $filters.priority || $filters.assignee || $filters.tag}
    <button class="clear-btn" on:click={() => { clearFilters(); refreshBoard(); }}>
      Clear
    </button>
  {/if}
</div>

<style>
  .filter-bar {
    display: flex;
    gap: 8px;
    align-items: center;
    padding: 8px 12px;
    border-bottom: 1px solid var(--border);
  }

  .filter-select, .filter-input {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 5px 10px;
    font-size: 13px;
    outline: none;
  }

  .filter-select:focus, .filter-input:focus {
    border-color: var(--accent);
  }

  .filter-input {
    width: 120px;
  }

  .clear-btn {
    background: none;
    color: var(--text-muted);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 5px 12px;
    font-size: 13px;
  }

  .clear-btn:hover {
    color: var(--text-primary);
    border-color: var(--border-hover);
  }
</style>
