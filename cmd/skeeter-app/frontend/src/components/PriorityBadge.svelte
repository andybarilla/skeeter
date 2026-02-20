<script lang="ts">
  export let priority: string;

  const colorMap: Record<string, string> = {
    critical: 'var(--priority-critical)',
    high: 'var(--priority-high)',
    medium: 'var(--priority-medium)',
    low: 'var(--priority-low)',
  };

  $: color = colorMap[priority] || 'var(--text-muted)';
</script>

<span class="badge" style="background: {color}">
  {priority}
</span>

<style>
  .badge {
    display: inline-block;
    padding: 1px 8px;
    border-radius: 10px;
    font-size: 11px;
    font-weight: 600;
    color: #fff;
    text-transform: lowercase;
    line-height: 1.6;
  }
</style>
