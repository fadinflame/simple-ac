<script>
  import Icon from './Icon.svelte';

  /** 'disconnected' | 'connecting' | 'disconnecting' | 'connected' */
  export let state = 'disconnected';
  /** Line under the title (server name, or an error). */
  export let subtitle = '';
  /** Render the subtitle in the danger colour. */
  export let error = false;
  /** "hh:mm:ss" — shown only while connected; the row keeps its height otherwise. */
  export let timer = '';

  const TITLES = {
    connected: 'Connected',
    connecting: 'Connecting…',
    disconnecting: 'Disconnecting…',
    disconnected: 'Not Connected',
  };
  $: busy = state === 'connecting' || state === 'disconnecting';
</script>

<div class="ac-status" class:ac-status--on={state === 'connected'} class:ac-status--busy={busy} aria-live="polite">
  <div class="ac-status-glyph">
    <Icon name={state === 'connected' ? 'shield-check' : 'shield'} size={30} stroke={1.4} />
  </div>
  <div class="ac-status-title">{TITLES[state] ?? TITLES.disconnected}</div>
  <div class="ac-status-sub" style={error ? 'color: var(--danger)' : ''}>{subtitle}</div>
  <div class="ac-status-timer">{state === 'connected' ? timer : ''}</div>
</div>
