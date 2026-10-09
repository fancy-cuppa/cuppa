# Community component research

Research for issue #26, performed 2026-10-09 from each project's GitHub page and the
[charm-and-friends/additional-bubbles](https://github.com/charm-and-friends/additional-bubbles)
list. "Verified" means stated on the project page; **unknown** means the page did not say, and the
component is added as a placeholder until checked against the module's `go.mod`.

Cuppa targets Bubble Tea / Lip Gloss **v2** (`charm.land/...`). Versions below were read from each module's `go.mod` on the Go proxy (latest release). A v1-only library can still be
*designed* in Cuppa (the preview is a drawing, not the real component), but generated code
(a later milestone) would need the v1 stack, so the status is recorded.

| Component | Repo | License | Bubble Tea | Decision |
|---|---|---|---|---|
| ntcharts (bar, line, sparkline, streamline, time series, heatmap, canvas) | NimbleMarkets/ntcharts | MIT | **v2 (verified)** | Added in the ntcharts entries (#25) |
| bubble-table | Evertras/bubble-table | MIT | **v2 (go.mod, v0.23.0)** | Add: `community.bubbletable` |
| stickers (FlexBox, Table) | 76creates/stickers | MIT | **v1 (verified)** | Add FlexBox: `community.flexbox` (v1) |
| bubbleboxer (layout tree) | treilik/bubbleboxer | MIT | **v1-era (v0.2.0 uses bubbletea v0.21)** | Add: `community.boxer` |
| bubble-datepicker | EthanEFung/bubble-datepicker | MIT | **v1-era (v0.1.1 uses bubbletea v0.24)** | Add: `community.datepicker` |
| bubbletea-overlay (modal) | rmhubbert/bubbletea-overlay | MIT | **v1 (v2 users should use Lip Gloss v2 compositing)** | Add as a modal: `community.overlay` (v1) |
| teacup (statusbar, filetree, ...) | knipferrc/teacup | MIT | **v1-era (v0.4.1 uses bubbletea v0.24)** | Add statusbar and filetree: `community.statusbar`, `community.filetree` |
| bubblezone | lrstanley/bubblezone | MIT | **v2 (verified)** | Non-visual (mouse regions): not a catalog entry; used as a note |
| promptkit | erikgeiser/promptkit | MIT | **v2 (go.mod, v0.12.0)** | Skip: duplicates Huh prompts (select, confirm, input) |
| bubbleup, bubbleo, bubblegum, bubbles-hlist, listExtensions, bubble-ssh, mritd/bubbles, ortizalec/bubbles, go-bubble-table, bubblelister | various | not checked | not checked | Backlog: re-evaluate when asked for; not verified |

## How an entry is chosen

1. It draws something visible (so it can be placed on a canvas). Non-visual helpers such as
   bubblezone are documented, not catalogued.
2. License is permissive (all verified ones are MIT).
3. It is not a duplicate of an official component.

## Follow-up

- Re-run the `go.mod` check (Go proxy `@latest`) when libraries release; v1-era rows are the ones to watch.
- Revisit the backlog rows once the catalog authoring guide (#41) exists.
