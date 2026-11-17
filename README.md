# reinvent: claims-reconciliation incident demo

A PagerDuty demo. A library upgrade makes the `claims-reconciliation` pod leak memory. The pod starves its Kubernetes node, which takes `patient-portal` and `scheduling` down and backs up claims with compliance exposure. Paige finds the change, a human approves, and workflows fix it in-cluster (restart the pod and apply a memory limit).

## Install (WSL2 on Windows 11)

This fork targets WSL2 only (tested on the Ubuntu distro it was developed in). You need Go and Git; Terraform is downloaded automatically if it isn't on your PATH.

```bash
git clone https://github.com/NickAtPD/reinvent.git ~/code/reinvent
cd ~/code/reinvent
./trigger
```

Keep the clone on the Linux filesystem (for example `~/code`), not under `/mnt/c` or `/mnt/g`: builds are much faster there, and the app stores an API token and Terraform state, which need real Linux file permissions.

**What you need:**
- A **PagerDuty API key** (PagerDuty → Integrations → API Access Keys). Reinvent asks for it, plus your region and login email, and checks each one.

Choose **SETUP → 2 Check this machine is ready** in the menu to confirm everything is in place. Links open in your Windows browser (via `wslview` from the `wslu` package, falling back to `explorer.exe`). Settings and state for a binary run outside a checkout live in `~/.local/share/reinvent`; in a checkout they stay in the repo.

Exports (menu **DEMO → 4**) go to `~/.local/share/reinvent/exports` and open in Explorer. They contain your real routing keys, so they stay on the Linux filesystem with mode 600 on purpose. For Postman, import the placeholder collection `payloads/rendered/reinvent.postman_collection.json` (web Postman works for this) and put your routing keys in a secret environment variable.

## Layout

| Dir | What |
|---|---|
| [`story/`](story/) | The narrative and demo beats (`README.md`), plus the original brief (`brief.md`) |
| [`instructions/`](instructions/) | Platform setup, in order, with the **automated vs manual matrix** |
| [`terraform/`](terraform/) | Services, business services, dependencies, EP, alert grouping, Event Orchestration, workflow shells |
| [`payloads/`](payloads/) | `scenario.json` (order, timing, shared vars) and the event templates. Edit these, not code. |
| [`trigger`](trigger) | Developers: run from a checkout. `./trigger` opens the same menu, using the repo's files |
| [`app/`](app/) | Go source (rebuilt automatically by `./trigger` when it changes) |
| [`runbook/`](runbook/) | The claims-reconciliation runbook to upload for Paige |
| [`k8s/`](k8s/) | Optional kind cluster and workloads for the remediation beat. Not used by the app for now; see instructions/03 |

## Quick start

You need Go (see Install). `./trigger` builds itself and fetches Terraform.

```bash
./trigger
```

That's it. A full-screen menu opens with two pages: **SETUP** (do once) and **DEMO** (use on stage). Switch with ← →. Once setup is done it opens straight on DEMO. ★ marks the next step. It covers everything:
settings (region US/EU, API key, your email, each checked against PagerDuty and saved to `.env`),
creating the PagerDuty objects, running, resolving, exporting the payloads, and teardown.
Anything that sends alerts or deletes something asks first.

Between `create` and the first run, do the manual steps in [instructions/](instructions/).

For scripts: `./trigger help` lists the subcommands (`setup`, `doctor`, `plan`, `create`, `run`, `resolve`, `export`, `destroy`).

## Editing the scenario

- **Wording, severities, custom details:** edit `payloads/alerts/*.json` or `payloads/change-event.json`. They are Go templates with `{{.Vars.x}}`, `{{now}}` and `{{ago "47m"}}`.
- **Order, timing, adding or removing an event:** edit `steps` in `payloads/scenario.json`.
- **Cluster, node, namespace, versions, library:** edit `vars` in `payloads/scenario.json`.
- **Then run `make test`** (or `./trigger run -dry-run`). It renders every step and checks that each one is a valid event.

The routing key is injected at send time from `terraform output routing_key` (or `PD_ROUTING_KEY`). It never goes in a payload file.

## Can and can't be automated

The short version, with the full matrix in [instructions/README.md](instructions/README.md):

- 🟢 **Automated:** services, business services, dependencies, routing key, change events, payloads, alert grouping, EO rule, EP (human), workflow names/descriptions/triggers, k8s workloads.
- 🔴 **Manual:** Paige on the EP, workflow steps, runner registration, Agent Connector, Slack/Paige scopes, Zoom/Meet, runbook upload, internal status page, console view, EA flags.
