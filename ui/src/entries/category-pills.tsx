import { hydrateRoot } from 'react-dom/client';
import { CategoryPills } from '../components/CategoryPills';

const container = document.getElementById('category-pills-root');
if (!container) {
  console.error('[island:category-pills] Container #category-pills-root not found');
  throw new Error('category-pills-root not found');
}

try {
  const props = JSON.parse(container.dataset.props || '{}');
  hydrateRoot(container, <CategoryPills {...props} />);
} catch (err) {
  console.error('[island:category-pills] Hydration failed:', err);
}
