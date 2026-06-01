import { hydrateRoot } from 'react-dom/client';
import { DescriptionPills } from '../components/DescriptionPills';

const container = document.querySelector('[data-island="description-pills"]') as HTMLElement | null;
if (!container) {
  console.error('[island:description-pills] Container [data-island="description-pills"] not found');
  throw new Error('description-pills island not found');
}

try {
  const props = JSON.parse(container.dataset.props || '{}');
  hydrateRoot(container, <DescriptionPills {...props} />);
} catch (err) {
  console.error('[island:description-pills] Hydration failed:', err);
}
