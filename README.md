## Описание

Веб-сервер планировщика задач (TODO). Хранит задачи с дедлайном и правилом повторения, поддерживает перенос даты при выполнении повторяющихся задач и удаление обычных. Есть HTTP API и веб-интерфейс для создания, получения, обновления, удаления и отметки задач как выполненных. Данные хранятся в SQLite.

## Задания со звездочкой

Выполнены задания со звездочкой:
- Поддержка `TODO_PORT` и `TODO_DBFILE`.
- Расширенные правила повторения `w` и `m`.
- Поиск задач по строке и по дате через `search` в `/api/tasks`.
- **Аутентификация по `TODO_PASSWORD` с JWT и middleware на API.**
- **Dockerfile и запуск в контейнере.**

## Запуск локально

1) Установите переменные окружения (**опционально**):

- `TODO_PORT` — порт сервера (по умолчанию 7540)
- `TODO_DBFILE` — путь к файлу SQLite (по умолчанию `./scheduler.db`)
- `TODO_PASSWORD` — пароль для авторизации (если не задан, авторизация отключена)

Пример для Windows PowerShell:

```powershell
$env:TODO_PORT="7540"
$env:TODO_DBFILE="./scheduler.db"
$env:TODO_PASSWORD="1234"
```

Пример для bash:

```bash
TODO_PORT=7540 TODO_DBFILE=./scheduler.db TODO_PASSWORD=1234 go run ./...
```

Также можно задать переменные через файл `.env` (bash):

```
TODO_PORT=7540
TODO_DBFILE=./scheduler.db
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
- `DBFile = "../scheduler.db"` — путь к файлу базы для тестов
- `FullNextDate = true` — включение полного набора проверок для `nextdate`
- `Search = true` — включение проверок поиска задач
- `Token = ""` — токен авторизации который сервер возвратил из /api/signin и которое хранится в куке token (оставьте пустым, если `TODO_PASSWORD` не задан)

2) Запустите тесты:

```bash
go test ./tests
```

## Docker

Сборка образа:

```bash
docker build -t todo-app .
```

Запуск без подключения базы (SQLite создается внутри контейнера):

```bash
docker run --rm -p 7540:7540 todo-app
```

Запуск с подключением базы (укажите абсолютный путь к файлу SQLite):

```bash
docker run --rm -p 7540:7540 \
  -v "<ABS_PATH_TO_DB>:/app/scheduler.db" \
  todo-app
```

Опционально можно переопределить переменные окружения:

```bash
docker run --rm -p 7540:7540 \
  -e TODO_PORT=7540 \
  -e TODO_DBFILE=/app/scheduler.db \
  -e TODO_PASSWORD=<value> \
  -v "<ABS_PATH_TO_DB>:/app/scheduler.db" \
  todo-app
```

Можно передать переменные через env-файл:

```bash
docker run --rm -p 7540:7540 \
  --env-file .env \
  -v "<ABS_PATH_TO_DB>:/app/scheduler.db" \
  todo-app
```