# Expensif Context

Language for the Expensif expense-tracking domain. This is a personal / small-household expense tracker, so terms are deliberately plain and user-facing.

## Language

**Expense**:
A single purchase or payment recorded by a user. It has an amount, category, description, date, currency, and optionally a payer.
_Avoid_: transaction, record, entry

**Category**:
A user-defined label that groups expenses, such as "food" or "transport".
_Avoid_: tag, bucket

**Description**:
A short free-text note explaining what the expense was for, such as "lunch with team".
_Avoid_: memo, note, detail

**Currency**:
The ISO 4217 code in which an expense amount is recorded, e.g. USD, MYR, JPY.
_Avoid_: money, denomination

**User**:
A person tracked in the system. A user can be selected as the payer for an expense.
_Avoid_: account, profile, person (when used as a generic substitute)

**Payer**:
The user who paid for an expense. Recorded as `paidBy` / `paidById` on an expense.
_Avoid_: paid by, spender

**Preferred Currency**:
The currency chosen in Preferences, used to display converted totals across expenses recorded in different currencies.
_Avoid_: base currency, home currency

**Preferences**:
The singleton settings for the application: preferred currency, default payer, and timezone.
_Avoid_: settings, config

**Exchange Rate**:
The conversion rate from the base currency (USD) to a target currency, fetched from Frankfurter and cached per day.
_Avoid_: fx rate, conversion factor

**Converted Amount**:
The expense amount expressed in the preferred currency for display. It is computed at render time, not stored.
_Avoid_: normalized amount, display amount

**Daily View**:
A page that groups expenses by date, newest first, with a subtotal per day.

**Calendar View**:
A monthly grid showing which days have expenses and how much was spent, using a heat map to encode magnitude.

**List View**:
A page showing all expenses in a single table with totals and conversion information.

**Category Suggestion**:
A recently used category offered as a quick-select pill in the add/edit form.

**Description Suggestion**:
A previously used description for a given category, offered as a quick-select pill in the add/edit form.

**Timezone**:
The IANA timezone (e.g. `Asia/Singapore`) used to compute "today", "yesterday", and calendar ranges. Dates are stored as bare `YYYY-MM-DD` strings; the timezone is applied at display time.
_Avoid_: locale, time offset

**Heat Level**:
A 1–5 quintile rank of a day's total spend relative to other days with any spend, used by the Calendar View to size and color the heat blob.
_Avoid_: intensity, score

**Island**:
A selectively hydrated React component mounted on an otherwise server-rendered HTML page. Each island has a named container and an entry script.
_Avoid_: widget, component (when the specific partial-hydration pattern matters)

## Cross-Context Notes

- **Date vs. Time**: Expensif stores only a calendar date (`YYYY-MM-DD`). The timezone is a display concern, not part of the stored domain fact.
- **Amount vs. Converted Amount**: The amount recorded on an expense is immutable relative to the expense's own currency. The converted amount is a view-model value derived from the exchange rate at display time.
- **User vs. Payer**: A user is a person in the system; a payer is the role that person plays for a specific expense.
