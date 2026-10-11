<script>
  export let checked = false;
  export let label = undefined;
  export let id = undefined;
</script>

<input {id} type="checkbox" role="switch" class="ac-switch" aria-label={label} bind:checked on:change />
