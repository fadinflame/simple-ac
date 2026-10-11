<script>
  import { onMount } from 'svelte';
  import Icon from './Icon.svelte';

  export let value = '';
  export let type = 'text';
  export let placeholder = '';
  /** 'regular' (28px) | 'small' (22px) */
  export let size = 'regular';
  export let mono = false;
  export let readonly = false;
  export let autofocus = false;
  /** Error: a string is shown under the field; `true` only turns the border red. */
  export let error = '';
  export let id = undefined;
  /** Adds a show/hide button for secrets (type="password"). */
  export let revealable = false;

  let input;
  let revealed = false;
  $: inputType = revealable && revealed ? 'text' : type;
  onMount(() => {
    if (autofocus) input.focus();
  });

  // `type` is dynamic (text/password), so bind:value is not allowed — sync manually.
  function onInput(e) {
    value = e.currentTarget.value;
  }

  export function focus() {
    input?.focus();
  }
</script>

<div class="ac-field-wrap">
  <label
    class="ac-field"
    class:ac-field--sm={size === 'small'}
    class:ac-field--mono={mono}
    class:ac-field--readonly={readonly}
    class:ac-field--error={!!error}
  >
    <input
      bind:this={input}
      {id}
      type={inputType}
      {placeholder}
      {readonly}
      {value}
      spellcheck="false"
      autocomplete="off"
      autocapitalize="off"
      aria-invalid={error ? 'true' : undefined}
      {...$$restProps}
      on:input={onInput}
      on:input
      on:keydown
      on:blur
    />
    {#if revealable}
      <button
        type="button"
        class="ac-field-btn"
        tabindex="-1"
        title={revealed ? 'Hide' : 'Show'}
        aria-label={revealed ? 'Hide' : 'Show'}
        on:click={() => (revealed = !revealed)}
      >
        <Icon name={revealed ? 'eye-slash' : 'eye'} size={16} />
      </button>
    {/if}
    <slot name="aside" />
  </label>
  {#if error && typeof error === 'string'}
    <div class="ac-field-msg ac-field-msg--error" role="alert">{error}</div>
  {/if}
</div>
