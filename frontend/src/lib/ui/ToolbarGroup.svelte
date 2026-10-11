<script>
  /** Accessible name for the group, e.g. "Window actions". */
  export let label = undefined;
</script>

<!-- Related toolbar buttons share one glass capsule. -->
<div class="ac-tgroup" role="group" aria-label={label}>
  <slot />
</div>
