# Findings: mnci Go flow (issue #5)

Observed with mnci 4.30.0 on Windows:

- `mnci new --into . --ci github --scope @cuppa --registry npm -y` bootstraps the repo and keeps `.git` and `LICENSE`.
- `mnci add go-internal-lib <name>` creates `libs/<name>` with its own `go.mod` and registers it in a root `go.work`. The starter slice is a package `libs/<name>/<name>/`.
- `mnci add go-app <name> --release` creates `apps/<name>` plus `build-all` / `package-all` targets that cross-compile six platforms and a `tools/go-app-release.cjs` that attaches the zips to the GitHub Release.
- Documentation mismatch: `mnci-details.md` says one root `go.mod`; the generator produces per-project modules and `go.work`.

Open items are tracked as issues on the MoNecromanCi repository when found.
