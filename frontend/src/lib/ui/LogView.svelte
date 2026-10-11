<script>
  import { beforeUpdate, afterUpdate, onMount } from 'svelte';

  /** Raw log text, one entry per line. A leading "HH:MM:SS" is shown as a timestamp. */
  export let text = '';
  export let emptyText = 'No log entries';

  let el;
  let stick = true;

  const TIME = /^\[?(\d{1,2}:\d{2}:\d{2}(?:\s?[AP]M)?)\]?\s+(.*)$/;
  const ERR = /\b(error|failed|failure)\b/i;

  $: lines = text
    ? text.split('\n').filter((l) => l.trim() !== '').map((l) => {
        const m = l.match(TIME);
        return { time: m ? m[1] : '', msg: m ? m[2] : l, err: ERR.test(l) };
      })
    : [];

  // Follow new output only if the user is already at the bottom.
  beforeUpdate(() => {
    if (el) stick = el.scrollHeight - el.scrollTop - el.clientHeight < 8;
  });
  afterUpdate(() => {
    if (el && stick) el.scrollTop = el.scrollHeight;
  });
  onMount(() => {
    if (el) el.scrollTop = el.scrollHeight;
  });
</script>

<div class="ac-log" role="log" bind:this={el}>
  {#each lines as line}
    <div class="ac-log-line" class:ac-log-line--err={line.err}>
      {#if line.time}<span class="ac-log-time">{line.time}</span>{/if}<span class="ac-log-msg">{line.msg}</span>
    </div>
  {:else}
    <div class="ac-log-empty">{emptyText}</div>
  {/each}
</div>
