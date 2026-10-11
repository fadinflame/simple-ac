<script>
  /** Text title; ignored when the `center` slot is used. */
  export let title = '';
  /** Hairline under the bar — turn on when content scrolls beneath it. */
  export let separator = false;

  $: bare = !title && !$$slots.left && !$$slots.center && !$$slots.right;
</script>

<header class="ac-toolbar" class:ac-toolbar--sep={separator} class:ac-toolbar--bare={bare} style="--wails-draggable: drag">
  <slot name="left" />
  {#if $$slots.center}
    <div class="ac-toolbar-center"><slot name="center" /></div>
  {:else}
    <div class="ac-toolbar-title">{title}</div>
  {/if}
  <slot name="right" />
</header>
