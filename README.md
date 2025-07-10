# Bizly - SaaS de Gestão para Microempreendedores

Bizly é um Sistema como Serviço (SaaS) focado em ajudar microempreendedores e autônomos no Brasil a gerenciar seus negócios de forma simples e eficiente. A plataforma visa resolver problemas comuns de gestão financeira, agendamento de clientes e vendas.

## Sobre o Projeto

O objetivo do Bizly é oferecer uma ferramenta "tudo em um" acessível e intuitiva para setores como estética, barbearias, freelancers, MEIs, etc. O aplicativo centraliza agenda, pagamentos, emissão de recibos e comunicação com o cliente.

## Tecnologias Utilizadas

**Backend:**
*   **Linguagem:** Go (Golang)
*   **Framework Web:** Gin Gonic
*   **ORM:** GORM
*   **Banco de Dados:** PostgreSQL
*   **Autenticação:** JWT (JSON Web Tokens)
*   **Arquitetura:** Baseada em princípios de Clean Architecture (Entidades -> Casos de Uso -> Repositórios -> Delivery/Infra)
*   **IDs:** UUID para todas as entidades principais.

**Frontend:**
*   **Framework:** Flutter (Dart)
*   **Gerenciamento de Estado:** Provider
*   **Armazenamento Seguro:** `flutter_secure_storage` para JWT.
*   **Plataforma Alvo:** Android (inicialmente), com potencial para iOS e Web.

## Estrutura do Projeto

O projeto é um monorepo com duas pastas principais:

```
SaaS_microbusiness_project/
├── backend_go/          # Backend em Go
│   ├── cmd/bizly_api/
│   ├── internal/
│   │   ├── config/
│   │   ├── delivery/ (HTTP Handlers, Middleware, Router)
│   │   ├── entity/
│   │   ├── infra/ (GORM, Security)
│   │   ├── repository/ (Interfaces de Repositório)
│   │   └── usecase/ (Lógica de Negócios)
│   └── ...
├── frontend_flutter/    # Frontend em Flutter
│   └── bizly_app/
│       ├── lib/
│       │   ├── core/
│       │   ├── features/ (Auth, Appointments, Home)
│       │   └── shared/ (Models, Services, Widgets)
│       └── ...
└── README.md
```

## Status do Projeto (O que já foi feito)

### Backend
- [x] Estrutura do projeto Go com base em Clean Architecture.
- [x] Configuração do servidor web com Gin Gonic.
- [x] Conexão com banco de dados PostgreSQL usando GORM.
- [x] Sistema de configuração flexível via variáveis de ambiente (`.env`).
- [x] Implementação de CORS para permitir comunicação com o frontend.
- [x] **Gerenciamento de Usuários:**
  - [x] Cadastro de usuário com hashing de senha (bcrypt).
  - [x] Login de usuário com verificação de senha.
- [x] **Autenticação e Autorização:**
  - [x] Geração de token JWT após login bem-sucedido.
  - [x] Middleware de autenticação para proteger rotas.
- [x] **Gerenciamento de Agendamentos (CRUD Básico):**
  - [x] Criação, listagem, busca por ID, atualização e cancelamento de agendamentos.
  - [x] Lógica de permissão (usuário só pode gerenciar seus próprios agendamentos).
- [x] Migração de IDs inteiros para UUIDs em todas as entidades e camadas.

### Frontend
- [x] Configuração inicial do projeto Flutter com estrutura de pastas organizada.
- [x] Implementação de gerenciamento de estado com Provider para autenticação.
- [x] Criação de `ApiService` para comunicação com o backend Go.
- [x] Armazenamento seguro de token JWT no dispositivo (`flutter_secure_storage`).
- [x] **Fluxo de Autenticação:**
  - [x] Telas de Login e Home (placeholder) funcionais.
  - [x] Lógica de login/logout e redirecionamento de tela.
- [x] **Conectividade:**
  - [x] Comunicação bem-sucedida com o backend tanto em ambiente Web quanto no Emulador Android.

## Próximos Passos (Roadmap)

### Frontend (Prioridade Alta)
- [ ] Implementar a tela de Cadastro de Usuário.
- [ ] Construir as telas do fluxo de Agendamentos (Listagem, Detalhes, Criação, Edição).
- [ ] Refinar a UI/UX para uma experiência mais limpa e profissional.
- [ ] Melhorar o fluxo de "auto-login" para buscar dados frescos do usuário ao iniciar o app.
- [ ] Implementar feedback visual para o usuário (loading spinners, snackbars de erro/sucesso).

### Backend (Refinamentos)
- [ ] Implementar lógica de negócio para **validação de conflito de horários**.
- [ ] Implementar **paginação** para listagens (ex: agendamentos).
- [ ] Refinar o tratamento de erros com tipos de erro customizados.
- [ ] Adicionar testes unitários e de integração.
- [ ] Gerar documentação da API com Swagger/OpenAPI.

### Novas Funcionalidades (Médio/Longo Prazo)
- [ ] Módulo de Gestão Financeira (registro de entradas e saídas).
- [ ] Emissão de Recibos em PDF.
- [ ] Integração com sistema de pagamento (PIX).
- [ ] Sistema de notificações (via Push Notification ou WhatsApp).
- [ ] Dashboard web para uma visão geral do negócio.

## Como Executar

**Pré-requisitos:**
*   Go (versão 1.18+)
*   Flutter (versão 3.x+)
*   PostgreSQL
*   Android Studio (para o SDK Android e emuladores)

**1. Backend:**
   - Navegue até `backend_go/`.
   - Crie um arquivo `.env` a partir do `.env.example` e configure suas variáveis (principalmente a senha do banco de dados).
   - Certifique-se de que o banco de dados (`bizly_db`) foi criado no PostgreSQL.
   - Habilite a extensão `uuid-ossp` no seu banco: `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
   - Navegue até `backend_go/cmd/bizly_api/`.
   - Execute: `go run main.go`
   - O servidor estará rodando em `http://localhost:8080`.

**2. Frontend:**
   - Navegue até `frontend_flutter/bizly_app/`.
   - Execute: `flutter pub get`
   - Inicie um emulador Android ou conecte um dispositivo.
   - Execute: `flutter run`
   - O app irá compilar e instalar no dispositivo/emulador.