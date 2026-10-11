<script>
  /** Text title; ignored when the `center` slot is used. */
  export let title = '';
  /** The title is the app name: hidden off macOS, where the native title bar already shows it. */
  export let appTitle = false;

  $: bare = !title && !$$slots.left && !$$slots.center && !$$slots.right;
</script>

<!-- Glass toolbar laid over the window; content scrolls under it. -->
<header class="ac-toolbar" class:ac-toolbar--bare={bare} style="--wails-draggable: drag">
  <slot name="left" />
  {#if $$slots.center}
    <div class="ac-toolbar-center"><slot name="center" /></div>
  {:else}
    <div class="ac-toolbar-title" class:ac-toolbar-title--app={appTitle}>{title}</div>
  {/if}
  <slot name="right" />
</header>
