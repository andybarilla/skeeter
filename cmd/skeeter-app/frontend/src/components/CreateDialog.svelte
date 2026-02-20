<script lang="ts">
  import { CreateTask } from '../../wailsjs/go/main/App';
  import { refreshBoard } from '../lib/stores/board';
  import { currentConfig } from '../lib/stores/config';
  import { notify, notifyError } from '../lib/stores/notifications';

  export let open = false;
  export let onClose: () => void;

  let title = '';
  let priority = '';
  let assignee = '';
  let tagsStr = '';
  let body = '';
  let submitting = false;

  $: config = $currentConfig;

  function reset() {
    title = '';
    priority = '';
    assignee = '';
    tagsStr = '';
    body = '';
  }

  async function handleSubmit() {
    if (!title.trim()) return;
    submitting = true;
    try {
      const tags = tagsStr ? tagsStr.split(',').map(t => t.trim()).filter(Boolean) : [];
      const created = await CreateTask({ title, priority, assignee, tags, body });
      notify('success', `Created ${created.id}: ${created.title}`);
      reset();
      onClose();
      await refreshBoard();
    } catch (e) {
      notifyError(e);
    } finally {
      submitting = false;
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose();
  }
</script>

{#if open}
  <div class="overlay" on:click={onClose} on:keydown={handleKeydown} role="presentation">
    <div class="dialog" on:click|stopPropagation role="dialog" aria-modal="true">
      <h2>New Task</h2>
      <form on:submit|preventDefault={handleSubmit}>
        <div class="field">
          <label for="title">Title</label>
          <input id="title" bind:value={title} placeholder="Task title" required autofocus />
        </div>
        <div class="row">
          <div class="field">
            <label for="priority">Priority</label>
            <select id="priority" bind:value={priority}>
              <option value="">Default</option>
              {#if config}
                {#each config.priorities as p}
                  <option value={p}>{p}</option>
                {/each}
              {/if}
            </select>
          </div>
          <div class="field">
            <label for="assignee">Assignee</label>
            <input id="assignee" bind:value={assignee} placeholder="@user" />
          </div>
        </div>
        <div class="field">
          <label for="tags">Tags (comma-separated)</label>
          <input id="tags" bind:value={tagsStr} placeholder="bug, frontend" />
        </div>
        <div class="field">
          <label for="body">Description</label>
          <textarea id="body" bind:value={body} rows="4" placeholder="Markdown description..."></textarea>
        </div>
        <div class="actions">
          <button type="button" class="btn-secondary" on:click={onClose}>Cancel</button>
          <button type="submit" class="btn-primary" disabled={submitting || !title.trim()}>
            {submitting ? 'Creating...' : 'Create'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: var(--bg-overlay);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
  }

  .dialog {
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius-lg);
    padding: 24px;
    width: 480px;
    max-width: 90vw;
    max-height: 90vh;
    overflow-y: auto;
    box-shadow: var(--shadow-lg);
  }

  h2 {
    font-size: 18px;
    margin-bottom: 16px;
    color: var(--text-primary);
  }

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

  label {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
  }

  input, select, textarea {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 8px 10px;
    outline: none;
    color-scheme: dark;
  }

  :global(.light) select {
    color-scheme: light;
  }

  input:focus, select:focus, textarea:focus {
    border-color: var(--accent);
  }

  textarea {
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
