<script>
  /** Full path. Shown with the home folder as ~ and shortened in the middle; full path in the tooltip. */
  export let path = '';
  export let placeholder = 'No location chosen';
  export let error = false;

  function pretty(p) {
    if (!p) return '';
    let s = p.replace(/^\/Users\/[^/]+/, '~').replace(/^\/home\/[^/]+/, '~');
    if (s.length <= 34) return s;
    const sep = s.includes('\\') ? '\\' : '/';
    const parts = s.split(sep);
    if (parts.length <= 3) return s;
    return [parts[0], '…', ...parts.slice(-2)].join(sep);
  }
</script>

<div class="ac-path">
  <label class="ac-field ac-field--readonly" class:ac-field--error={error} title={path}>
    <input type="text" readonly value={pretty(path)} {placeholder} tabindex="-1" />
  </label>
  <slot />
</div>
