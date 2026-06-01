import { hydrateRoot } from 'react-dom/client';
import { MobileNav } from '../components/MobileNav';

const container = document.querySelector('[data-island="mobile-nav"]') as HTMLElement | null;
if (!container) {
  console.error('[island:mobile-nav] Container [data-island="mobile-nav"] not found');
  throw new Error('mobile-nav island not found');
}

try {
  const props = JSON.parse(container.dataset.props || '{}');
  hydrateRoot(container, <MobileNav {...props} />);
} catch (err) {
  console.error('[island:mobile-nav] Hydration failed:', err);
}
