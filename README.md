<div align="center">

# 🟩 Habit-Graph `v1.0.0`

**Full-Stack сервис трекинга привычек и визуализации активности в стиле GitHub Activity Heatmap**

[![Release](https://img.shields.io/badge/Release-v1.0.0-10B981?style=for-the-badge&logo=github&logoColor=white)](https://github.com/N1ktarchik/habbit-graph/releases)
[![CI/CD](https://img.shields.io/badge/GitHub_Actions-CI%2FCD-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)](https://github.com/N1ktarchik/habbit-graph/actions)
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![React](https://img.shields.io/badge/React-18.3-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.6-3178C6?style=for-the-badge&logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-3.4-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)](https://tailwindcss.com/)
[![Docker](https://img.shields.io/badge/GHCR_%26_Compose-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

</div>

---

**Habit-Graph** — это современное full-stack веб-приложение для ежедневного учёта времени, потраченного на полезные привычки, и наглядной визуализации прогресса с помощью тепловых карт (Heatmaps) и графиков динамики. В основе проекта — высокопроизводительный REST API на **Go** (Clean Architecture, пул соединений PostgreSQL без ORM, сквозное кэширование в Redis, Graceful Shutdown, покрытие Unit и интеграционными тестами с `testcontainers-go`) и реактивный SPA-интерфейс в эстетике GitHub Dark.

### ✨ Ключевые возможности

* **⏱️ Встроенный таймер с персистентностью:** секундомер активной сессии сохраняет состояние в `localStorage`, не сбрасывается при перезагрузке вкладки и автоматически конвертирует секунды в минуты при завершении.
* **📅 Двойной режим тепловой карты (Месяц / Год):**
  * **Годовой вид (GitHub-style):** 52/53 недели с цветовой градацией интенсивности и всплывающими тултипами.
  * **Месячный вид:** детализированная календарная сетка с прогресс-барами выполнения дневной нормы внутри каждой ячейки.
* **📈 График динамики (Trend Chart):** комбинированный интерактивный SVG/Bar-график за выбранный месяц с расчётом среднего времени в день и подсветкой пикового рекорда.
* **⚡ Атомарное суммирование (Upsert):** несколько сессий за один день автоматически складываются на уровне БД (`ON CONFLICT DO UPDATE`), а кэш в Redis инвалидируется с механизмом фоновых ретраев.
* **🧪 Комплексное тестирование и CI/CD:** покрытие слоёв `domain`, `middleware`, `service` и `transportHTTP` Unit-тестами (93–100%), интеграционные тесты репозиториев с реальными контейнерами PostgreSQL и Redis через `testcontainers-go`, статический анализ `golangci-lint` v2 и автоматическая сборка/публикация релизов в **GitHub Container Registry (`ghcr.io`)** с автообновлением через **Watchtower**.
* **🤖 Инженерный подход:** сочетание классической Backend-разработки и практик AI-Assisted Engineering:
  > **Фокус на Backend + AI:** Вся серверная архитектура (Backend), проектирование БД, кэширование, тестирование, CI/CD-пайплайн и Docker-инфраструктура написаны вручную на Go. Клиентская часть (**Frontend на React + TypeScript + Tailwind CSS**) спроектирована с применением практик **AI-Assisted Development (Prompt Engineering & Iterative AI Prototyping)**, интегрирована с REST API и упакована в production-контейнер с Nginx. Это позволило кратно сократить *Time-to-Market* и получить полноценный UI для демонстрации работы бэкенда.

---

### 🛠 Технологический стек

| Слой | Технология / Библиотека | Версия | Назначение |
| :--- | :--- | :--- | :--- |
| **Backend** | **Go (Golang)** | `1.27.1` | Основной язык серверной части |
| | `net/http` (`http.ServeMux`) | stdlib | Маршрутизация Go 1.22+ с поддержкой `PathValue` |
| | `log/slog` | stdlib | Структурированное JSON-логирование (`debug`, `info`, `warn`, `error`) |
| | `github.com/jackc/pgx/v5` | `v5.11.0` | Высокопроизводительный пул соединений (`pgxpool`) для PostgreSQL |
| | `github.com/redis/go-redis/v9` | `v9.23.0` | Клиент Redis (Pipeline, HSet/HGetAll, Scan-итератор) |
| | `github.com/google/uuid` | `v1.6.0` | Валидация UUID-идентификаторов в запросах |
| | `github.com/joho/godotenv` | `v1.5.1` | Загрузка конфигурации из `.env` |
| **Testing & QA** | `github.com/stretchr/testify` | `v1.12.1` | Табличные Unit-тесты (`assert`, `require`) с детектором гонок (`-race`) |
| | `testcontainers-go` | `v0.45.0` | Интеграционные тесты репозиториев на реальных контейнерах Postgres и Redis |
| | `golangci-lint` | `v2` | Статический анализ кода и проверка форматирования (`goimports`) |
| **Database & Cache** | **PostgreSQL** | `16-alpine` | Основное реляционное хранилище (`habits`, `activities`) |
| | **Redis** | `7-alpine` | Кэширование годовых графов активности (`habit:{id}:graph:{year}`) |
| | **golang-migrate** | `v4.18.1` | Автоматическое применение и откат SQL-миграций |
| **Frontend** | **React + TypeScript** | `18.3.1` / `5.6.3` | Типизированное SPA-приложение |
| | **Vite** | `5.4.10` | Сборщик и dev-сервер |
| | **Tailwind CSS + Lucide** | `3.4.14` / `0.454.0` | Стилизация интерфейса (Dark UI) и иконки |
| **DevOps & CD** | **Docker & Compose** | Multi-stage | Сборка API в легковесный `FROM scratch`, фронтенда в `nginx:1.27-alpine` |
| | **GitHub Actions & GHCR** | CI/CD | Автопрогон тестов и публикация образов в `ghcr.io` по ветке `main` и тегам `v*.*.*` |
| | **Watchtower** | `latest` | Автоматическое обновление запущенных контейнеров `api` и `frontend` из реестра |

---

### 🏗 Архитектура проекта (Clean Architecture)

Проект разделён на независимые слои с инверсией зависимостей через интерфейсы:

```text
Habit-Graph/
├── .github/workflows/
│   └── ci-cd.yml                      # Пайплайн CI/CD (Linter, Unit+Integration Tests, GHCR Push)
├── cmd/
│   └── main.go                        # Точка входа, DI-сборка и Graceful Shutdown
├── internal/
│   ├── core/                          # Ядро системы и общие утилиты
│   │   ├── domain/                    # Доменные модели (Habit, Track), валидация и Unit-тесты
│   │   ├── errors/                    # Кастомные доменные ошибки (ErrorApp)
│   │   ├── logger/                    # Инициализация JSON-логгера log/slog
│   │   ├── repository/                # Фабрики подключений pgxpool и Redis
│   │   └── transport/                 # HTTP-сервер, CORS middleware (+ тесты), JSON-хелперы
│   └── habits/                        # Доменный модуль привычек
│       ├── transportHTTP/             # HTTP-хендлеры, DTO и Unit-тесты (httptest)
│       ├── service/                   # Бизнес-логика (Usecases), кэш-стратегии и Unit-тесты
│       └── repository/
│           ├── postgres/              # SQL-репозиторий (CRUD, Upsert) + интеграционные тесты
│           └── redis/                 # Кэш-репозиторий (Pipeline, Scan) + интеграционные тесты
├── migrations/                        # SQL-миграции (000001_init.up.sql / .down.sql)
├── frontend/                          # SPA-приложение на React + TypeScript + Vite
├── compose.yml                        # Оркестрация (postgres, redis, migrate, api, frontend, watchtower)
├── Dockerfile                         # Multi-stage сборка Go-бинарника в scratch-образ
├── .golangci.yml                      # Конфигурация линтера golangci-lint v2
└── Makefile                           # Автоматизация запуска, миграций, тестов и линтера

```

#### Инженерные решения под капотом:

1. **Cache-Aside + Background Retries:** При запросе графа (`GetGraph`) сервис сначала обращается в Redis (`HGETALL`). При промахе данные читаются из PostgreSQL и асинхронно сохраняются в кэш на 24 часа через `context.WithoutCancel` и Redis Pipeline. При обновлении или удалении привычки кэш сбрасывается; если Redis временно недоступен, запускается фоновая горутина на 5 повторных попыток с интервалом 2 секунды.
2. **Атомарный Upsert в PostgreSQL:** Таблица `activities` использует составной первичный ключ `(habit_id, date)` и запрос `ON CONFLICT (habit_id, date) DO UPDATE SET duration_minutes = activities.duration_minutes + EXCLUDED.duration_minutes`, что исключает состояние гонки при параллельных запросах.
3. **Ультра-легковесный Docker-образ:** Бэкенд компилируется статически (`CGO_ENABLED=0`) с флагами `-ldflags="-w -s"`, после чего упаковывается в пустой базовый образ `FROM scratch` вместе с TLS-сертификатами и базой часовых поясов.
4. **Непрерывная доставка (CD) через GHCR + Watchtower:** При мерже в `main` или создании тега релиза (`v1.0.0`) GitHub Actions прогоняет все проверки и публикует готовые образы в `ghcr.io`. Контейнер `watchtower` отслеживает обновления помеченных лейблами сервисов (`api` и `frontend`) и бесшовно перезапускает их без участия разработчика.

---

### 🚀 Быстрый старт

#### Предварительные требования

* **Docker** и **Docker Compose**
* **Make** *(опционально, для удобного вызова команд)*

#### Запуск через Docker Compose

1. **Клонируйте репозиторий:**
```bash
git clone [https://github.com/N1ktarchik/habbit-graph.git](https://github.com/N1ktarchik/habbit-graph.git)
cd habbit-graph

```


2. **Создайте файл окружения:**
Скопируйте `.env.example` в `.env` (для запуска внутри Docker-сети убедитесь, что `POSTGRES_HOST=postgres`, `POSTGRES_PORT=5432` и `REDIS_HOST=redis`):
```bash
cp .env.example .env

```


3. **Запустите инфраструктуру:**
```bash
make docker-up
# или без make:
docker compose up -d

```


*Docker Compose скачает готовые образы из `ghcr.io`, поднимет `postgres` и `redis`, дождётся прохождения `healthcheck`, автоматически применит миграции через `migrate`, после чего запустит `api`, `frontend` и `watchtower`.*
4. **Откройте приложение:**
* 🌐 **Веб-интерфейс (Frontend):** [http://localhost:3000](http://localhost:3000)
* ⚙️ **Backend REST API:** [http://localhost:8080](http://localhost:8080)
* 🐘 **PostgreSQL (снаружи контейнера):** `localhost:5433` (БД: `habit-graph_app`, пользователь: `postgres`)
* 🟥 **Redis (снаружи контейнера):** `localhost:6379`



---

### 🧪 Тестирование и контроль качества

Проект покрыт модульными и интеграционными тестами с проверкой состояния гонок (`-race`):

| Пакет / Слой | Тип тестов | Покрытие (`statements`) | Что проверяется |
| --- | --- | --- | --- |
| `internal/core/domain` | Unit | **100.0%** | Валидация длины рун (UTF-8/кириллица), дат и диапазона годов |
| `internal/core/transport/middleware` | Unit | **100.0%** | CORS-заголовки и перехват preflight `OPTIONS` запросов |
| `internal/habits/service` | Unit + Mocks | **96.2%** | Бизнес-логика, Cache-Aside, фоновые горутины ретраев Redis |
| `internal/habits/transportHTTP` | Unit (`httptest`) | **93.8%** | Маршрутизация, валидация UUID/JSON, HTTP-коды и маппинг DTO |
| `internal/habits/repository/redis` | Integration (`testcontainers`) | **91.5%** | Реальный контейнер `redis:7-alpine`, Pipeline, TTL, `SCAN`-удаление |
| `internal/habits/repository/postgres` | Integration (`testcontainers`) | **85.9%** | Реальный контейнер `postgres:16-alpine`, миграции, `COALESCE`, `ON CONFLICT`, каскадное удаление |

Запуск тестов локально:

```bash
# Только быстрые Unit-тесты (без Docker)
make unit-test

# Интеграционные тесты репозиториев с автоподнятием контейнеров через Testcontainers
make test-integration

# Полный прогон всех тестов (Unit + Integration)
make test-all

# Статический анализ линтером golangci-lint v2
make lint

```

---

### 🛠 Команды Makefile

В проект включён `Makefile` для автоматизации рутинных операций:

| Команда | Описание |
| --- | --- |
| `make help` | Показать список всех доступных команд с подсказками |
| `make docker-up` | Собрать и поднять всю инфраструктуру в фоновом режиме (`docker compose up -d --build`) |
| `make docker-down` | Остановить и удалить контейнеры (`docker compose down`) |
| `make docker-logs` | Просмотр live-логов всех сервисов (`docker compose logs -f`) |
| `make migrate-up` | Накатить миграции через контейнер `migrate` |
| `make migrate-down` | Откатить последнюю применённую миграцию |
| `make migrate-create NAME=...` | Создать новую пару SQL-миграций в папке `./migrations` |
| `make db-reset` | **Полный сброс:** удаляет Docker-тома (`down -v`), заново поднимает БД и накатывает миграции |
| `make run-back` | Запустить Go-бэкенд локально без Docker (`go run ./cmd/main.go`) |
| `make run-front` | Установить зависимости и запустить Vite dev-сервер на порту `3000` |
| `make unit-test` | Запустить быстрые Unit-тесты с детектором гонок (`-race`) без Docker |
| `make test-integration` | Запустить интеграционные тесты PostgreSQL и Redis через `testcontainers-go` |
| `make test-all` | Запустить все тесты проекта (Unit + Integration) с подсчётом покрытия |
| `make lint` | Запустить проверку кода линтером `golangci-lint` |

---

### 📡 Документация REST API

Все запросы и ответы используют формат `application/json`.
Ограничения валидации: `title` — от 3 до 120 символов, `description` — до 200 символов, `id` — валидный UUID, год — `>= 2026` и не больше текущего, дата активности — формат `YYYY-MM-DD`.

| Метод | Эндпоинт | Описание | Тело запроса (JSON) | Успешный ответ |
| --- | --- | --- | --- | --- |
| `POST` | `/api/habits` | Создать новую привычку | `{"title": "Чтение книг", "description": "По 30 минут вечером"}` | `201 Created` и объект `HabitResponseDTO` |
| `GET` | `/api/habits` | Получить список всех привычек | — | `200 OK` и массив `[]HabitResponseDTO` |
| `PATCH` | `/api/habits/{id}` | Частично обновить привычку | `{"title": "Новое название", "description": "..."}` | `200 OK` и обновлённый `HabitResponseDTO` |
| `DELETE` | `/api/habits/{id}` | Удалить привычку (каскадно удаляет историю в БД и кэш в Redis) | — | `204 No Content` |
| `POST` | `/api/habits/{id}/activities` *(алиас: `/api/habits/{id}/track`)* | Записать / прибавить время активности за день | `{"date": "2026-11-16", "value": 45}` *(в минутах)* | `200 OK` и `{"status": "ok"}` |
| `GET` | `/api/habits/{id}/graph?year=2026` | Получить карту активности за год (из Redis или БД) | — | `200 OK` и словарь `{"2026-11-16": 45, "2026-11-17": 90}` |

#### Пример ответа `HabitResponseDTO`:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Программирование на Go",
  "description": "Разработка микросервисов и pet-проектов",
  "created_at": "2026-11-16T15:04:05Z"
}

```

#### Пример ответа при ошибке (`ErrorApp`):

```json
{
  "code": 400,
  "message": "invalid title"
}

```

---

### 👤 Автор

**Никита Клеймёнов** — Backend Developer (Go)

GitHub: [@N1ktarchik](https://github.com/N1ktarchik)

