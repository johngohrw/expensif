import { type ComponentType } from 'react';
import { hydrateRoot } from 'react-dom/client';

export function hydrateIsland(name: string, Component: ComponentType<any>) {
  const container = document.querySelector(`[data-island="${name}"]`) as HTMLElement | null;
  if (!container) {
    console.error(`[island:${name}] Container [data-island="${name}"] not found`);
    throw new Error(`${name} island not found`);
  }

  try {
    const props = JSON.parse(container.dataset.props || '{}');
    hydrateRoot(container, <Component {...props} />);
  } catch (err) {
    console.error(`[island:${name}] Hydration failed:`, err);
  }
}
