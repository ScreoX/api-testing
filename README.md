# Автотесты REST и GraphQL API на Go

Тесты используют [Playground REST API](https://playground.nileslabs.com/api/v1) и [Postman Echo GraphQL](https://graphql.postman-echo.com/graphql). Нужен Go 1.22+

## Сценарии

- `TestRESTPostPatchGet`: POST `/posts` -> извлечение `id` -> PATCH `/posts/{id}` -> GET `/posts/{id}`. Проверяются код ответа, `id`, новый `title` и сохранность `body` и `userId`. Каждый запуск использует отдельную песочницу `_sandbox`
- `TestBearerAuthentication`: POST `/auth/login` -> извлечение `access_token` -> GET `/auth/me` с Bearer-токеном. Проверяются данные пользователя. `TestRejectsInvalidBearerToken` проверяет отказ при неверном токене
- `TestGraphQLQuery` проверяет `data.hello`, `TestGraphQLMutation` - `id`, `name` и `age` созданного человека. Оба теста проверяют HTTP 200 и отсутствие `errors`. `TestGraphQLValidationError` - **дополнительный негативный** сценарий для неизвестного поля
- Subscription `Greetings { greetings }` находится в [Postman collection](postman/GraphQL.postman_collection.json)

Тесты лежат в `tests/*_test.go`. Общие HTTP-настройки, типы и логирование - в `tests/*_helper.go`. При падении теста выводятся запрос и ответ. Пароль и токен в логах скрыты

## Запуск

Из корня проекта:

```bash
go test -count=1 ./...
```

По умолчанию используются `Bret / Password@123`. Указанный в задании пароль `password123` сейчас возвращает HTTP 401. При изменении данных можно задать `API_USERNAME` и `API_PASSWORD` через переменные окружения

## Subscription в Postman

Импортируйте collection, откройте запрос `greetings` и нажмите Send. Убедитесь, что с заголовком `Accept: text/event-stream` приходят события `next` с `data.greetings`, затем `complete`
