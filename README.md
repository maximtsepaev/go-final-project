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

После запуска откройте в браузере http://localhost:7540.