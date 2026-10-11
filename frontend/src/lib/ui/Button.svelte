<script>
  /** 'regular' | 'default' | 'danger' — use 'default' once per screen (the Return action). */
  export let variant = 'regular';
  /** 'regular' (28px) | 'small' (22px, inside group rows) */
  export let size = 'regular';
  export let type = 'button';
  export let disabled = false;
  export let wide = false;
  export let block = false;
</script>

<button
  {type}
  {disabled}
  class="ac-btn"
  class:ac-btn--default={variant === 'default'}
  class:ac-btn--danger={variant === 'danger'}
  class:ac-btn--sm={size === 'small'}
  class:ac-btn--wide={wide}
  class:ac-btn--block={block}
  {...$$restProps}
  on:click
>
  <slot />
</button>
