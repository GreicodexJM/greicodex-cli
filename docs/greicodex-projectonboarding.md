
**"Protocolo de Iniciación de Proyectos"**. 

Su objetivo es simple: estandarizar la creación de un entorno de desarrollo completo y profesional para que, cuando el primer desarrollador reciba el ticket `DEV-001`, tenga todo lo que necesita para empezar a crear valor, no a configurar herramientas.

Aquí está la plantilla. Este documento se convierte en el primer artefacto del Módulo 0 de tu GOS (Greicodex Operating System).

---
### **Plantilla: Protocolo de Iniciación de Proyectos Greicodex (M0P1)**

**Objetivo:** Estandarizar la creación de un nuevo proyecto, asegurando que cada repositorio nazca con la documentación, infraestructura y automatización necesarias para cumplir con los estándares de "Ejecución Impecable" del GOS.

**Responsable:** Líder de Proyecto (PM) / Líder Técnico

---
### **Fase 1: Setup Inicial y Control de Acceso (Día 1)**

* `[ ]` **Crear Repositorio:** Crear un nuevo repositorio en Bitbucket/GitHub siguiendo la convención de nombres: `[nombre-constelacion]-[adjetivo-proyecto]`.
* `[ ]` **Inicializar con README:** Inicializar el repositorio con un archivo `README.md` que contenga una breve descripción del proyecto y un link al Acta de Constitución en Confluence.
* `[ ]` **Configurar Ramas:**
    * Crear la rama `develop` a partir de `main`.
    * Configurar `main` como la rama protegida (solo se puede hacer merge a través de Pull Requests).
* `[ ]` **Gestionar Accesos:** Dar acceso al repositorio al equipo de desarrollo asignado con los permisos correspondientes (ej: "Write" para Devs, "Admin" para Líder Técnico).
* `[ ]` **Crear Proyecto en Jira:** Crear el proyecto en Jira, configurar el tablero (Kanban/Scrum) y enlazarlo con el repositorio de Bitbucket.

---
### **Fase 2: La Fundación de la Confianza - El "Memory Bank" (Día 1-2)**

* `[ ]` **Crear la Estructura de Carpetas:** En la raíz del repositorio, crear la carpeta `/memory-bank`.
* `[ ]` **Poblar el `01_PROJECT_CHARTER.md`:** Copiar el contenido final del "Acta de Constitución del Proyecto" acordada con el cliente. Este documento define el **"Porqué"**.
* `[ ]` **Poblar el `02_ARCHITECTURE_PRINCIPLES.md`:** Copiar el contenido de tu documento `hexagonal-solid-architecture.md`. Este documento define el **"Cómo"** construimos.
* `[ ]` **Poblar el `03_AGENTIC_WORKFLOW.md`:** Copiar el contenido de tu documento `sparc-agentic-development.md`. Este documento define **nuestra metodología**.
* `[ ]` **Crear el `04_USER_STORIES.md` (Borrador):** Crear el archivo y añadir las primeras 2-3 historias de usuario de más alto nivel que se desprenden del Acta de Constitución. Este será un documento vivo que el PM actualizará.

---
### **Fase 3: La Arquitectura como Código - IaaC (Día 2-3)**

* `[ ]` **Definir la Infraestructura Inicial:** En una nueva carpeta `/infrastructure`, crear los archivos de IaaC (Terraform/Pulumi/Bicep) para provisionar el entorno de pruebas inicial.
    * **Mínimo indispensable:** Base de datos, red virtual, y el contenedor o servicio donde correrá la aplicación.
* `[ ]` **Variables de Entorno:** Crear un archivo `env.template` en la raíz del proyecto que liste todas las variables de entorno que la aplicación necesitará para correr, pero sin los valores secretos (ej: `DB_HOST=`, `DB_PASSWORD=`).
* `[ ]` **Gestión de Secretos:** Configurar el gestor de secretos (ej: Vault, AWS/Azure Key Vault, Bitbucket Pipelines secrets) con los valores iniciales para el entorno de pruebas.

---
### **Fase 4: La Automatización - Pipelines CI/CD (Día 3-4)**

* `[ ]` **Crear el Pipeline de CI (Integración Continua):**
    * Crear el archivo `bitbucket-pipelines.yml` o el workflow de GitHub Actions.
    * Configurar el pipeline para que se active con cada `push` a una rama de desarrollo.
    * **Pasos mínimos del pipeline de CI:**
        1.  Instalar dependencias.
        2.  Ejecutar el linter (análisis de estilo de código).
        3.  Ejecutar las pruebas unitarias.

* `[ ]` **Crear el Pipeline de CD (Despliegue Continuo a Pruebas):**
    * Configurar el pipeline para que se active cuando se hace `merge` a la rama `develop`.
    * **Pasos mínimos del pipeline de CD:**
        1.  Todos los pasos del pipeline de CI.
        2.  Construir la imagen de Docker de la aplicación.
        3.  Publicar la imagen en el registro de contenedores.
        4.  Desplegar la nueva imagen en el ambiente de pruebas.
        5.  Enviar una notificación a un canal de Slack/Teams del equipo de QA.

---
**Resultado Final:**
Al completar este protocolo, el Líder de Proyecto ha creado un **"Proyecto Listo para el Desarrollo"**. Cuando el primer desarrollador tome el ticket `DEV-001`, no perderá tiempo en configuración. Simplemente creará su rama, escribirá el código, hará un `push`, y la máquina de CI/CD que tú has construido se encargará del resto.
