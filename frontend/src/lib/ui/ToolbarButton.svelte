<script>
  import Icon from './Icon.svelte';
  export let icon;
  /** Shown as the tooltip and used as the accessible name. */
  export let label;
  export let disabled = false;
  /** A lone action: a 34px glass circle. Inside a ToolbarGroup leave it false. */
  export let round = false;
</script>

<button
  type="button"
  class="ac-tbtn"
  class:ac-tbtn--round={round}
  title={label}
  aria-label={label}
  {disabled}
  on:click
>
  <Icon
    name={icon}
    size={round ? 16 : 18}
    stroke={round && icon.startsWith('chevron') ? 1.8 : 1.5}
  />
</button>
