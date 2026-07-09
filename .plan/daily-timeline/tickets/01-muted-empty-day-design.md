# Muted empty-day card design

Type: prototype
Status: open
Blocked by: none

## Question

What does a day with no expenses look like, and what is its add affordance?

`templates/daily.html:44` already holds a first draft — a full card with an italic
grey line and a primary "Add one" button. It was written for the `?date=` single-day
path, where it is the only thing on the page. It has never been seen in a stack of
twenty.

Make the cheap concrete artifact and react to it. The tension to resolve: a day card
today is a bordered white panel with a header, a `data-table`, and a footer total. An
empty day has none of those. If it keeps the same chrome it is visually loud —
twenty of them drown the days that matter. If it sheds the chrome it stops reading as
a day at all.

Decide, against something you can actually look at:

- Full card, thin row, or bare date line?
- Is the add affordance a button, a hover-revealed `+`, or the whole row being a
  click target to `/expenses/new?date=`?
- Does the affordance survive on mobile, where hover does not exist?
- Does a non-empty day keep its own footer "+ Add expense" button
  (`daily.html:41`), or does that move/go?

Prototype at least two treatments in the template and view them against a seeded run
of empty and non-empty days. The output is the chosen markup, linked as an asset.

Constraint from the map: no islands. This must be static markup and CSS.
