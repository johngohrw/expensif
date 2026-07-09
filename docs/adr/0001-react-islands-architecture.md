# React Islands Architecture

We render the application as server-side Go HTML templates and hydrate only small, interactive islands with React. Go owns routing, data fetching, and the initial document; React owns the few widgets that need client-side interactivity (pills, tables, mobile navigation).

We chose this over a full SPA because the app's pages are document-like (daily list, calendar, tables) and most of the UI does not need client-side state. We chose it over vanilla JS because the interactive widgets are complex enough that React's component model and testability are worth the build step.

This decision is the foundation of the current frontend architecture. Reversing it would require rebuilding routing, state management, and data loading on the client.
