<script>
  export let label = '';
  /** Secondary line under the label. */
  export let sub = '';
  /** Read-only value on the right. */
  export let value = '';
  export let mono = false;
</script>

<div class="ac-row">
  <div class="ac-row-text">
    <span class="ac-row-label"><slot name="label">{label}</slot></span>
    {#if sub}<span class="ac-row-sub">{sub}</span>{/if}
  </div>
  {#if value !== ''}<span class="ac-row-value" class:ac-row-value--mono={mono}>{value}</span>{/if}
  {#if $$slots.default}<div class="ac-row-aside"><slot /></div>{/if}
</div>
