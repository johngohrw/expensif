import type { ReactNode } from "react";
import type { Column, ColumnType, ActionConfig } from "./DataTable";

const actionBase =
  "inline-flex items-center gap-1 font-medium transition focus:outline-none";
const actionVariantClasses: Record<string, string> = {
  neutral:
    "bg-gray-100 text-gray-700 hover:bg-gray-200 rounded px-2 py-1 text-xs",
  danger: "bg-red-50 text-red-600 hover:bg-red-100 rounded px-2 py-1 text-xs",
};

export const Icons: Record<string, ReactNode> = {
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

export function humanDate(dateStr: string): string {
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

export function currencySymbol(code: string): string {
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

export function formatCurrency(amount: number, symbol: string): string {
  return `${symbol}${amount.toFixed(2)}`;
}

export function ActionCell({
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

export const defaultRenderers: Record<
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
