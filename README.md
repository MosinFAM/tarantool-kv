## kv-хранилище с использованием Tarantool


#### `make build`

Создаёт приложение в bin/redditclone.

#### `make test`

Запускает тесты.

#### `make lint`

Запускает линтер.

## Запуск контейнера

```bash
docker compose -f build/docker-compose.yml up -d --build
```

## API

- POST /kv body: {key: "test", "value": {SOME ARBITRARY JSON}} 

- PUT kv/{id} body: {"value": {SOME ARBITRARY JSON}}

- GET kv/{id} 

- DELETE kv/{id}


- POST  возвращает 409 если ключ уже существует, 

- POST, PUT возвращают 400 если боди некорректное

- PUT, GET, DELETE возвращает 404 если такого ключа нет

- все операции логируются


## примеры запросов

Получение

```bash
curl -X GET "http://localhost:8080/kv/test"   
```

Отправка 

```bash
curl -X POST "http://localhost:8080/kv" \                              
     -H "Content-Type: application/json" \
     -d '{"key": "test", "value": {"1": "1"}}'
```

Изменение 

```bash
curl -X PUT "http://localhost:8080/kv/test" \                          
     -H "Content-Type: application/json" \
     -d '{"value": {"2": "2"}}'

```

Удаление 

```bash
curl -X DELETE "http://localhost:8080/kv/test"     
```