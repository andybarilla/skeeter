<script lang="ts">
  import { selectedTask, detailOpen, closeDetail } from '../lib/stores/taskDetail';
  import { currentConfig } from '../lib/stores/config';
  import { refreshBoard } from '../lib/stores/board';
  import { notify, notifyError } from '../lib/stores/notifications';
  import { UpdateTask, GetTask, EnhanceTask } from '../../wailsjs/go/main/App';
  import PriorityBadge from './PriorityBadge.svelte';

  let editing = false;
  let title = '';
  let status = '';
  let priority = '';
  let assignee = '';
  let tagsStr = '';
  let body = '';
  let saving = false;
  let enhancing = false;

  $: config = $currentConfig;
  $: task = $selectedTask;
  $: if (task && !editing) {
    title = task.title;
    status = task.status;
    priority = task.priority;
    assignee = task.assignee || '';
    tagsStr = (task.tags || []).join(', ');
    body = task.body || '';
  }

  function startEdit() {
    editing = true;
  }

  function cancelEdit() {
    editing = false;
    if (task) {
      title = task.title;
      status = task.status;
      priority = task.priority;
      assignee = task.assignee || '';
      tagsStr = (task.tags || []).join(', ');
      body = task.body || '';
    }
  }

  async function saveEdit() {
    if (!task) return;
    saving = true;
    try {
      const tags = tagsStr ? tagsStr.split(',').map(t => t.trim()).filter(Boolean) : [];
      await UpdateTask({ id: task.id, title, status, priority, assignee, tags, body });
      const updated = await GetTask(task.id);
      selectedTask.set(updated);
      editing = false;
      notify('success', `Updated ${task.id}`);
      await refreshBoard();
    } catch (e) {
      notifyError(e);
    } finally {
      saving = false;
    }
  }

  async function handleEnhance() {
    if (!task) return;
    enhancing = true;
    try {
      await EnhanceTask(task.id);
      const updated = await GetTask(task.id);
      selectedTask.set(updated);
      body = updated.body || '';
      notify('success', `Enhanced ${task.id}`);
      await refreshBoard();
    } catch (e) {
      notifyError(e);
    } finally {
      enhancing = false;
    }
  }

  function handleClose() {
    editing = false;
    closeDetail();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') handleClose();
  }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if $detailOpen && task}
  <div class="overlay" on:click={handleClose} role="presentation">
    <div class="panel" on:click|stopPropagation role="dialog" aria-modal="true">
      <div class="panel-header">
        <span class="task-id">{task.id}</span>
        <button class="close-btn" on:click={handleClose}>&times;</button>
      </div>

      {#if editing}
        <div class="edit-form">
          <div class="field">
            <label>Title</label>
            <input bind:value={title} />
          </div>
          <div class="row">
            <div class="field">
              <label>Status</label>
              <select bind:value={status}>
                {#if config}
                  {#each config.statuses as s}
                    <option value={s}>{s}</option>
                  {/each}
                {/if}
              </select>
            </div>
            <div class="field">
              <label>Priority</label>
              <select bind:value={priority}>
                {#if config}
                  {#each config.priorities as p}
                    <option value={p}>{p}</option>
                  {/each}
                {/if}
              </select>
            </div>
          </div>
          <div class="field">
            <label>Assignee</label>
            <input bind:value={assignee} placeholder="@user" />
          </div>
          <div class="field">
            <label>Tags (comma-separated)</label>
            <input bind:value={tagsStr} placeholder="bug, frontend" />
          </div>
          <div class="field">
            <label>Description</label>
            <textarea bind:value={body} rows="8"></textarea>
          </div>
          <div class="actions">
            <button class="btn-secondary" on:click={cancelEdit}>Cancel</button>
            <button class="btn-primary" on:click={saveEdit} disabled={saving}>
              {saving ? 'Saving...' : 'Save'}
            </button>
          </div>
        </div>
      {:else}
        <div class="view">
          <h2>{task.title}</h2>
          <div class="meta">
            <div class="meta-row">
              <span class="label">Status</span>
              <span class="value">{task.status}</span>
            </div>
            <div class="meta-row">
              <span class="label">Priority</span>
              <PriorityBadge priority={task.priority} />
            </div>
            {#if task.assignee}
              <div class="meta-row">
                <span class="label">Assignee</span>
                <span class="value assignee">@{task.assignee}</span>
              </div>
            {/if}
            {#if task.tags && task.tags.length > 0}
              <div class="meta-row">
                <span class="label">Tags</span>
                <span class="value">{task.tags.join(', ')}</span>
              </div>
            {/if}
            <div class="meta-row">
              <span class="label">Created</span>
              <span class="value">{task.created}</span>
            </div>
            <div class="meta-row">
              <span class="label">Updated</span>
              <span class="value">{task.updated}</span>
            </div>
          </div>
          {#if task.body}
            <div class="body">
              <pre>{task.body}</pre>
            </div>
          {/if}
          <div class="actions">
            <button class="btn-secondary" on:click={handleEnhance} disabled={enhancing}>
              {enhancing ? 'Enhancing...' : 'Enhance'}
            </button>
            <button class="btn-primary" on:click={startEdit}>Edit</button>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: var(--bg-overlay);
    z-index: 50;
    display: flex;
    justify-content: flex-end;
  }

  .panel {
    width: 480px;
    max-width: 90vw;
    height: 100%;
    background: var(--bg-primary);
    border-left: 1px solid var(--border);
    overflow-y: auto;
    padding: 20px;
    animation: slideIn 0.2s ease-out;
  }

  @keyframes slideIn {
    from { transform: translateX(100%); }
    to { transform: translateX(0); }
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .task-id {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-weight: 700;
    font-size: 14px;
    color: var(--text-muted);
  }

  .close-btn {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 24px;
    padding: 0 4px;
    line-height: 1;
  }

  .close-btn:hover {
    color: var(--text-primary);
  }

  h2 {
    font-size: 18px;
    margin-bottom: 16px;
  }

  .meta {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }

  .meta-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .label {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
    min-width: 70px;
  }

  .value {
    font-size: 13px;
    color: var(--text-primary);
  }

  .assignee {
    color: var(--accent);
  }

  .body {
    margin-bottom: 16px;
  }

  .body pre {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
    color: var(--text-secondary);
    white-space: pre-wrap;
    word-wrap: break-word;
    line-height: 1.6;
  }

  /* Edit form styles */
  .field {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 12px;
    flex: 1;
  }

  .row {
    display: flex;
    gap: 12px;
  }

  .edit-form label {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
  }

  .edit-form input, .edit-form select, .edit-form textarea {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 8px 10px;
    outline: none;
    color-scheme: dark;
  }

  :global(.light) .edit-form select {
    color-scheme: light;
  }

  .edit-form input:focus, .edit-form select:focus, .edit-form textarea:focus {
    border-color: var(--accent);
  }

  .edit-form textarea {
    resize: vertical;
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
  }

  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 16px;
  }

  .btn-primary, .btn-secondary {
    padding: 8px 16px;
    border-radius: var(--radius);
    font-weight: 500;
    font-size: 13px;
    border: none;
  }

  .btn-primary {
    background: var(--accent);
    color: var(--accent-text);
  }

  .btn-primary:hover:not(:disabled) {
    background: var(--accent-hover);
  }

  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border);
  }

  .btn-secondary:hover {
    background: var(--bg-hover);
  }
</style>
