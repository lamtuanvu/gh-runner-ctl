---
title: ghr
layout: hextra-home
---

{{< hextra/hero-badge >}}
  <div class="hx-w-2 hx-h-2 hx-rounded-full hx-bg-primary-400"></div>
  <span>Open Source</span>
  {{< icon name="arrow-circle-right" attributes="height=14" >}}
{{< /hextra/hero-badge >}}

<div class="hx-mt-6 hx-mb-6">
{{< hextra/hero-headline >}}
  Manage GitHub Actions&nbsp;<br class="sm:hx-block hx-hidden" />self-hosted runners with Docker
{{< /hextra/hero-headline >}}
</div>

<div class="hx-mb-12">
{{< hextra/hero-subtitle >}}
  Replace verbose docker-compose files with a single CLI.&nbsp;<br class="sm:hx-block hx-hidden" />Launch, scale, and tear down runners in seconds.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx-mb-6">
{{< hextra/hero-button text="Get Started" link="docs/getting-started/installation" >}}
</div>

```bash
ghr init        # one-time setup
ghr up 10       # launch 10 runners
ghr scale 20    # scale to 20
ghr down --all  # tear down everything
```

<div class="hx-mt-6"></div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Single Binary"
    subtitle="Install via Homebrew, apt, or a single shell script. No dependencies beyond Docker."
  >}}
  {{< hextra/feature-card
    title="Stateless Design"
    subtitle="No state files or databases. All runner state is tracked via Docker labels."
  >}}
  {{< hextra/feature-card
    title="Declarative Scaling"
    subtitle="Scale to an exact runner count with one command. ghr adds or removes as needed."
  >}}
  {{< hextra/feature-card
    title="Org & Repo Runners"
    subtitle="Register runners at the organization level or scoped to a single repository."
  >}}
  {{< hextra/feature-card
    title="Docker Context Aware"
    subtitle="Auto-detects OrbStack, colima, Docker Desktop, and other Docker contexts."
  >}}
  {{< hextra/feature-card
    title="Ephemeral Runners"
    subtitle="Built-in support for ephemeral runners that clean up after each job."
  >}}
{{< /hextra/feature-grid >}}
