<script lang="ts">
  import { AddRepo, BrowseDirectory } from '../../wailsjs/go/main/App';
  import { refreshRepos } from '../lib/stores/repos';
  import { refreshBoard } from '../lib/stores/board';
  import { notify, notifyError } from '../lib/stores/notifications';

  export let open = false;
  export let onClose: () => void;

  let tab: 'local' | 'remote' = 'local';
  let name = '';
  let path = '';
  let remote = '';
  let dir = '';
  let submitting = false;

  function reset() {
    name = '';
    path = '';
    remote = '';
    dir = '';
  }

  async function handleSubmit() {
    submitting = true;
    try {
      if (tab === 'local') {
        await AddRepo({ name, path, remote: '', dir: '' });
      } else {
        await AddRepo({ name, path: '', remote, dir });
      }
      notify('success', `Added repo${name ? ': ' + name : ''}`);
      reset();
      onClose();
      await refreshRepos();
      await refreshBoard();
    } catch (e) {
      notifyError(e);
    } finally {
      submitting = false;
    }
  }

  async function browse() {
    try {
      const selected = await BrowseDirectory();
      if (selected) path = selected;
    } catch (e) {
      notifyError(e);
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') onClose();
  }
</script>

{#if open}
  <div class="overlay" on:click={onClose} on:keydown={handleKeydown} role="presentation">
    <div class="dialog" on:click|stopPropagation role="dialog" aria-modal="true">
      <h2>Add Repository</h2>

      <div class="tabs">
        <button class="tab" class:active={tab === 'local'} on:click={() => tab = 'local'}>
          Local Path
        </button>
        <button class="tab" class:active={tab === 'remote'} on:click={() => tab = 'remote'}>
          GitHub Remote
        </button>
      </div>

      <form on:submit|preventDefault={handleSubmit}>
        <div class="field">
          <label for="repo-name">Display Name (optional)</label>
          <input id="repo-name" bind:value={name} placeholder="Auto-detected from config" />
        </div>

        {#if tab === 'local'}
          <div class="field">
            <label for="repo-path">Path to .skeeter directory</label>
            <div class="path-row">
              <input id="repo-path" bind:value={path} placeholder="/home/user/project/.skeeter" required />
              <button type="button" class="btn-secondary" on:click={browse}>Browse</button>
            </div>
          </div>
        {:else}
          <div class="field">
            <label for="repo-remote">GitHub owner/repo</label>
            <input id="repo-remote" bind:value={remote} placeholder="owner/repo" required />
          </div>
          <div class="field">
            <label for="repo-dir">Directory in repo (optional)</label>
            <input id="repo-dir" bind:value={dir} placeholder=".skeeter (default)" />
          </div>
        {/if}

        <div class="actions">
          <button type="button" class="btn-secondary" on:click={onClose}>Cancel</button>
          <button type="submit" class="btn-primary" disabled={submitting}>
            {submitting ? 'Adding...' : 'Add'}
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
    width: 440px;
    max-width: 90vw;
    box-shadow: var(--shadow-lg);
  }

  h2 {
    font-size: 18px;
    margin-bottom: 16px;
  }

  .tabs {
    display: flex;
    gap: 4px;
    margin-bottom: 16px;
    background: var(--bg-secondary);
    border-radius: var(--radius);
    padding: 3px;
  }

  .tab {
    flex: 1;
    padding: 6px 12px;
    border: none;
    background: none;
    color: var(--text-secondary);
    font-size: 13px;
    font-weight: 500;
    border-radius: 4px;
    transition: all 0.15s;
  }

  .tab.active {
    background: var(--bg-primary);
    color: var(--text-primary);
    box-shadow: var(--shadow);
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 12px;
  }

  label {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-secondary);
  }

  input {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    padding: 8px 10px;
    outline: none;
  }

  input:focus {
    border-color: var(--accent);
  }

  .path-row {
    display: flex;
    gap: 8px;
  }

  .path-row input {
    flex: 1;
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
