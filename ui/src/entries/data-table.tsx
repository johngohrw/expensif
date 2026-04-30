import { hydrateRoot } from 'react-dom/client';
import { DataTable } from '../components/DataTable';

// Simple action registry for declarative onClick handlers.
// Extend this map at build time if you need custom client-side actions.
const actionRegistry: Record<string, (row: Record<string, unknown>) => void> = {};

function init() {
  const roots = document.querySelectorAll<HTMLDivElement>('[data-table-root]');

  roots.forEach((container) => {
    const script = container.querySelector('script[type="application/json"]');
    if (!script) {
      console.error('[island:data-table] Missing JSON props in', container);
      return;
    }

    try {
      const props = JSON.parse(script.textContent || '{}');
      hydrateRoot(container, <DataTable {...props} actions={actionRegistry} />);
    } catch (err) {
      console.error('[island:data-table] Hydration failed:', err);
    }
  });
}

init();
