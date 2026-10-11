<script>
  export let label;
  /** id of the control, so clicking the label focuses it. */
  export let forId = undefined;
</script>

<div class="ac-row ac-row--field">
  <label class="ac-row-label" for={forId}>{label}</label>
  <slot />
</div>
