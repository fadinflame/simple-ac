<script>
  import Icon from './Icon.svelte';

  export let value = 0;
  export let min = 0;
  export let max = 99;
  export let size = 'regular';
  export let id = undefined;

  const clamp = (n) => Math.min(max, Math.max(min, Number.isFinite(n) ? Math.round(n) : min));

  function step(d) {
    value = clamp(Number(value) + d);
  }
  function onKey(e) {
    if (e.key === 'ArrowUp') { e.preventDefault(); step(1); }
    if (e.key === 'ArrowDown') { e.preventDefault(); step(-1); }
  }
  function onInput(e) {
    const n = parseInt(e.currentTarget.value, 10);
    if (!Number.isNaN(n)) value = n;
  }
  function onBlur(e) {
    value = clamp(Number(value));
    e.currentTarget.value = String(value);
  }
</script>

<div class="ac-field-wrap">
  <label class="ac-field" class:ac-field--sm={size === 'small'}>
    <input {id} type="text" inputmode="numeric" {value} on:input={onInput} on:keydown={onKey} on:blur={onBlur} />
    <span class="ac-stepper">
      <button type="button" tabindex="-1" aria-label="Increment" on:click={() => step(1)}>
        <Icon name="chevron-up" size={10} stroke={2} />
      </button>
      <button type="button" tabindex="-1" aria-label="Decrement" on:click={() => step(-1)}>
        <Icon name="chevron-down" size={10} stroke={2} />
      </button>
    </span>
  </label>
</div>
