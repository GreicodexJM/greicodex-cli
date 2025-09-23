# Receta de Configuración y Mantenimiento de Proyectos

*(Para proyectos en Dart, TypeScript, PHP, Python, Golang, Java — SaaS & Outsourcing)*

---

## 1. Checklist para **Nuevo Proyecto**

### 🏗 **Fundación del Proyecto**

* [ ] **Épica y Tickets en JIRA**: Crear épica para el proyecto y dividir en hitos.
* [ ] **Repositorio**: Crear repo en Git (Bitbucket/GitHub). Inicializar ramas `main` + `develop` (gitflow).
* [ ] **README.md**: Breve pero completo: propósito, stack técnico, pasos de instalación local.
* [ ] **.gitignore & .editorconfig**: Config estándar según lenguaje.
* [ ] **LICENSE & CONTRIBUTING.md**: Usar los de la compañía.
* [ ] **Docker/Docker-Compose**: Entorno local con DB, mensajería y dependencias. Debe funcionar offline.

### ⚙️ **Lenguajes & Frameworks Estándar**

* **Dart/Flutter**: Arquitectura por capas.
* **TypeScript**: Preferir Next/Angular/Nuxt/React según el caso.
* **PHP**: Symfony/Laravel.
* **Python**: FastAPI (microservicios).
* **Golang**: Para backends de alto rendimiento.
* **Java**: Spring Boot.
* [ ] Siempre usar **arquitectura hexagonal** (Puertos y Adaptadores).

### 📦 **Dependencias**

* [ ] Versiones fijas en `package.json`, `composer.json`, `requirements.txt`, `go.mod`, `pom.xml`.
* [ ] Usar espejos privados/cachés para trabajo sin internet (npm, PyPI, Maven, Composer, etc.).
* [ ] Escaneo de vulnerabilidades (`npm audit`, `snyk`, etc.).

### 🧪 **Pruebas**

* [ ] Framework de pruebas unitarias configurado (jest, phpunit, codeception, JUnit, pytest, etc.).
* [ ] Cobertura mínima del 80% obligatoria.
* [ ] Adaptadores simulados para servicios externos.

### 🚦 **Pipelines (CI/CD)**

* [ ] Bitbucket Pipelines/GitHub Actions configurados:

  * Escaneo de secretos (Trufflehog, Gitleaks).
  * Linting (ESLint, Flake8, GoLint, Checkstyle, etc.).
  * Tests con cobertura.
  * Generación de artefactos (imagen Docker, chart de Helm).
* [ ] Job de despliegue (DigitalOcean → AWS → GCP → Azure en ese orden).
* [ ] Helm chart dentro de `/deploy/helm`.

### 📑 **Documentación**

* [ ] Documentos en repo (`/docs`).
* [ ] ADRs (Architecture Decision Records) para decisiones importantes.
* [ ] Diagramas o archivos pesados en Confluence o Google Drive.

---

## 2. Checklist para **Mantenimiento de Proyectos Existentes**

### 🔍 **Auditoría & Actualización**

* [ ] Traer últimos cambios de `main` & `develop`.
* [ ] Actualizar dependencias (minor mensual, major trimestral).
* [ ] Re-ejecutar escaneos de vulnerabilidades.
* [ ] Revisar pipelines CI/CD (acciones, plugins actualizados).
* [ ] Validar imágenes Docker (base actualizada).

### 🧹 **Calidad de Código & Arquitectura**

* [ ] Confirmar hexagonal + principios SOLID.
* [ ] Refactorizar código acoplado a frameworks.
* [ ] Asegurar que adaptadores sean reemplazables/testeables.

### 🧪 **QA & Testing**

* [ ] Mantener cobertura >80%.
* [ ] Pruebas E2E automáticas.
* [ ] QA manual solo si es imposible automatizar.
* [ ] QA debe sincronizarse con REQ para actualizar protocolos de pruebas.

### 🚀 **Operaciones & Entrega**

* [ ] Validar despliegues vía Helm charts en staging y prod.
* [ ] Revisar IaC (OpenTofu) contra la infraestructura real.
* [ ] Probar backup/restore de bases de datos.

### 📊 **Seguimiento & Métricas**

* [ ] Registrar horas en Cronus (diario).
* [ ] Generar métricas automáticas semanales (velocidad, cobertura, frecuencia de despliegue, MTTR).
* [ ] Revisar anomalías (retrasos, fallos de pruebas, sobrecarga).

---

## 3. Reglas Transversales (**Aplica a Todos los Proyectos**)

* Cumplir con los **12 Factores** y principios **SOLID**.
* Documentación **junto al código**, no en correos ni en cabezas.
* Ningún merge sin PR + revisión de pares + QA.
* Cada release debe ser reproducible (pipelines + contenedores).
* Nada de hotfix en producción sin ticket/branch.
* Manejo de secretos solo en vaults, nunca en repositorio.

---

