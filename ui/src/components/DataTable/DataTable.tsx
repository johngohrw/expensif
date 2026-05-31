import { type ReactNode } from "react";
import { defaultRenderers } from "./renderers";

export type ColumnType = "text" | "date" | "currency" | "badge" | "actions";

export interface ActionConfig {
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
      "bg-white rounded-none sm:rounded-xl shadow-sm border border-gray-200 overflow-hidden",
    thead: "bg-gray-50 border-b border-gray-200",
    th: "text-left px-1 first:pl-3 py-1 font-semibold text-gray-600 whitespace-nowrap",
    tr: "border-b border-gray-100 hover:bg-gray-50",
    td: "px-1 first:pl-3 py-1 text-xs",
  },
  ghost: {
    container: "",
    thead: "bg-[#fdfdfd]",
    th: "text-left px-1 first:pl-3 py-1 font-semibold text-gray-500 text-xs whitespace-nowrap",
    tr: "hover:bg-gray-50",
    td: "px-1 first:pl-3 py-1 text-xs whitespace-nowrap",
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
