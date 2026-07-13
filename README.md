## Описание проекта

REST API сервер планировщика задач (TODO). Хранит задачи с дедлайном и правилом повторения, поддерживает перенос даты при выполнении повторяющихся задач и удаление, если правила повторения нет. API позволяет создавать, получать, обновлять, удалять задачи и отмечать их выполненными. Данные хранятся в PostgreSQL.

## Задания повышенной трудности

Выполнены задания со звездочкой:
- Поддержка переменных окружения `TODO_PORT` и `DATABASE_URL` (DSN для PostgreSQL).
- Расширенные правила повторения `w` и `m`.
- Поиск задач по строке и по дате через `search` в `/api/tasks`.
- **Аутентификация по `TODO_PASSWORD` с JWT и middleware на API.**
- **Dockerfile, Docker Compose и запуск в контейнерах с базой PostgreSQL.**

## Запуск локально

1) Установите переменные окружения (**опционально**):

- `TODO_PORT` — порт сервера (по умолчанию 7540)
- `DATABASE_URL` — DSN подключения к PostgreSQL (например, `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`). **Обязательно для заполнения!** При пустом значении сервер завершит работу с ошибкой.
- `TODO_PASSWORD` — пароль для авторизации (если не задан, авторизация отключена)

Пример для Windows PowerShell:

```powershell
$env:TODO_PORT="7540"
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
$env:TODO_PASSWORD="1234"
```

Пример для bash:

```bash
TODO_PORT=7540 DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" TODO_PASSWORD=1234 go run ./...
```

Также можно задать переменные через файл `.env` (bash):

```
TODO_PORT=7540
DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
TODO_PASSWORD=1234
```

```bash
set -a && source .env && set +a
go run ./...
```

2) Запустите сервер:

```bash
go run ./...
```

3) Откройте в браузере http://localhost:7540

## Запуск тестов

1) Проверьте параметры в [tests/settings.go](tests/settings.go):

- `Port = 7540` — порт запущенного сервера
- `DBFile = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"` — DSN базы для тестов
- `FullNextDate = true` — включение полного набора проверок для `nextdate`
- `Search = true` — включение проверок поиска задач
- `Token = ""` — токен авторизации который сервер возвратил из /api/signin (оставьте пустым, если `TODO_PASSWORD` не задан)

2) Запустите тесты (перед запуском тестов задайте переменную окружения `DATABASE_URL`):

```bash
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
go test ./tests
```

## Запуск в Docker

### Вариант 1. Запуск через Dockerfile (требуется внешняя БД)

Если у вас уже есть запущенная база данных PostgreSQL (на хосте или в облаке), вы можете собрать и запустить только контейнер с приложением:

1) Соберите образ:

```bash
docker build -t todo-app .
```

2) Запустите контейнер, передав DSN внешней базы данных в переменную `DATABASE_URL`:

```bash
docker run --rm -p 7540:7540 -e DATABASE_URL="postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable" todo-app
```

*Примечание: `host.docker.internal` используется в Windows/macOS для обращения к localhost хостовой машины из контейнера.*

### Вариант 2. Запуск через Docker Compose (рекомендуется)

Если у вас нет внешней базы данных, вы можете запустить приложение вместе с PostgreSQL одной командой:

```bash
docker-compose up --build -d
```

Compose автоматически:
1. Запустит контейнер `todo_db` с PostgreSQL.
2. Проверит его готовность с помощью `healthcheck`.
3. Соберет образ приложения и запустит контейнер `todo_app`, подключив его к базе данных.

Открыть приложение: http://localhost:7540

Остановить работу и очистить ресурсы:

```bash
docker-compose down -v
```

## Документация (godoc)

Запуск сервера документации:

```bash
godoc -http=:6060
```

Откройте в браузере http://localhost:6060/pkg/github.com/maximtsepaev/go-final-project/.

Пакеты из `internal/` не показываются в общем списке. Их можно открыть по прямым ссылкам:

- http://localhost:6060/pkg/github.com/maximtsepaev/go-final-project/internal/api/
- http://localhost:6060/pkg/github.com/maximtsepaev/go-final-project/internal/db/
- http://localhost:6060/pkg/github.com/maximtsepaev/go-final-project/internal/server/
