<script>
  /** 'accent' | 'danger' */
  export let tone = 'accent';
  export let small = false;
</script>

<button
  type="button"
  class="ac-link"
  class:ac-link--danger={tone === 'danger'}
  class:ac-link--sm={small}
  {...$$restProps}
  on:click
>
  <slot />
</button>
