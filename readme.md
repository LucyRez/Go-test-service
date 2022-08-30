# Тестовый сервис на Go
(порт для старта кафки :9092 , для keycloak :8080, каждый сервис запускается отдельно на своём порту)

## Service-producer
### Ендпоинты

### localhost:9090/login
- Post запрос с телом для аутентификации через keycloak
- Пример тела запроса: 
```
{
    "username": "user",
    "password": "password"
}

```
### localhost:9090/register
- Только POST запрос
- Доступен только для пользователей, прошедших авторизацию с ролью ADMIN
- Пример тела запроса: 

```
{
    "firstName": "",
    "lastName": "",
    "email": "",
    "enabled": true,
    "username": "user2",
    "password": "password2"
}
```

### localhost:9090/entity/
- Доступен только после аутентификации
- Get запрос - получение списка всех объектов (доступно для всех аутентифицированных пользователей)
- Post запрос с телом - добавление нового объекта в список (только для пользователей с ролью ADMIN)
### localhost:9090/submit
- Get запрос с телом - отправка сообщения на другой сервис (сообщение содержит объект, отправленный в теле запроса)
- Пример тела запроса: 
```
{
    "id": 1,
    "name": "Second Entity"
}

```

## Service-receiver
### Ендпоинты
### localhost:9091/received
- Get запрос - получение списка всех объектов, доставленных сервисом-продюсером
