# Kubernetes Manifests — Kustomize Foundation

This directory contains Kubernetes manifests for ms-casino-go-v2 organized using Kustomize (Phase 1 of ADR-001).

## Overview

**Kustomize** allows us to maintain a single source of truth (`base/`) with environment-specific variations (`overlays/`) without template logic or code duplication.

- **`base/`** — Shared resources (Deployment, Service, HTTPRoute, ConfigMap)
- **`overlays/{dev,hml,prd}/`** — Environment-specific patches (replicas, resources, environment variables)
- **`components/`** — Reusable feature bundles (Datadog sidecar, NetworkPolicy)

### Phase 1 vs Phase 2

**Phase 1 (Current):** Standard `Deployment`, single `Service`, per-env overlays.  
**Phase 2 (Deferred to TEMPL-001.9):** Argo Rollouts for progressive delivery (canary steps, automatic rollback, metric-based promotion).

See [ADR-001](../../../docs/architecture/project-decisions/ADR-001-k8s-manifest-strategy.md) for the full strategy.

---

## Directory Structure

```
k8s/
├── base/
│   ├── deployment.yaml          # Standard Deployment (replicas patched per env)
│   ├── service.yaml             # ClusterIP Service
│   ├── httproute.yaml           # Gateway API HTTPRoute (routes /api/... traffic)
│   ├── configmap.yaml           # Shared env vars (LOG_LEVEL, APP_ENV, etc.)
│   └── kustomization.yaml       # Kustomize base definition
│
├── overlays/
│   ├── dev/                     # Development environment
│   │   ├── kustomization.yaml   # 2 replicas, 250m CPU, 256Mi mem, LOG_LEVEL=DEBUG, Datadog disabled
│   │   ├── deployment-patch.yaml
│   │   └── configmap-patch.yaml
│   │
│   ├── hml/                     # Staging (homologation) environment
│   │   ├── kustomization.yaml   # 3 replicas, 500m CPU, 512Mi mem, LOG_LEVEL=INFO, Datadog enabled
│   │   ├── deployment-patch.yaml
│   │   ├── configmap-patch.yaml
│   │   └── pdb.yaml             # PodDisruptionBudget minAvailable: 2
│   │
│   └── prd/                     # Production environment
│       ├── kustomization.yaml   # 3 replicas (HPA 3-10), 1.0 CPU, 1Gi mem, LOG_LEVEL=INFO, Datadog enabled
│       ├── deployment-patch.yaml
│       ├── configmap-patch.yaml
│       ├── hpa.yaml             # HorizontalPodAutoscaler (3-10 replicas, 80% CPU/mem targets)
│       ├── pdb.yaml             # PodDisruptionBudget minAvailable: 3
│       └── (networkpolicy.yaml is applied via component)
│
└── components/
    ├── datadog-sidecar/         # Datadog APM + logging agent sidecar
    │   └── kustomization.yaml   # Referenced by hml and prd overlays
    │
    └── networkpolicy/           # default-deny + allow-from-gateway NetworkPolicy
        ├── kustomization.yaml   # Referenced by prd overlay only
        └── networkpolicy.yaml
```

---

## Deploying to Each Environment

### Development
```bash
kubectl apply -k overlays/dev
```

Expected output:
- Deployment with 2 replicas
- Service (ClusterIP)
- HTTPRoute
- ConfigMap with LOG_LEVEL=DEBUG, DATADOG_ENABLED=false

### Staging (Homologation)
```bash
kubectl apply -k overlays/hml
```

Expected output:
- Deployment with 3 replicas
- Service (ClusterIP)
- HTTPRoute
- ConfigMap with LOG_LEVEL=INFO, DATADOG_ENABLED=true
- Datadog agent sidecar container
- PodDisruptionBudget (minAvailable: 2)

### Production
```bash
kubectl apply -k overlays/prd
```

Expected output:
- Deployment with 3 initial replicas
- Service (ClusterIP)
- HTTPRoute
- ConfigMap with LOG_LEVEL=INFO, DATADOG_ENABLED=true
- Datadog agent sidecar container
- HorizontalPodAutoscaler (3-10 replicas, 80% CPU/mem targets)
- PodDisruptionBudget (minAvailable: 3)
- NetworkPolicy (default-deny + allow-from-gateway-namespace)

---

## Viewing Rendered Manifests

Before applying, render the manifests to verify patches and components:

```bash
# Dev
kubectl kustomize overlays/dev

# Staging
kubectl kustomize overlays/hml

# Production
kubectl kustomize overlays/prd
```

---

## Common Patching Patterns

### 1. Changing Replica Count

Edit the `deployment-patch.yaml` in the desired overlay:

```yaml
patches:
  - target:
      group: apps
      version: v1
      kind: Deployment
      name: ms-casino-go-v2
    patch: |-
      - op: replace
        path: /spec/replicas
        value: 5
```

### 2. Changing Resource Requests/Limits

Edit the `deployment-patch.yaml` in the desired overlay:

```yaml
patches:
  - target:
      group: apps
      version: v1
      kind: Deployment
      name: ms-casino-go-v2
    patch: |-
      - op: replace
        path: /spec/template/spec/containers/0/resources
        value:
          requests:
            memory: "2Gi"
            cpu: "2000m"
          limits:
            memory: "4Gi"
            cpu: "4000m"
```

### 3. Adding Environment Variables

Edit the `configmap-patch.yaml` in the desired overlay:

```yaml
patchesJson6902:
  - target:
      group: ""
      version: v1
      kind: ConfigMap
      name: ms-casino-go-v2-config
    patch: |-
      - op: add
        path: /data/MY_NEW_VAR
        value: "my-value"
```

### 4. Enabling/Disabling Datadog

The `datadog-sidecar` component is referenced in `hml` and `prd` overlays:

```yaml
# In overlays/{hml,prd}/kustomization.yaml:
components:
  - ../../components/datadog-sidecar
```

To disable Datadog in an overlay, remove this line.

### 5. Adding NetworkPolicy

The `networkpolicy` component is referenced in the `prd` overlay:

```yaml
# In overlays/prd/kustomization.yaml:
components:
  - ../../components/networkpolicy
```

To add NetworkPolicy to another environment, add this line to that overlay's `kustomization.yaml`.

---

## Adding a New Component

1. Create a directory under `components/`:
   ```bash
   mkdir components/my-feature
   cd components/my-feature
   ```

2. Create `kustomization.yaml` (kind: Component):
   ```yaml
   apiVersion: kustomize.config.k8s.io/v1beta1
   kind: Component

   resources:
     - my-resource.yaml

   # Or use patches:
   patchesJson6902:
     - target:
         group: apps
         version: v1
         kind: Deployment
         name: ms-casino-go-v2
       patch: |-
         - op: add
           path: /spec/template/spec/...
           value: ...
   ```

3. Reference the component in the overlay's `kustomization.yaml`:
   ```yaml
   components:
     - ../../components/my-feature
   ```

4. Test rendering:
   ```bash
   kubectl kustomize overlays/{dev,hml,prd}
   ```

---

## Extending an Overlay

To add a new file to a specific environment (e.g., a custom HPA for `hml`):

1. Create the file in the overlay directory:
   ```bash
   cat > overlays/hml/custom-hpa.yaml << EOF
   apiVersion: autoscaling/v2
   kind: HorizontalPodAutoscaler
   metadata:
     name: ms-casino-go-v2-custom-hpa
   spec:
     ...
   EOF
   ```

2. Reference it in the overlay's `kustomization.yaml`:
   ```yaml
   resources:
     - pdb.yaml
     - custom-hpa.yaml
   ```

3. Test rendering:
   ```bash
   kubectl kustomize overlays/hml
   ```

---

## Template Variables

All template variables are preserved:

- `ms-casino-go-v2` — Project name (e.g., `my-api-service`)
- `default` — Kubernetes namespace (e.g., `default`)
- `<ECR_REPO>` — ECR repository URI (e.g., `123456789.dkr.ecr.us-east-1.amazonaws.com`)
- `8081` — Application port (e.g., `8080`)
- `v2` — API version prefix (e.g., `v1`)
- `internal-gateway` — Gateway name (e.g., `internal-gateway`)
- `gateway` — Gateway namespace (e.g., `ingress-system`)

These are automatically replaced during project scaffolding via Copier.

---

## Validation

### kubectl kustomize

Render manifests without applying:

```bash
kubectl kustomize overlays/{dev,hml,prd}
```

### kubectl apply --dry-run

Validate against the cluster schema:

```bash
kubectl apply -k overlays/dev --dry-run=client
```

---

## Troubleshooting

### Kustomize render fails

```bash
# Check syntax errors
kubectl kustomize overlays/dev 2>&1 | head -20

# Verify resource references
kubectl kustomize overlays/dev | grep -A5 "apiVersion"
```

### Apply fails with schema errors

```bash
# Validate specific resource
kubectl kustomize overlays/prd | kubectl apply -f - --dry-run=client

# Check logs
kubectl logs -n default deployment/ms-casino-go-v2
```

### Template variables not replaced

Template variables are replaced during project scaffolding, not at deployment time. If you see `ms-casino-go-v2` in rendered output, the project was not scaffolded correctly.

---

## Progressive Delivery (Phase 2 — TEMPL-001.9)

This Phase 1 uses standard `Deployment` and single `Service`. **Phase 2** (deferred to story TEMPL-001.9) will:

- Replace `Deployment` → `Rollout` (Argo Rollouts CRD)
- Split `Service` → `service-stable.yaml` + `service-canary.yaml`
- Add weighted backends in HTTPRoute
- Add `AnalysisTemplate` for metric-based automatic rollback
- Add manual promotion gates for canary steps

Phase 2 will **not** restructure the overlays—it layers on top without breaking existing patches.

See [ADR-001 Sub-Decision B](../../../docs/architecture/project-decisions/ADR-001-k8s-manifest-strategy.md#sub-decision-b-progressive-delivery-argo-rollouts) for the full plan.

---

## References

- **Kustomize Documentation:** https://kustomize.io/
- **ADR-001 (Architecture Decision):** `docs/architecture/project-decisions/ADR-001-k8s-manifest-strategy.md`
- **Story TEMPL-001.6 (This story):** `docs/stories/active/TEMPL-001.6.k8s-kustomize-foundation.story.md`
- **Story TEMPL-001.9 (Phase 2 — Argo Rollouts):** `docs/stories/active/TEMPL-001.9.k8s-progressive-delivery.story.md`
- **Gateway API HTTPRoute:** https://gateway-api.sigs.k8s.io/

---

**Last Updated:** 2026-05-14  
**Status:** Phase 1 (Kustomize Foundation)
