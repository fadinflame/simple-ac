<script>
  /** Heading above the box (sentence case). */
  export let title = '';
  /** Footnote under the box. Use the `footer` slot for rich content. */
  export let footer = '';
</script>

<section class="ac-section">
  {#if title}<h3 class="ac-group-title">{title}</h3>{/if}
  <div class="ac-group"><slot /></div>
  {#if $$slots.footer}
    <div class="ac-group-footer"><slot name="footer" /></div>
  {:else if footer}
    <p class="ac-group-footer">{footer}</p>
  {/if}
</section>
