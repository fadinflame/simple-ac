<script>
  // Outline icons from the SimpleAC design system: 20×20 grid, 1.5 stroke, currentColor.
  export let name;
  export let size = 16;
  export let stroke = 1.5;

  const PATHS = {
    'sliders': '<path d="M3.5 6h7M14.5 6h2M3.5 14h2M9.5 14h7"/><circle cx="12.5" cy="6" r="2"/><circle cx="7.5" cy="14" r="2"/>',
    'doc': '<path d="M5.5 2.75h6l3.5 3.5v10a1 1 0 0 1-1 1h-8.5a1 1 0 0 1-1-1v-12.5a1 1 0 0 1 1-1z"/><path d="M11.5 2.75v3.5h3.5M7.25 10h5.5M7.25 13h5.5"/>',
    'lock': '<rect x="4.25" y="8.75" width="11.5" height="8.5" rx="2"/><path d="M6.75 8.75v-2.5a3.25 3.25 0 0 1 6.5 0v2.5"/>',
    'lock-open': '<rect x="4.25" y="8.75" width="11.5" height="8.5" rx="2"/><path d="M6.75 8.75v-2.5a3.25 3.25 0 0 1 6.3-1.1"/>',
    'shield': '<path d="M10 2.5l6 2.25v4.6c0 3.7-2.5 6.6-6 8.15-3.5-1.55-6-4.45-6-8.15v-4.6z"/>',
    'shield-check': '<path d="M10 2.5l6 2.25v4.6c0 3.7-2.5 6.6-6 8.15-3.5-1.55-6-4.45-6-8.15v-4.6z"/><path d="M7.4 10l1.85 1.85 3.4-3.6"/>',
    'power': '<path d="M10 3v6.5"/><path d="M6.1 5.6a6 6 0 1 0 7.8 0"/>',
    'trash': '<path d="M4 5.75h12M8 5.75v-1.5a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1v1.5M5.75 5.75l.7 10.6a1 1 0 0 0 1 .9h5.1a1 1 0 0 0 1-.9l.7-10.6M8.5 9v5M11.5 9v5"/>',
    'folder': '<path d="M2.75 6v8.75a1.5 1.5 0 0 0 1.5 1.5h11.5a1.5 1.5 0 0 0 1.5-1.5v-7a1.5 1.5 0 0 0-1.5-1.5h-6l-1.6-1.75h-3.9a1.5 1.5 0 0 0-1.5 1.5z"/>',
    'chevron-left': '<path d="M12.25 4.5L6.75 10l5.5 5.5"/>',
    'chevron-up': '<path d="M5 12.5l5-5 5 5"/>',
    'chevron-down': '<path d="M5 7.5l5 5 5-5"/>',
    'info': '<circle cx="10" cy="10" r="7.25"/><path d="M10 9v4.75"/><circle cx="10" cy="6.4" r=".4" fill="currentColor"/>',
    'terminal': '<rect x="2.75" y="3.75" width="14.5" height="12.5" rx="2"/><path d="M6 8l2.25 2L6 12M10 12.5h4"/>',
    'eye': '<path d="M1.75 10s3-5.75 8.25-5.75S18.25 10 18.25 10s-3 5.75-8.25 5.75S1.75 10 1.75 10z"/><circle cx="10" cy="10" r="2.5"/>',
    'eye-slash': '<path d="M7.6 4.6A8.6 8.6 0 0 1 10 4.25c5.25 0 8.25 5.75 8.25 5.75a14.6 14.6 0 0 1-2.2 2.9M12.1 15.4a8.6 8.6 0 0 1-2.1.35C4.75 15.75 1.75 10 1.75 10a14.8 14.8 0 0 1 3.1-3.75M8.25 8.25a2.5 2.5 0 0 0 3.5 3.5M3 3l14 14"/>',
    'alert': '<circle cx="10" cy="10" r="7.25"/><path d="M10 6.25v4.5"/><circle cx="10" cy="13.6" r=".4" fill="currentColor"/>',
    'arrow-down-circle': '<circle cx="10" cy="10" r="7.25"/><path d="M10 6.25v7M7 10.5l3 3 3-3"/>',
    'xmark': '<path d="M5.75 5.75l8.5 8.5M14.25 5.75l-8.5 8.5"/>',
    'app': '<rect x="2.75" y="3.75" width="14.5" height="12.5" rx="2"/><path d="M2.75 7.25h14.5"/>'
  };
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 20 20"
  fill="none"
  stroke="currentColor"
  stroke-width={stroke}
  stroke-linecap="round"
  stroke-linejoin="round"
  aria-hidden="true">{@html PATHS[name] ?? ''}</svg>
