import { type ReactNode } from "react";

type ColumnType = "text" | "date" | "currency" | "badge" | "actions";

interface ActionConfig {
  type: "link" | "form" | "button";
  href?: string;
  action?: string;
  text?: string;
  icon?: "pencil" | "trash";
  confirm?: string;
  onClick?: string;
  variant?: "neutral" | "danger";
}

export interface Column {
  key: string;
  title: string;
  type?: ColumnType;
  width?: string;
  actions?: ActionConfig[];
  render?: (value: unknown, row: Record<string, unknown>) => ReactNode;
}

interface DataTableProps {
  columns: Column[];
  data: Record<string, unknown>[];
  variant?: "default" | "ghost";
  meta?: Record<string, unknown>;
  actions?: Record<string, (row: Record<string, unknown>) => void>;
}

const variantStyles = {
  default: {
    container:
      "bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden",
    thead: "bg-gray-50 border-b border-gray-200",
    th: "text-left px-4 py-3 font-semibold text-gray-600 whitespace-nowrap",
    tr: "border-b border-gray-100 hover:bg-gray-50",
    td: "px-4 py-3",
  },
  ghost: {
    container: "",
    thead: "bg-[#fdfdfd]",
    th: "text-left px-6 py-2 font-semibold text-gray-500 text-xs whitespace-nowrap",
    tr: "hover:bg-gray-50",
    td: "px-6 py-3",
  },
};

const Icons: Record<string, ReactNode> = {
  pencil: (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className="h-4 w-4"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
      />
    </svg>
  ),
  trash: (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className="h-4 w-4"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={2}
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
      />
    </svg>
  ),
};

const actionBase =
  "inline-flex items-center gap-1 font-medium transition focus:outline-none";
const actionVariantClasses: Record<string, string> = {
  neutral:
    "bg-gray-100 text-gray-700 hover:bg-gray-200 rounded px-3 py-1.5 text-xs",
  danger: "bg-red-50 text-red-600 hover:bg-red-100 rounded px-3 py-1.5 text-xs",
};

function humanDate(dateStr: string): string {
  const t = new Date(dateStr + "T00:00:00");
  if (isNaN(t.getTime())) return dateStr;

  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const date = new Date(t);
  date.setHours(0, 0, 0, 0);

  const diffMs = today.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return "Today";
  if (diffDays === 1) return "Yesterday";

  const rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" });
  return rtf.format(-diffDays, "day");
}

function currencySymbol(code: string): string {
  const map: Record<string, string> = {
    MYR: "RM",
    USD: "$",
    JPY: "¥",
    CNY: "¥",
    THB: "฿",
    EUR: "€",
    GBP: "£",
    SGD: "S$",
    KRW: "₩",
    AUD: "A$",
    CAD: "C$",
    INR: "₹",
    VND: "₫",
    PHP: "₱",
    IDR: "Rp",
    HKD: "HK$",
    TWD: "NT$",
  };
  return map[code] || code;
}

function formatCurrency(amount: number, symbol: string): string {
  return `${symbol}${amount.toFixed(2)}`;
}

function ActionCell({
  row,
  actions,
  registry,
}: {
  row: Record<string, unknown>;
  actions: ActionConfig[];
  registry?: Record<string, (row: Record<string, unknown>) => void>;
}) {
  const replaceParams = (str?: string) => {
    if (!str) return str;
    return str.replace(/\{(\w+)\}/g, (_match, key) => String(row[key] ?? ""));
  };

  return (
    <div className="flex gap-2 justify-end">
      {actions.map((action, i) => {
        const className = `${actionBase} ${actionVariantClasses[action.variant || "neutral"]}`;
        const content = (
          <>
            {action.icon && Icons[action.icon]}
            {action.text && <span>{replaceParams(action.text)}</span>}
          </>
        );

        if (action.type === "link" && action.href) {
          return (
            <a key={i} href={replaceParams(action.href)} className={className}>
              {content}
            </a>
          );
        }

        if (action.type === "form" && action.action) {
          return (
            <form
              key={i}
              method="POST"
              action={replaceParams(action.action)}
              className="inline"
              onSubmit={
                action.confirm
                  ? (e) => {
                      if (!confirm(replaceParams(action.confirm)!))
                        e.preventDefault();
                    }
                  : undefined
              }
            >
              <button type="submit" className={className}>
                {content}
              </button>
            </form>
          );
        }

        const handleClick =
          action.onClick && registry?.[action.onClick]
            ? () => {
                const handler = action.onClick
                  ? registry[action.onClick]
                  : undefined;
                handler?.(row);
              }
            : undefined;

        return (
          <button
            key={i}
            type="button"
            className={className}
            onClick={handleClick}
          >
            {content}
          </button>
        );
      })}
    </div>
  );
}

const defaultRenderers: Record<
  ColumnType,
  (
    value: unknown,
    row: Record<string, unknown>,
    meta?: Record<string, unknown>,
    column?: Column,
    registry?: Record<string, (row: Record<string, unknown>) => void>,
  ) => ReactNode
> = {
  text: (value) =>
    value === undefined || value === null || value === "" ? "-" : String(value),
  date: (value) => {
    if (typeof value !== "string") return "-";
    return <span title={value}>{humanDate(value)}</span>;
  },
  currency: (_value, row, meta) => {
    const converted =
      typeof row.convertedAmount === "number"
        ? row.convertedAmount
        : Number(row.amount ?? 0);
    const original = typeof row.amount === "number" ? row.amount : 0;
    const originalCurrency = String(row.currency ?? "");
    const prefCurrency = String(meta?.currency ?? "USD");
    const prefSymbol = String(
      meta?.currencySymbol ?? currencySymbol(prefCurrency),
    );

    return (
      <div>
        <div className="font-medium">
          {formatCurrency(converted, prefSymbol)}
        </div>
        {originalCurrency && originalCurrency !== prefCurrency && (
          <div className="text-xs text-gray-400">
            {formatCurrency(original, currencySymbol(originalCurrency))}
          </div>
        )}
      </div>
    );
  },
  badge: (value) => (
    <span className="inline-block bg-blue-50 text-blue-700 px-2 py-0.5 rounded text-xs font-medium">
      {String(value ?? "")}
    </span>
  ),
  actions: (_value, row, _meta, column, registry) => {
    if (!column?.actions) return null;
    return (
      <ActionCell row={row} actions={column.actions} registry={registry} />
    );
  },
};

export function DataTable({
  columns,
  data,
  variant = "default",
  meta,
  actions: registry,
}: DataTableProps) {
  const styles = variantStyles[variant];

  const renderCell = (column: Column, row: Record<string, unknown>) => {
    if (column.render) {
      return column.render(row[column.key], row);
    }
    const type = column.type || "text";
    const renderer = defaultRenderers[type];
    return renderer(row[column.key], row, meta, column, registry);
  };

  const tableContent = (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead className={styles.thead}>
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                className={styles.th}
                style={col.width ? { width: col.width } : undefined}
              >
                {col.title}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => (
            <tr key={rowIndex} className={styles.tr}>
              {columns.map((col) => (
                <td key={col.key} className={styles.td}>
                  {renderCell(col, row)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );

  if (variant === "ghost") {
    return tableContent;
  }

  return <div className={styles.container}>{tableContent}</div>;
}
