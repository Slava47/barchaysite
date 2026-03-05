# Barchay API Server

Go HTTP-сервер для чайного бара «Ли Бо». Раздаёт меню и квиз, принимает загрузку новых позиций от администратора.

## Быстрый старт

### 1. Сборка

```bash
cd api
go build -o barchay-api .
```

### 2. Запуск (локально для разработки)

```bash
# Сгенерируйте секретный ключ:
export ADMIN_API_KEY=$(openssl rand -hex 32)

# Разрешите CORS для локальной разработки:
export ALLOWED_ORIGINS="*"

# Запустите сервер:
cd api
go run .
```

Сервер запустится на `http://localhost:8080`.

### 3. Проверка

```bash
# Проверка работы:
curl http://localhost:8080/health
# → {"status":"ok"}

# Получить меню:
curl http://localhost:8080/menu

# Получить квиз:
curl http://localhost:8080/quiz
```

### 4. Добавление позиции меню

```bash
curl -X POST http://localhost:8080/admin/items-with-images \
  -H "X-API-Key: $ADMIN_API_KEY" \
  -F 'item={"id":"mojito-1","category":"cold","name":"Мохито","price":"420","description":"Классический мятный коктейль","tags":["освежающий","холодный"]};type=application/json' \
  -F "images=@/path/to/photo.jpg;type=image/jpeg"
```

## Подключение frontend к API

### Вариант 1: Через `window.__CONFIG__` (рекомендуется)

В файле `index.html`, **до** подключения скриптов, раскомментируйте блок конфигурации и укажите адрес вашего API:

```html
<script>
  window.__CONFIG__ = {
    API_BASE: 'http://localhost:8080'          <!-- для разработки -->
    // API_BASE: 'https://apibarchay.rudenko.tech' <!-- для продакшена -->
  };
</script>
```

### Вариант 2: Прямая правка `js/app.js`

Найдите строку:
```javascript
const API_BASE = (window.__CONFIG__ && window.__CONFIG__.API_BASE) || '';
```
И замените пустую строку на URL вашего API:
```javascript
const API_BASE = (window.__CONFIG__ && window.__CONFIG__.API_BASE) || 'http://localhost:8080';
```

> Если `API_BASE` пуст или не задан — фронтенд использует локальные данные из `js/menu-data.js` и `js/quiz-data.js`.

## Переменные окружения

| Переменная | Обязательна | По умолчанию | Описание |
|---|---|---|---|
| `ADMIN_API_KEY` | **да** | *(нет)* | Ключ авторизации для `/admin/` эндпоинтов |
| `PORT` | нет | `8080` | Порт сервера |
| `DATA_DIR` | нет | `data` | Путь к директории с `menu.json` |
| `UPLOADS_DIR` | нет | `uploads` | Путь к директории с изображениями |
| `PUBLIC_URL` | нет | `http://localhost:PORT` | Базовый URL для формирования ссылок на загруженные изображения |
| `ALLOWED_ORIGINS` | нет | `https://barchay.rudenko.tech` | Разрешённые CORS-источники (через запятую). Используйте `*` для разработки |

## Эндпоинты API

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/health` | Проверка работоспособности |
| `GET` | `/menu` | Полное меню (категории + позиции) |
| `GET` | `/quiz` | Данные для квиза (вопросы) |
| `POST` | `/admin/items-with-images` | Добавление позиции с изображениями (требует `X-API-Key`) |
| `GET` | `/uploads/{filename}` | Статические файлы изображений |

## Деплой на сервер

Подробная инструкция: [`docs/API_RU.md`](../docs/API_RU.md) (раздел 7).

Краткие шаги:

1. Соберите бинарник: `go build -o barchay-api .`
2. Скопируйте `barchay-api` на сервер в `/opt/barchay-api/`
3. Создайте директории `data/` и `uploads/`
4. Настройте env-файл с `ADMIN_API_KEY` и `PUBLIC_URL`
5. Создайте systemd unit и запустите сервис
6. Настройте nginx/Caddy как обратный прокси с TLS

## Тесты

```bash
cd api
go test -v ./...
```
