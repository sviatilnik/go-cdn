# Генерация Swagger документации

Проект использует [swaggo/swag](https://github.com/swaggo/swag) для генерации Swagger/OpenAPI документации.

## Генерация документации

Из корневой директории проекта:

```bash
swag init -g cmd/server/main.go -o ./docs
```

Это создаст файлы:
- `docs/swagger.json`
- `docs/swagger.yaml`
- `docs/docs.go`

## Формат аннотаций

Проект использует стандартные аннотации swaggo:

- `@title` - заголовок API
- `@version` - версия API
- `@Summary` - краткое описание эндпоинта
- `@Description` - подробное описание
- `@Tags` - группировка эндпоинтов
- `@Accept` - принимаемый Content-Type
- `@Produce` - возвращаемый Content-Type
- `@Param` - параметры запроса
- `@Success` - успешный ответ
- `@Failure` - ошибки
- `@Router` - путь и метод

Подробнее: https://github.com/swaggo/swag#declarative-comments-format
