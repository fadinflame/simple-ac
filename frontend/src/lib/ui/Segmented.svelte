<script>
  import Icon from './Icon.svelte';

  /** [{ value, label, icon? }] */
  export let items = [];
  export let value;
  export let small = false;
  export let label = undefined;

  function onKey(e) {
    const i = items.findIndex((x) => x.value === value);
    if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
      e.preventDefault();
      const d = e.key === 'ArrowRight' ? 1 : -1;
      value = items[(i + d + items.length) % items.length].value;
    }
  }
</script>

<div class="ac-seg" class:ac-seg--sm={small} role="group" aria-label={label} on:keydown={onKey}>
  {#each items as item (item.value)}
    <button type="button" aria-pressed={item.value === value} on:click={() => (value = item.value)}>
      {#if item.icon}<Icon name={item.icon} size={14} />{/if}
      {item.label}
    </button>
  {/each}
</div>
