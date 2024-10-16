# Тестовое для BearHead Studio

## Запуск приложения

```bash
docker-compose up --build
```

После запуска api будет доступен по адресу http://localhost:8080/v1/

## Тестирование
Для запуска юнит тестов выполните:

```
make test
```

## Документация API
Swagger-документация доступна по адресу:

http://localhost:8080/swagger/index.html

## API Эндпоинты
### Регистрация пользователя
Запрос:
```bash
curl -X POST \
http://localhost:8080/v1/auth/register \
-H 'Content-Type: application/json' \
-d '{
"username": "yourusername",
"password": "yourpassword"
}'
```

Ответ:
```
HTTP/1.1 201 Created
```

### Вход в систему
Запрос:
```bash
curl -X POST \
http://localhost:8080/v1/auth/login \
-H 'Content-Type: application/json' \
-d '{
"username": "yourusername",
"password": "yourpassword"
}'
```

Ответ:
```
{
"token": "ваш_jwt_токен"
}
```
Сохраните полученный token для последующих запросов.

### Добавление ассета

Запрос:

```bash
curl -X POST \
http://localhost:8080/v1/assets \
-H 'Content-Type: application/json' \
-H 'Authorization: Bearer ваш_jwt_токен' \
-d '{
"name": "Asset Name",
"description": "Asset Description",
"price": 150.00
}'
```

Ответ:
```
HTTP/1.1 201 Created
```

### Получение ассетов пользователя
Запрос:
```bash
curl -X GET \
http://localhost:8080/v1/assets \
-H 'Authorization: Bearer ваш_jwt_токен'
```
Ответ:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "name": "Asset Name",
    "description": "Asset Description",
    "price": 150.00
  },
  ...
]
```

### Удаление ассета
Запрос:
```bash
curl -X DELETE \
http://localhost:8080/v1/assets/1 \
-H 'Authorization: Bearer ваш_jwt_токен'
```

Ответ:
```
HTTP/1.1 200 OK
```

### Покупка ассета

Запрос:

```bash
curl -X POST \
http://localhost:8080/v1/assets/purchase/2 \
-H 'Authorization: Bearer ваш_jwt_токен'
```

Ответ:
```
HTTP/1.1 200 OK
```