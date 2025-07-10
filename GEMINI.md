# GEMINI.md - Contexto do Projeto Bizly para IA

Este documento serve como um guia de contexto para um assistente de IA (como o Google Gemini) para fornecer ajuda relevante e consistente para o projeto Bizly.

## 1. Resumo do Projeto

*   **Nome:** Bizly
*   **Tipo:** SaaS (Software as a Service)
*   **Objetivo:** Ajudar microempreendedores e autônomos no Brasil com gestão de agendamentos, finanças e clientes.
*   **Plataforma Inicial:** Aplicativo Android, com backend em Go.

## 2. Stack Tecnológica Principal

*   **Backend:**
    *   **Linguagem:** Go (Golang)
    *   **Framework:** Gin Gonic
    *   **ORM:** GORM
    *   **Banco de Dados:** PostgreSQL
*   **Frontend:**
    *   **Framework:** Flutter (Dart)
    *   **Gerenciamento de Estado:** `provider`
*   **Autenticação:**
    *   **Método:** Token JWT enviado via cabeçalho `Authorization: Bearer <token>`.

## 3. Decisões de Arquitetura

### Backend (Clean Architecture)

O backend segue uma estrutura de Clean Architecture com as seguintes camadas:

*   **`entity`**: Structs Go puras representando as entidades de negócio (User, Appointment). **Não contêm tags de banco (gorm) ou JSON.**
*   **`repository`**: **Apenas interfaces** que definem os contratos para operações de banco de dados (ex: `UserRepository`).
*   **`usecase`**: Contém a lógica de negócios central. Depende das interfaces do `repository`, não das implementações concretas.
*   **`infra`**: Implementações concretas.
    *   `persistence/gorm`: Implementação dos repositórios usando GORM. Contém os "GORM Models" (structs com tags `gorm:"..."`) e funções de mapeamento para/de `entity`.
    *   `security`: Lógica de hashing de senhas (bcrypt) e geração/validação de tokens JWT.
*   **`delivery/http`**: Handlers do Gin, DTOs de request/response (com tags `json` e `binding`), e configuração de rotas. Esta camada chama os `usecases`.

### Frontend (Flutter)

A estrutura do Flutter é baseada em funcionalidades:

*   **`features/`**: Cada funcionalidade (auth, appointments) tem suas próprias telas, widgets e lógica.
*   **`shared/`**: Contém código reutilizável entre as features.
    *   `models/`: Modelos de dados do Dart (ex: `UserModel`) com métodos `fromJson`.
    *   `services/`: Serviços que encapsulam lógica, como `AuthService` (lida com estado de login) e `ApiService` (lida com chamadas HTTP).
    *   `widgets/`: Widgets customizados e reutilizáveis.
*   **`core/`**: Configurações centrais do app (tema, constantes, etc.).

### Identificadores (IDs)

*   **Todos os IDs de entidades primárias (Users, Appointments, etc.) são `uuid.UUID`**. Não usamos inteiros auto-incrementais.

## 4. Estrutura de Diretórios de Referência

```
.
├── backend_go/
│   ├── cmd/bizly_api/main.go
│   └── internal/
│       ├── config
│       ├── delivery/http
│       ├── entity
│       ├── infra/persistence/gorm
│       ├── infra/security
│       ├── repository
│       └── usecase
└── frontend_flutter/
    └── bizly_app/
        └── lib/
            ├── core
            ├── features/
            │   ├── auth/screens
            │   └── appointments/screens
            └── shared/
                ├── models
                └── services
```
