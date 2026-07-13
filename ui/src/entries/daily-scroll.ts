// The daily view's only JavaScript: it appends older windows of the ledger as
// the foot scrolls into view. Not a React island — there is nothing to render.
//
// Everything it shows comes from the server. The window it fetches is HTML from
// the same Go partial that drew the first one, and the foot's four states are
// server-rendered and revealed by data-state, so no markup and no Tailwind live
// in here. The cursor is server-issued too: its absence is the terminal, so
// this file owns no stop condition and does no date arithmetic.

const foot = document.querySelector<HTMLElement>('[data-daily-scroll]');

if (foot) {
  let start = foot.dataset.nextStart;
  let end = foot.dataset.nextEnd;

  const observer = new IntersectionObserver((entries) => {
    if (entries.some((entry) => entry.isIntersecting)) void load();
  });

  const load = async () => {
    if (!start || !end) return;
    observer.disconnect(); // one fetch in flight; re-armed only on success
    foot.dataset.state = 'loading';

    try {
      const res = await fetch(`/daily/older?start=${start}&end=${end}`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      foot.insertAdjacentHTML('beforebegin', await res.text());

      const appended = foot.previousElementSibling as HTMLElement;
      start = appended.dataset.nextStart;
      end = appended.dataset.nextEnd;
    } catch {
      // The observer stays disconnected through a failure: retrying is the
      // user's call, so a flaky connection cannot emit a burst of failing
      // requests as the sentinel drifts in and out of view.
      foot.dataset.state = 'error';
      return;
    }

    if (start && end) {
      foot.dataset.state = 'idle';
      observer.observe(foot);
    } else {
      foot.dataset.state = 'end';
    }
  };

  foot.querySelector('[data-retry]')?.addEventListener('click', () => void load());

  if (start && end) observer.observe(foot);
}
