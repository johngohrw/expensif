import { hydrateRoot } from 'react-dom/client';
import { CategoryPills } from '../components/CategoryPills';

const container = document.querySelector('[data-island="category-pills"]') as HTMLElement | null;
if (!container) {
  console.error('[island:category-pills] Container [data-island="category-pills"] not found');
  throw new Error('category-pills island not found');
}

try {
  const props = JSON.parse(container.dataset.props || '{}');
  hydrateRoot(container, <CategoryPills {...props} />);
} catch (err) {
  console.error('[island:category-pills] Hydration failed:', err);
}
