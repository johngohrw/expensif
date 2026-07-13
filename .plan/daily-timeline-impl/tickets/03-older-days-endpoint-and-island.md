---
type: task
blocked_by: [02]
claimed_by: claude-code-session-b9f33780
claimed_at: 2026-07-13T00:00:00Z
---

# The older-days endpoint and the scroll island

## Question

Deliver the spec's **Window and pagination** decisions: `GET /daily/older?start=&end=`
in the HTML handler layer, serving a fragment rendered by the same day-entry partial
as the first window, carrying the server-issued cursor as `data-next-start`/`-end`
(absent at the earliest expense — the terminal rule, final-window clamping included,
lives server-side). The four foot states server-rendered in the daily template and
toggled by `data-state`. The island itself: ~30 lines of vanilla TypeScript — sentinel
`IntersectionObserver`, one fetch in flight, append via `insertAdjacentHTML`, manual
Retry with the observer disconnected on error — plus its Vite entry. Zero Tailwind and
zero markup in TS.

Tests per the spec's seam 2: the fragment through the real renderer via `httptest`,
asserting only `data-next-*` and `data-state`; bad or inverted range → 400. Verify the
scroll live on a seeded scratch DB: append at the window seam (rails unbroken, month
dividers correct), terminal at the earliest expense, Retry on a killed server, and the
first-window terminal for a future-only account.

Done when: scrolling walks back to the earliest expense and stops, errors offer
Retry, and a JS-less load still shows a fully working first window.
