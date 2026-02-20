<script lang="ts">
  import { notifications } from '../lib/stores/notifications';

  const iconMap: Record<string, string> = {
    success: '\u2713',
    error: '\u2717',
    info: '\u24D8',
  };
</script>

<div class="toast-container">
  {#each $notifications as note (note.id)}
    <div class="toast toast-{note.type}">
      <span class="icon">{iconMap[note.type] || ''}</span>
      <span class="message">{note.message}</span>
    </div>
  {/each}
</div>

<style>
  .toast-container {
    position: fixed;
    bottom: 16px;
    right: 16px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    gap: 8px;
    max-width: 360px;
  }

  .toast {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    border-radius: var(--radius);
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    box-shadow: var(--shadow-lg);
    animation: slideIn 0.2s ease-out;
  }

  .toast-success { border-left: 3px solid var(--success); }
  .toast-error { border-left: 3px solid var(--error); }
  .toast-info { border-left: 3px solid var(--accent); }

  .icon {
    font-size: 16px;
    flex-shrink: 0;
  }

  .toast-success .icon { color: var(--success); }
  .toast-error .icon { color: var(--error); }
  .toast-info .icon { color: var(--accent); }

  .message {
    font-size: 13px;
    color: var(--text-primary);
  }

  @keyframes slideIn {
    from { transform: translateX(100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }
</style>
