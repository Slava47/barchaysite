# Документация: API и деплой Barchaysite

> Язык: **Русский**  
> Репозиторий: [Slava47/barchaysite](https://github.com/Slava47/barchaysite)  
> Сайт (frontend): **https://barchay.rudenko.tech**  
> API: **https://apibarchay.rudenko.tech**

---

## Содержание

1. [Архитектура](#1-архитектура)
2. [Контракт API](#2-контракт-api)
   - [GET /menu](#get-menu)
   - [GET /quiz](#get-quiz)
   - [GET /health](#get-health)
   - [POST /admin/items-with-images](#post-adminitems-with-images)
3. [Формат данных](#3-формат-данных)
4. [CORS](#4-cors)
5. [Service Worker и кеширование](#5-service-worker-и-кеширование)
6. [Конфигурация frontend](#6-конфигурация-frontend)
7. [Деплой сайта](#7-деплой-сайта)
   - [Статический frontend](#71-статический-frontend)
   - [DNS](#72-dns)
   - [TLS и обратный прокси (nginx)](#73-tls-и-обратный-прокси-nginx)
   - [Структура директорий API](#74-структура-директорий-api)
   - [Переменные окружения и секреты](#75-переменные-окружения-и-секреты)

---

## 1. Архитектура

```
Браузер пользователя
    │
    ▼
┌─────────────────────────────────────────────────────┐
│  Static frontend (barchay.rudenko.tech)             │
│  index.html + css/ + js/ + icons/ + sw.js           │
│  Хостинг: любой статический хостинг / nginx         │
└────────────────────────┬────────────────────────────┘
                         │  GET /menu, GET /quiz
                         ▼
┌─────────────────────────────────────────────────────┐
│  Go API  (apibarchay.rudenko.tech)                  │
│  ┌─────────────┐  ┌──────────────────────────────┐  │
│  │  data/      │  │  uploads/                    │  │
│  │  menu.json  │  │  <uuid>.jpg, <uuid>.png, ... │  │
│  └─────────────┘  └──────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                         ▲
                         │  POST /admin/items-with-images
                    Администратор (curl / CLI)
```

**Компоненты:**

| Компонент | Описание |
|-----------|----------|
| **Frontend** | Статический SPA (HTML/CSS/JS). При наличии `API_BASE` загружает меню и квиз из API; иначе использует локальные `MENU_DATA` / `QUIZ_DATA`. |
| **Go API** | HTTP-сервер на Go. Раздаёт меню и квиз, принимает загрузку товаров от администратора. |
| **menu.json** | Единственное хранилище данных меню. Читается при старте и обновляется атомарно (запись в tmp-файл + переименование). |
| **uploads/** | Статические файлы изображений. Раздаются по пути `/uploads/<имя-файла>`. |

---

## 2. Контракт API

Базовый URL: `https://apibarchay.rudenko.tech`

### GET /menu

Возвращает полное меню.

**Запрос:**
```http
GET /menu HTTP/1.1
Host: apibarchay.rudenko.tech
```

**Ответ `200 OK`:**
```json
{
  "categories": [
    { "id": "tea", "name": "Чаи", "nameZh": "茶" }
  ],
  "items": [
    {
      "id": "lapsang-1",
      "category": "tea",
      "name": "Лапсанг Сушонг",
      "nameZh": "正山小种",
      "price": "350",
      "description": "Дымный красный чай из провинции Фуцзянь.",
      "fullDescription": "Полное описание...",
      "image": "https://apibarchay.rudenko.tech/uploads/lapsang.jpg",
      "images": [
        "https://apibarchay.rudenko.tech/uploads/lapsang.jpg",
        "https://apibarchay.rudenko.tech/uploads/lapsang2.jpg"
      ],
      "tags": ["красный", "дымный", "горячий"]
    }
  ]
}
```

---

### GET /quiz

Возвращает данные для квиза (вопросы и алгоритм рекомендаций).

**Запрос:**
```http
GET /quiz HTTP/1.1
Host: apibarchay.rudenko.tech
```

**Ответ `200 OK`:** структура `QUIZ_DATA` (вопросы + метки для алгоритма).

---

### GET /health

*(опционально)* Проверка работоспособности сервиса.

**Ответ `200 OK`:**
```json
{ "status": "ok" }
```

---

### POST /admin/items-with-images

Защищённый эндпоинт для добавления нового пункта меню вместе с изображениями.

#### Авторизация

Заголовок `X-API-Key` со значением из переменной окружения `ADMIN_API_KEY`.

```http
X-API-Key: <ваш-секретный-ключ>
```

При отсутствии или неверном ключе — `401 Unauthorized`.

#### Формат запроса

`Content-Type: multipart/form-data`

| Часть (part) | Обязательна | Content-Type | Описание |
|---|---|---|---|
| `item` | да | `application/json` | JSON-объект нового пункта меню (без полей `image`/`images` — они заполняются сервером). |
| `images` | нет | `image/jpeg`, `image/png`, `image/webp`, `image/gif` | Изображения (1..N файлов). Каждый файл — отдельная часть с именем `images`. |

> **Лимиты:** максимальный размер одного файла — 10 МБ; максимальный суммарный размер запроса — 50 МБ (настраивается в прокси).

#### Алгоритм обработки на сервере

1. Проверить заголовок `X-API-Key`.
2. Прочитать часть `item` (JSON), валидировать обязательные поля.
3. Для каждого файла из частей `images`:
   - Проверить MIME-тип (только `image/*`).
   - Сгенерировать имя файла: `<uuid>.<ext>` (UUID v4, расширение из Content-Type).
   - Сохранить в директорию `uploads/`.
   - Сформировать публичный URL: `https://apibarchay.rudenko.tech/uploads/<uuid>.<ext>`.
4. Установить `item.image` = первый URL, `item.images` = все URL.
5. Загрузить `data/menu.json`, добавить item в `menu.items`.
6. Атомарно сохранить `menu.json` (запись в `data/menu.json.tmp` → `rename`).
7. Вернуть `201 Created` с добавленным объектом.

#### Пример curl

**Добавление коктейля с двумя изображениями:**

```bash
curl -X POST https://apibarchay.rudenko.tech/admin/items-with-images \
  -H "X-API-Key: YOUR_ADMIN_API_KEY" \
  -F 'item={"id":"mojito-1","category":"cocktails","name":"Мохито","nameZh":"莫希托","price":"420","description":"Классический мятный коктейль","fullDescription":"Ром, лайм, мята, содовая","tags":["освежающий","холодный","алкогольный"]};type=application/json' \
  -F "images=@/path/to/mojito-main.jpg;type=image/jpeg" \
  -F "images=@/path/to/mojito-detail.jpg;type=image/jpeg"
```

**Добавление позиции без изображений:**

```bash
curl -X POST https://apibarchay.rudenko.tech/admin/items-with-images \
  -H "X-API-Key: YOUR_ADMIN_API_KEY" \
  -F 'item={"id":"green-tea-1","category":"tea","name":"Лунцзин","nameZh":"龙井","price":"280","description":"Классический зелёный чай","tags":["зелёный","горячий"]};type=application/json'
```

#### Коды ответов

| Код | Описание |
|-----|----------|
| `201 Created` | Позиция успешно добавлена. Тело: добавленный объект Item. |
| `400 Bad Request` | Неверный формат запроса, отсутствуют обязательные поля или недопустимый тип файла. |
| `401 Unauthorized` | Отсутствует или неверный `X-API-Key`. |
| `409 Conflict` | Позиция с таким `id` уже существует в меню. |
| `500 Internal Server Error` | Ошибка сохранения файла или `menu.json`. |

---

## 3. Формат данных

### Объект Item (минимальные обязательные поля)

```json
{
  "id":          "уникальный-идентификатор",
  "category":    "id-категории-из-menu.categories",
  "name":        "Название на русском",
  "price":       "350"
}
```

### Рекомендуемые дополнительные поля

```json
{
  "nameZh":          "中文名称",
  "description":     "Краткое описание (до 80 символов)",
  "fullDescription": "Полное описание для детальной карточки",
  "image":           "https://apibarchay.rudenko.tech/uploads/<uuid>.jpg",
  "images":          ["https://apibarchay.rudenko.tech/uploads/<uuid>.jpg"],
  "tags":            ["тег1", "тег2"]
}
```

> **Теги** используются алгоритмом рекомендаций квиза. Рекомендуется использовать: вкусовые (`сладкий`, `кислый`, `крепкий`, `пряный`), температуру (`тёплый`, `холодный`), алкоголь (`алкогольный`, `безалкогольный`), настроение (`мягкий`, `бодрящий`).

### Структура menu.json

```json
{
  "categories": [
    { "id": "tea",       "name": "Чаи",     "nameZh": "茶"  },
    { "id": "cocktails", "name": "Коктейли","nameZh": "鸡尾酒" }
  ],
  "items": [ /* массив объектов Item */ ]
}
```

---

## 4. CORS

### Разрешённые источники

Для публичных эндпоинтов `GET /menu` и `GET /quiz` необходимо разрешить только frontend-домен:

```
Access-Control-Allow-Origin: https://barchay.rudenko.tech
Access-Control-Allow-Methods: GET, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

Пример конфигурации CORS в Go (middleware):

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if origin == "https://barchay.rudenko.tech" {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        }
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Почему CORS не нужен для /admin/items-with-images

Эндпоинт `/admin/items-with-images` предназначен для использования только через `curl` или CLI-инструменты. Такие запросы не являются браузерными cross-origin запросами и не подпадают под политику CORS. Настраивать `Access-Control-Allow-Origin` для этого эндпоинта не требуется — это было бы излишним расширением поверхности атаки.

---

## 5. Service Worker и кеширование

### Как работает sw.js

Service Worker (`sw.js`) кеширует статические файлы сайта (HTML, CSS, JS, иконки). При каждом запросе сначала пытается получить свежий ответ из сети, при неудаче отдаёт из кеша (стратегия "Network first with cache fallback").

Версия кеша: `libo-v1` (константа `CACHE_NAME` в `sw.js`).

### Влияние на обновление данных меню

**Важно:** API-запросы к `https://apibarchay.rudenko.tech` не должны кешироваться Service Worker'ом. В текущей реализации sw.js кеширует все запросы через `fetch` включая потенциальные API-запросы.

#### Что делать при проблемах с обновлением данных

**1. Исключить API из кеша в sw.js** *(рекомендуется)*:

```javascript
self.addEventListener('fetch', event => {
  // Не кешировать запросы к API
  if (event.request.url.startsWith('https://apibarchay.rudenko.tech')) {
    return; // браузер обрабатывает запрос самостоятельно
  }
  event.respondWith(
    fetch(event.request)
      .then(response => {
        const clone = response.clone();
        caches.open(CACHE_NAME).then(cache => cache.put(event.request, clone));
        return response;
      })
      .catch(() => caches.match(event.request))
  );
});
```

**2. Очистка кеша вручную** (DevTools → Application → Storage → Clear site data).

**3. Обновление версии кеша:** при изменении данных меню (если данные хранятся локально) увеличьте версию в `sw.js`:
```javascript
const CACHE_NAME = 'libo-v2'; // было v1
```
Это заставит Service Worker удалить старый кеш при следующей активации.

---

## 6. Конфигурация frontend

### Способ 1: window.__CONFIG__ (рекомендуется для деплоя)

В `index.html`, **до** подключения скриптов, добавьте блок конфигурации:

```html
<!-- API Configuration -->
<script>
  window.__CONFIG__ = {
    API_BASE: 'https://apibarchay.rudenko.tech'
  };
</script>
```

Файл `js/app.js` автоматически считает значение:

```javascript
const API_BASE = (window.__CONFIG__ && window.__CONFIG__.API_BASE) || '';
```

Если `window.__CONFIG__` не задан или `API_BASE` пуст — приложение работает с локальными данными (`MENU_DATA` / `QUIZ_DATA`).

### Способ 2: Прямая правка js/app.js

Найдите строку:
```javascript
const API_BASE = (window.__CONFIG__ && window.__CONFIG__.API_BASE) || '';
```
И замените пустую строку на URL API:
```javascript
const API_BASE = (window.__CONFIG__ && window.__CONFIG__.API_BASE) || 'https://apibarchay.rudenko.tech';
```

> **Внимание:** При этом подходе сборка с жёстко вшитым URL сложнее обновлять без правки исходника.

---

## 7. Деплой сайта

### 7.1 Статический frontend

**Файлы для загрузки на хостинг `barchay.rudenko.tech`:**

```
index.html
css/
  style.css
js/
  app.js
  menu-data.js
  quiz-data.js
icons/
  icon-192.png
  icon-512.png
manifest.json
sw.js
CNAME
```

**Варианты размещения:**

- **GitHub Pages:** загрузите файлы в репозиторий, включите Pages в настройках, укажите домен `barchay.rudenko.tech` в поле Custom Domain (файл `CNAME` уже присутствует в репозитории).
- **nginx / Caddy:** скопируйте файлы в директорию сервера, настройте виртуальный хост (см. раздел 7.3).
- **Netlify / Vercel / Cloudflare Pages:** drag-and-drop или подключение репозитория через GUI.

**Настройка API_BASE перед загрузкой:**

Отредактируйте `index.html` — добавьте блок перед `<script src="js/menu-data.js">`:

```html
<!-- API Configuration -->
<script>
  window.__CONFIG__ = {
    API_BASE: 'https://apibarchay.rudenko.tech'
  };
</script>
```

### 7.2 DNS

Для обоих поддоменов добавьте A/AAAA записи, указывающие на IP вашего сервера.

| Запись | Тип | Значение |
|--------|-----|---------|
| `barchay.rudenko.tech` | A | `<IP-сервера-frontend>` |
| `barchay.rudenko.tech` | AAAA | `<IPv6-сервера-frontend>` *(если есть)* |
| `apibarchay.rudenko.tech` | A | `<IP-сервера-API>` |
| `apibarchay.rudenko.tech` | AAAA | `<IPv6-сервера-API>` *(если есть)* |

> Frontend и API могут находиться на одном или на разных серверах.

Проверка DNS-записей:
```bash
dig barchay.rudenko.tech A
dig apibarchay.rudenko.tech A
```

### 7.3 TLS и обратный прокси (nginx)

#### Получение сертификатов Let's Encrypt (Certbot)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d barchay.rudenko.tech -d apibarchay.rudenko.tech
```

#### Конфигурация nginx для frontend

```nginx
server {
    listen 80;
    server_name barchay.rudenko.tech;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name barchay.rudenko.tech;

    ssl_certificate     /etc/letsencrypt/live/barchay.rudenko.tech/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/barchay.rudenko.tech/privkey.pem;

    root /var/www/barchaysite;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    # Кешируем статику надолго
    location ~* \.(css|js|png|ico|webp|json|woff2?)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

#### Конфигурация nginx для API

```nginx
server {
    listen 80;
    server_name apibarchay.rudenko.tech;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name apibarchay.rudenko.tech;

    ssl_certificate     /etc/letsencrypt/live/apibarchay.rudenko.tech/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/apibarchay.rudenko.tech/privkey.pem;

    # Лимит размера загружаемых файлов (должен совпадать с лимитами Go API)
    client_max_body_size 55M;

    location / {
        proxy_pass         http://127.0.0.1:8080;
        proxy_set_header   Host $host;
        proxy_set_header   X-Real-IP $remote_addr;
        proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_read_timeout 60s;
    }
}
```

> **Важно:** `client_max_body_size` в nginx должен быть не меньше суммарного лимита загрузки в Go API (рекомендуется 55 МБ при лимите API 50 МБ).

#### Альтернатива: Caddy

Caddy автоматически получает и продлевает сертификаты Let's Encrypt.

```caddyfile
barchay.rudenko.tech {
    root * /var/www/barchaysite
    file_server
    try_files {path} /index.html
}

apibarchay.rudenko.tech {
    request_body {
        max_size 55MB
    }
    reverse_proxy 127.0.0.1:8080
}
```

### 7.4 Структура директорий API

Рекомендуемая структура файлов на сервере:

```
/opt/barchay-api/
├── barchay-api          # скомпилированный бинарник Go
├── data/
│   └── menu.json        # хранилище меню
└── uploads/             # загруженные изображения
    ├── <uuid1>.jpg
    └── <uuid2>.png
```

Настройка прав:
```bash
sudo useradd -r -s /sbin/nologin barchay
sudo chown -R barchay:barchay /opt/barchay-api
sudo chmod 750 /opt/barchay-api
sudo chmod 755 /opt/barchay-api/uploads
```

Пример systemd unit-файла (`/etc/systemd/system/barchay-api.service`):

```ini
[Unit]
Description=Barchay API Server
After=network.target

[Service]
User=barchay
WorkingDirectory=/opt/barchay-api
ExecStart=/opt/barchay-api/barchay-api
Restart=on-failure
RestartSec=5s
EnvironmentFile=/etc/barchay-api/env

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now barchay-api
sudo systemctl status barchay-api
```

### 7.5 Переменные окружения и секреты

Переменные окружения API-сервера:

| Переменная | Обязательна | Пример | Описание |
|------------|-------------|--------|----------|
| `ADMIN_API_KEY` | **да** | `super-secret-key-32chars` | Ключ авторизации для `/admin/` эндпоинтов. |
| `PORT` | нет | `8080` | Порт, на котором слушает API (по умолчанию `8080`). |
| `DATA_DIR` | нет | `/opt/barchay-api/data` | Путь к директории с `menu.json`. |
| `UPLOADS_DIR` | нет | `/opt/barchay-api/uploads` | Путь к директории с изображениями. |
| `PUBLIC_URL` | нет | `https://apibarchay.rudenko.tech` | Базовый URL для формирования ссылок на изображения. |

#### Генерация безопасного ключа

```bash
openssl rand -hex 32
# или
python3 -c "import secrets; print(secrets.token_hex(32))"
```

#### Хранение env-файла

```bash
sudo mkdir -p /etc/barchay-api
sudo tee /etc/barchay-api/env > /dev/null <<EOF
ADMIN_API_KEY=ваш-сгенерированный-ключ
PORT=8080
DATA_DIR=/opt/barchay-api/data
UPLOADS_DIR=/opt/barchay-api/uploads
PUBLIC_URL=https://apibarchay.rudenko.tech
EOF
sudo chmod 600 /etc/barchay-api/env
sudo chown barchay:barchay /etc/barchay-api/env
```

> **Никогда не коммитьте `ADMIN_API_KEY` в репозиторий.** Используйте `.gitignore` или secrets-хранилище (Vault, GitHub Secrets, etc.).

---

## Быстрый старт (checklist)

- [ ] Клонировать репозиторий, добавить блок `window.__CONFIG__` с `API_BASE` в `index.html`
- [ ] Загрузить статические файлы на хостинг `barchay.rudenko.tech`
- [ ] Настроить DNS A-записи для обоих поддоменов
- [ ] Получить TLS-сертификаты (Certbot / Caddy)
- [ ] Настроить nginx/Caddy: проксирование API, `client_max_body_size 55M`
- [ ] Создать директории `/opt/barchay-api/data/` и `/opt/barchay-api/uploads/`
- [ ] Создать начальный `menu.json` с категориями
- [ ] Сгенерировать `ADMIN_API_KEY` и сохранить в `/etc/barchay-api/env`
- [ ] Запустить `barchay-api` через systemd
- [ ] Проверить работу: `curl https://apibarchay.rudenko.tech/health`
- [ ] Добавить первый пункт меню через `POST /admin/items-with-images`
- [ ] Открыть `https://barchay.rudenko.tech` и убедиться, что меню загружается из API
