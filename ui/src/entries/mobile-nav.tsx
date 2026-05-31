import { hydrateRoot } from 'react-dom/client';
import { MobileNav } from '../components/MobileNav';

const container = document.getElementById('mobile-nav-root');
if (!container) {
  console.error('[island:mobile-nav] Container #mobile-nav-root not found');
  throw new Error('mobile-nav-root not found');
}

try {
  const props = JSON.parse(container.dataset.props || '{}');
  hydrateRoot(container, <MobileNav {...props} />);
} catch (err) {
  console.error('[island:mobile-nav] Hydration failed:', err);
}
