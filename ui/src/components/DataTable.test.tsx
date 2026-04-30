import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { DataTable } from './DataTable';

describe('DataTable', () => {
  it('renders headers and rows for default variant', () => {
    const columns = [
      { key: 'name', title: 'Name', type: 'text' as const },
      { key: 'age', title: 'Age', type: 'text' as const },
    ];
    const data = [
      { name: 'Alice', age: 30 },
      { name: 'Bob', age: 25 },
    ];

    render(<DataTable columns={columns} data={data} variant="default" />);

    expect(screen.getByText('Name')).toBeInTheDocument();
    expect(screen.getByText('Age')).toBeInTheDocument();
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
  });

  it('renders badge type correctly', () => {
    const columns = [
      { key: 'status', title: 'Status', type: 'badge' as const },
    ];
    const data = [{ status: 'Active' }];

    render(<DataTable columns={columns} data={data} />);

    const badge = screen.getByText('Active');
    expect(badge).toBeInTheDocument();
    expect(badge.tagName).toBe('SPAN');
  });

  it('renders currency with meta', () => {
    const columns = [
      { key: 'amount', title: 'Amount', type: 'currency' as const },
    ];
    const data = [
      { amount: 100, currency: 'USD', convertedAmount: 420 },
    ];
    const meta = { currency: 'MYR', currencySymbol: 'RM' };

    render(<DataTable columns={columns} data={data} meta={meta} />);

    expect(screen.getByText('RM420.00')).toBeInTheDocument();
    expect(screen.getByText('$100.00')).toBeInTheDocument();
  });

  it('renders ghost variant without outer card wrapper', () => {
    const columns = [
      { key: 'name', title: 'Name', type: 'text' as const },
    ];
    const data = [{ name: 'Test' }];

    const { container } = render(<DataTable columns={columns} data={data} variant="ghost" />);

    // Ghost variant should not have the bg-white rounded-xl container
    expect(container.querySelector('.bg-white.rounded-xl')).not.toBeInTheDocument();
    expect(screen.getByText('Test')).toBeInTheDocument();
  });

  it('renders action buttons', () => {
    const columns = [
      {
        key: 'actions',
        title: '',
        type: 'actions' as const,
        actions: [
          { type: 'link' as const, href: '/edit/{id}', text: 'Edit', variant: 'neutral' as const },
        ],
      },
    ];
    const data = [{ id: 1 }];

    render(<DataTable columns={columns} data={data} />);

    const link = screen.getByText('Edit');
    expect(link).toBeInTheDocument();
    expect(link.closest('a')).toHaveAttribute('href', '/edit/1');
  });

  it('renders custom renderer when provided', () => {
    const columns = [
      {
        key: 'name',
        title: 'Name',
        render: (_value: unknown, row: Record<string, unknown>) => (
          <strong>{String(row.name)}</strong>
        ),
      },
    ];
    const data = [{ name: 'Custom' }];

    render(<DataTable columns={columns} data={data} />);

    const el = screen.getByText('Custom');
    expect(el.tagName).toBe('STRONG');
  });

  it('formats dates with humanDate', () => {
    const today = new Date().toISOString().slice(0, 10);
    const columns = [
      { key: 'date', title: 'Date', type: 'date' as const },
    ];
    const data = [{ date: today }];

    render(<DataTable columns={columns} data={data} />);

    expect(screen.getByText('Today')).toBeInTheDocument();
  });
});
