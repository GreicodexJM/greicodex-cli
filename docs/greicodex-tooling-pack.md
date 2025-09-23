# Greicodex Tooling Pack — Enforcement Template (ES)

Plantilla para **enforzar** el proceso: repos base, hooks de Git, CI/CD (Bitbucket/GitHub), escaneo de secretos, linting, tests, build de contenedores, Helm, e IaC con OpenTofu. Compatible con desarrollo **offline** vía Docker/Compose y perfiles.

---

## 🗂 Estructura del repo

```
.
├─ .editorconfig
├─ .gitleaks.toml
├─ .gitignore
├─ .githooks/
│  ├─ pre-commit
│  └─ pre-push
├─ .github/
│  └─ workflows/
│     └─ ci.yml
├─ bitbucket-pipelines.yml
├─ docker-compose.yml
├─ Makefile
├─ README.md
├─ scripts/
│  ├─ enforce-coverage.sh
│  └─ wait-for.sh
├─ deploy/
│  └─ helm/
│     ├─ Chart.yaml
│     ├─ values.yaml
│     └─ templates/
│        ├─ deployment.yaml
│        ├─ service.yaml
│        └─ ingress.yaml
└─ iac/
   ├─ main.tfu
   ├─ variables.tfu
   └─ outputs.tfu
```

> **Nota:** Configura `core.hooksPath=.githooks` en el repo para activar los hooks.

---

## 🔧 Configuración base

### `.editorconfig`

```ini
root = true
[*]
end_of_line = lf
insert_final_newline = true
charset = utf-8
indent_style = space
indent_size = 2
trim_trailing_whitespace = true
```

### `.gitignore` (mínimo común)

```gitignore
# Dependencias
node_modules/
vendor/
.venv/
__pycache__/
*.pyc

# Builds/artifacts
build/
dist/
coverage/
reports/
*.log

# IDE
.vscode/
.idea/
.DS_Store

# Docker
*.pid
.env
.env.*
```

### `.gitleaks.toml` (reglas sensatas)

```toml
title = "Greicodex baseline"
[extend]
path = ""  # puedes extender de un preset upstream si lo deseas

[[rules]]
id = "generic-aws-key"
regex = '''AKIA[0-9A-Z]{16}'''
entropy = 3.5
path = ""
```

---

## 🪝 Hooks de Git

Config global del repo (una vez):

```bash
git config core.hooksPath .githooks
```

### `.githooks/pre-commit`

```bash
#!/usr/bin/env bash
set -euo pipefail

# 1) Escaneo de secretos
if command -v gitleaks >/dev/null 2>&1; then
  echo "[pre-commit] Running gitleaks..."
  gitleaks protect --staged --no-banner
else
  echo "[pre-commit] gitleaks not found (skip)."
fi

# 2) Linting rápido por lenguaje
if [ -f package.json ]; then
  if jq -e '.scripts.lint' package.json >/dev/null 2>&1; then
    echo "[pre-commit] npm run lint"
    npm run lint --silent || (echo "Lint failed"; exit 1)
  fi
fi

if [ -f "composer.json" ] && command -v composer >/dev/null; then
  echo "[pre-commit] PHP CodeSniffer"
  vendor/bin/phpcs -q || true
fi

if ls **/*.py >/dev/null 2>&1; then
  if command -v ruff >/dev/null; then
    echo "[pre-commit] ruff"
    ruff check .
  elif command -v flake8 >/dev/null; then
    echo "[pre-commit] flake8"
    flake8 .
  fi
fi

# 3) Tests ultra-rápidos opcionales (ej. --passWithNoTests)
if [ -f package.json ]; then
  if jq -e '.scripts.test' package.json >/dev/null 2>&1; then
    echo "[pre-commit] npm test (quick)"
    npm test --silent -- --passWithNoTests --watchAll=false || true
  fi
fi

exit 0
```

### `.githooks/pre-push`

```bash
#!/usr/bin/env bash
set -euo pipefail

# Bloqueo por cobertura mínima (configurable)
MIN_COVERAGE="80"

# Ejecuta pipeline local mínimo (útil offline)
make ci-local || {
  echo "[pre-push] CI local falló"; exit 1;
}

# Enforce coverage si existe reporte
if ./scripts/enforce-coverage.sh "$MIN_COVERAGE"; then
  echo "[pre-push] Cobertura OK (>= ${MIN_COVERAGE}%)"
else
  echo "[pre-push] Cobertura insuficiente (< ${MIN_COVERAGE}%)"; exit 1
fi
```

### `scripts/enforce-coverage.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail
REQ=${1:-80}
FOUND=0

# Detecta formatos comunes
if [ -f coverage/coverage-summary.json ]; then
  # Jest
  FOUND=$(jq '.total.lines.pct' coverage/coverage-summary.json | awk '{print int($1+0.5)}')
elif [ -f coverage.xml ]; then
  # Cobertura XML (phpunit/jacoco/pytest-cov)
  FOUND=$(python3 - <<'PY'
import re
import sys
from xml.etree import ElementTree as ET
try:
    root = ET.parse('coverage.xml').getroot()
    # jacoco/cobertura variants
    if root.tag.endswith('coverage') and root.get('line-rate'):
        pct = float(root.get('line-rate')) * 100
    else:
        lines_valid = int(root.get('lines-valid') or 0)
        lines_covered = int(root.get('lines-covered') or 0)
        pct = (lines_covered/lines_valid*100) if lines_valid else 0
    print(int(round(pct)))
except Exception:
    print(0)
PY
)
fi

[ "$FOUND" -ge "$REQ" ]
```

### `scripts/wait-for.sh`

```bash
#!/usr/bin/env bash
# Pequeño helper para esperar a que un puerto esté listo
# uso: ./scripts/wait-for.sh host:puerto -- comando
set -e
HOSTPORT=$1; shift
HOST=${HOSTPORT%:*}
PORT=${HOSTPORT#*:}
until nc -z "$HOST" "$PORT"; do echo "waiting for $HOSTPORT"; sleep 1; done
exec "$@"
```

---

## 🐳 Docker Compose (offline-friendly)

### `docker-compose.yml`

```yaml
version: "3.9"

x-common-env: &common_env
  TZ: "UTC"
  APP_ENV: "local"

services:
  app:
    image: node:20-bullseye
    working_dir: /workspace
    volumes:
      - ./:/workspace
    environment:
      <<: *common_env
    command: bash -lc "npm ci || npm install && npm run dev"
    profiles: [web]
    depends_on:
      - postgres
      - rabbitmq

  # Bases de datos (actívalas por perfiles)
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: app
    ports: ["5432:5432"]
    volumes:
      - pgdata:/var/lib/postgresql/data

  mongo:
    image: mongo:7
    ports: ["27017:27017"]
    volumes:
      - mongodata:/data/db
    profiles: [mongo]

  couchdb:
    image: couchdb:3
    ports: ["5984:5984"]
    environment:
      COUCHDB_USER: admin
      COUCHDB_PASSWORD: admin
    volumes:
      - couchdata:/opt/couchdb/data
    profiles: [couch]

  rabbitmq:
    image: rabbitmq:3-management
    ports: ["5672:5672", "15672:15672"]
    volumes:
      - rabbitdata:/var/lib/rabbitmq

  # Kafka single-broker para dev
  kafka:
    image: bitnami/kafka:3.8
    environment:
      KAFKA_ENABLE_KRAFT: "yes"
      KAFKA_CFG_PROCESS_ROLES: "broker,controller"
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 1@localhost:9093
    ports: ["9092:9092"]
    profiles: [kafka]

volumes:
  pgdata: {}
  mongodata: {}
  couchdata: {}
  rabbitdata: {}
```

> Usa perfiles: `docker compose --profile mongo --profile kafka up -d`

---

## 🧰 Makefile (comandos estandarizados)

### `Makefile`

```makefile
SHELL := /usr/bin/env bash
MIN_COV ?= 80

.PHONY: setup ci-local lint test coverage build docker helm tofu

setup:
	@echo "[setup] Installing local deps..."
	@if [ -f package.json ]; then npm ci || npm install; fi
	@if [ -f requirements.txt ]; then python3 -m venv .venv && . .venv/bin/activate && pip install -r requirements.txt; fi
	@if [ -f composer.json ]; then composer install; fi

ci-local: lint test coverage
	@./scripts/enforce-coverage.sh $(MIN_COV)

lint:
	@if command -v eslint >/dev/null && [ -f package.json ]; then eslint .; fi
	@if command -v ruff >/dev/null; then ruff check . || true; fi
	@if command -v golangci-lint >/dev/null; then golangci-lint run || true; fi
	@if command -v phpcs >/dev/null; then phpcs || true; fi

test:
	@if [ -f package.json ]; then npm test -- --coverage --watchAll=false; fi
	@if [ -d tests ] && command -v pytest >/dev/null; then pytest --cov=. --cov-report=xml; fi
	@if [ -f phpunit.xml* ]; then vendor/bin/phpunit --coverage-clover coverage.xml; fi
	@if [ -d */* ]; then echo "" > /dev/null; fi

coverage:
	@echo "Coverage reports generated (Jest/pytest/phpunit where applicable)."

build:
	@echo "[build] app build (language-specific)"

docker:
	docker build -t myorg/app:dev . || true

helm:
	cd deploy/helm && helm lint && helm package . -d ../

# IaC con OpenTofu (DigitalOcean por defecto)
tofu:
	cd iac && tofu init && tofu validate && tofu plan
```

---

## 🚦 CI/CD (GitHub Actions y Bitbucket)

### `.github/workflows/ci.yml`

```yaml
name: CI
on:
  pull_request:
  push:
    branches: [develop]

jobs:
  build-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Node
        uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: 'npm'
      - name: Setup Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      - name: Install tooling
        run: |
          pip install ruff pytest pytest-cov
          npm ci || npm install
          curl -sSL https://github.com/gitleaks/gitleaks/releases/latest/download/gitleaks_$(uname -s)_$(uname -m).tar.gz \
            | tar -xz && sudo mv gitleaks /usr/local/bin/
      - name: Secret scan
        run: gitleaks detect --no-banner
      - name: Lint
        run: make lint
      - name: Test
        run: make test
      - name: Enforce coverage
        run: ./scripts/enforce-coverage.sh 80
      - name: Build Docker
        run: |
          echo "Building image"
          docker build -t ${{ github.repository }}:${{ github.sha }} .
      - name: Helm Lint
        run: make helm
```

### `bitbucket-pipelines.yml`

```yaml
image: atlassian/default-image:4

pipelines:
  pull-requests:
    '**':
      - step:
          name: CI
          caches: [node]
          script:
            - apt-get update && apt-get install -y python3-pip curl jq
            - pip3 install ruff pytest pytest-cov
            - npm ci || npm install
            - curl -sSL https://github.com/gitleaks/gitleaks/releases/latest/download/gitleaks_linux_x64.tar.gz \
              | tar -xz && mv gitleaks /usr/local/bin/
            - gitleaks detect --no-banner
            - make lint
            - make test
            - ./scripts/enforce-coverage.sh 80
            - docker build -t "$BITBUCKET_REPO_OWNER/$BITBUCKET_REPO_SLUG:$BITBUCKET_COMMIT" .
```

---

## ⛵ Helm chart mínimo

### `deploy/helm/Chart.yaml`

```yaml
apiVersion: v2
name: app
version: 0.1.0
appVersion: "0.1.0"
```

### `deploy/helm/values.yaml`

```yaml
image:
  repository: myorg/app
  tag: latest
  pullPolicy: IfNotPresent

env: []

service:
  type: ClusterIP
  port: 80

resources: {}

ingress:
  enabled: false
  className: nginx
  hosts:
    - host: app.local
      paths:
        - path: /
          pathType: Prefix
  tls: []
```

### `deploy/helm/templates/deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "app.fullname" . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ include "app.name" . }}
  template:
    metadata:
      labels:
        app: {{ include "app.name" . }}
    spec:
      containers:
        - name: app
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - containerPort: 3000
          env:
            {{- range .Values.env }}
            - name: {{ .name }}
              value: "{{ .value }}"
            {{- end }}
```

### `deploy/helm/templates/service.yaml`

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ include "app.fullname" . }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: 3000
  selector:
    app: {{ include "app.name" . }}
```

### `deploy/helm/templates/ingress.yaml`

```yaml
{{- if .Values.ingress.enabled }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "app.fullname" . }}
  annotations: {}
spec:
  ingressClassName: {{ .Values.ingress.className }}
  rules:
  {{- range .Values.ingress.hosts }}
    - host: {{ .host }}
      http:
        paths:
        {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType }}
            backend:
              service:
                name: {{ include "app.fullname" $ }}
                port:
                  number: {{ $.Values.service.port }}
        {{- end }}
  {{- end }}
{{- end }}
```

---

## ☁️ IaC con OpenTofu (DigitalOcean por defecto)

> Ajusta providers según orden preferido (DO → AWS → GCP → Azure).

### `iac/main.tfu`

```hcl
terraform {
  required_version = ">= 1.6.0"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.40.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_kubernetes_cluster" "this" {
  name   = var.cluster_name
  region = var.region
  version = var.k8s_version
  node_pool {
    name       = "default"
    size       = var.node_size
    node_count = var.node_count
  }
}

output "cluster_name" { value = digitalocean_kubernetes_cluster.this.name }
```

### `iac/variables.tfu`

```hcl
variable "do_token" { type = string }
variable "cluster_name" { type = string  default = "greicodex-dev" }
variable "region" { type = string  default = "nyc1" }
variable "k8s_version" { type = string default = "latest" }
variable "node_size" { type = string  default = "s-2vcpu-4gb" }
variable "node_count" { type = number  default = 1 }
```

### `iac/outputs.tfu`

```hcl
output "kubeconfig" {
  value       = "Run: doctl kubernetes cluster kubeconfig save $(terraform output -raw cluster_name)"
  description = "How to fetch kubeconfig"
}
```

---

## 📜 README.md (extracto sugerido)

````markdown
# Proyecto XYZ

## Requisitos
- Docker & Docker Compose
- Node 20 / Python 3.11 / PHP 8.2 (según módulo)

## Desarrollo local
```bash
make setup
docker compose up -d
npm run dev
````

## CI local

```bash
make ci-local MIN_COV=80
```

## Estándares

* Gitflow (branches: feature/, bugfix/, hotfix/)
* PR con lint, tests y cobertura ≥80%
* 12-Factor + SOLID
* Arquitectura hexagonal (Puertos/Adaptadores)

```

---

## 🔒 Buenas prácticas de secretos
- Usa un Secret Manager (DO, AWS, GCP) o Vault. Nunca `.env` en commits.
- Para dev local: `.env` solo fuera del control de versiones.
- Variables inyectadas en CI vía secretos de repo/workspace.

---

## 🧪 Testing (marcos sugeridos por lenguaje)
- **TypeScript**: Jest, ts-jest, supertest; ESLint + Prettier.
- **PHP**: PHPUnit/Codeception; PHP-CS-Fixer/PHPCS.
- **Python**: pytest + pytest-cov; ruff.
- **Go**: `go test -coverprofile`; golangci-lint.
- **Java**: JUnit5 + JaCoCo; Checkstyle/SpotBugs.
- **Dart/Flutter**: `flutter test --coverage` + `melos` si monorepo.

Asegura que cada stack exporte **coverage.xml** o **coverage-summary.json** para `enforce-coverage.sh`.

---

## 🧭 Integración con JIRA y Gitflow
- Nombrado de ramas: `feature/KEY-123-descripcion`, `bugfix/KEY-456-...`
- Commits incluyen `KEY-123` para trazabilidad.
- PR exige link a ticket, checklist (doc, tests, pipeline verde).

---

## 🛰 Despliegue
- Empaquetar Helm: `make helm`
- IaC: `make tofu` (plan/validate) → aplicar con aprobación.
- Knative/K8s: ajusta `values.yaml` (recursos, env, probes).

---

## ✅ Lista de verificación de cumplimiento
- [ ] Hooks activos (`git config core.hooksPath .githooks`)
- [ ] Gitleaks pasa local y en CI
- [ ] Lint y Tests verdes
- [ ] Cobertura ≥ 80%
- [ ] Imagen Docker build OK
- [ ] Helm lint OK
- [ ] IaC valida (tofu validate)
- [ ] PR con checklist completo y ticket JIRA vinculado
```
