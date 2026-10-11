<script>
  import Icon from './Icon.svelte';
  export let icon;
  /** Shown as the tooltip and used as the accessible name. */
  export let label;
  export let disabled = false;
</script>

<button type="button" class="ac-tbtn" title={label} aria-label={label} {disabled} on:click>
  <Icon name={icon} size={18} />
</button>
