# Docker и деплой на VPS

Это руководство объясняет, как запустить **Barchaysite** (фронтенд + Go-API) в Docker на любом VPS-сервере с доменом и HTTPS.

---

## Содержание

1. [Архитектура контейнеров](#1-архитектура-контейнеров)
2. [Быстрый старт (локально)](#2-быстрый-старт-локально)
3. [Деплой на VPS — пошагово](#3-деплой-на-vps-пошагово)
   - [3.1 Подготовка сервера](#31-подготовка-сервера)
   - [3.2 Клонирование репозитория](#32-клонирование-репозитория)
   - [3.3 Настройка окружения (.env)](#33-настройка-окружения-env)
   - [3.4 Запуск контейнеров](#34-запуск-контейнеров)
   - [3.5 HTTPS через Caddy (рекомендуется)](#35-https-через-caddy-рекомендуется)
   - [3.6 HTTPS через nginx + Certbot (альтернатива)](#36-https-через-nginx--certbot-альтернатива)
4. [Переменные окружения](#4-переменные-окружения)
5. [Управление данными и обновление](#5-управление-данными-и-обновление)
6. [Проверка работоспособности](#6-проверка-работоспособности)
7. [Частые проблемы](#7-частые-проблемы)

---

## 1. Архитектура контейнеров

```
Браузер
  │
  ▼ HTTPS :443
┌─────────────────────┐
│  Caddy / host nginx │  ← обратный прокси на хосте (SSL-терминация)
└──────────┬──────────┘
           │ HTTP
           ▼ :80
┌─────────────────────┐
│  frontend (nginx)   │  ← Docker-контейнер, порт 80
│  статические файлы  │
│  + proxy API-путей  │
└──────────┬──────────┘
           │ внутренняя сеть Docker (barchay)
           ▼ :8080
┌─────────────────────┐
│  api (Go)           │  ← Docker-контейнер, порт 8080
│  /menu /quiz /admin │
│  /uploads           │
└─────────────────────┘
           │
    Docker volumes
    api_data / api_uploads
```

**Что делает каждый контейнер:**

| Контейнер  | Образ         | Роль                                                       |
|------------|---------------|------------------------------------------------------------|
| `api`      | golang→alpine | Go HTTP-сервер: меню, квиз, загрузка изображений          |
| `frontend` | nginx:alpine  | Раздаёт статику, проксирует `/menu`, `/quiz`, `/uploads` к `api` |

Фронтенд и API находятся в одной внутренней Docker-сети `barchay` и не видны снаружи напрямую — доступ через обратный прокси на хосте.

---

## 2. Быстрый старт (локально)

```bash
# 1. Скопируйте пример env-файла
cp .env.example .env

# 2. Заполните ADMIN_API_KEY (можно любую строку для теста)
#    Остальные значения подходят для локальной разработки

# 3. Соберите и запустите
docker compose up --build

# 4. Откройте браузер
open http://localhost
```

Сайт доступен на `http://localhost`, API на `http://localhost/menu`.

---

## 3. Деплой на VPS — пошагово

### 3.1 Подготовка сервера

Минимальные требования: **1 CPU / 1 GB RAM / Ubuntu 22.04+**.

```bash
# Обновляем пакеты
sudo apt update && sudo apt upgrade -y

# Устанавливаем Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker

# Проверяем
docker --version
docker compose version
```

Настройте **DNS A-записи** у вашего регистратора доменов, указав IP сервера:

```
yourdomain.com      A  <IP сервера>
www.yourdomain.com  A  <IP сервера>   # опционально
```

> Если API будет на отдельном поддомене — добавьте ещё одну запись:
> ```
> api.yourdomain.com  A  <IP сервера>
> ```

### 3.2 Клонирование репозитория

```bash
git clone https://github.com/Slava47/barchaysite.git /opt/barchaysite
cd /opt/barchaysite
```

### 3.3 Настройка окружения (.env)

```bash
cp .env.example .env
nano .env          # или любой другой редактор
```

Минимально необходимые значения:

```dotenv
# Секретный ключ для эндпоинта /admin/items-with-images
ADMIN_API_KEY=сюда_длинную_случайную_строку

# Публичный URL сайта (используется в URL изображений, которые возвращает API)
PUBLIC_URL=https://yourdomain.com

# Разрешённые CORS-источники (обычно совпадает с PUBLIC_URL)
ALLOWED_ORIGINS=https://yourdomain.com

# Оставьте пустым — тогда nginx проксирует API через тот же домен
API_BASE=

# Порт, на котором frontend-контейнер слушает на хосте
# Caddy/nginx будут проксировать именно на этот порт
FRONTEND_PORT=8000
```

> **Совет:** Для генерации надёжного `ADMIN_API_KEY`:
> ```bash
> openssl rand -hex 32
> ```

### 3.4 Запуск контейнеров

```bash
cd /opt/barchaysite

# Сборка образов и запуск в фоне
docker compose up --build -d

# Просмотр логов
docker compose logs -f
```

Убедитесь, что оба контейнера работают:

```bash
docker compose ps
```

Вывод должен выглядеть примерно так:

```
NAME                  STATUS
barchaysite-api-1     Up (healthy)
barchaysite-frontend-1  Up
```

### 3.5 HTTPS через Caddy (рекомендуется)

Caddy автоматически получает и обновляет TLS-сертификаты (Let's Encrypt).

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \
  | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \
  | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy -y
```

Создайте `/etc/caddy/Caddyfile`:

```caddyfile
yourdomain.com {
    reverse_proxy localhost:8000
}
```

> Замените `yourdomain.com` и порт `8000` (должен совпадать с `FRONTEND_PORT` в `.env`).

```bash
sudo systemctl enable --now caddy
sudo systemctl reload caddy
```

Готово — сайт доступен по `https://yourdomain.com`.

### 3.6 HTTPS через nginx + Certbot (альтернатива)

```bash
sudo apt install -y nginx certbot python3-certbot-nginx
```

Создайте `/etc/nginx/sites-available/barchaysite`:

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass         http://localhost:8000;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        client_max_body_size 55M;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/barchaysite /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx

# Получить сертификат
sudo certbot --nginx -d yourdomain.com

# Certbot добавит HTTPS-блок автоматически
sudo systemctl reload nginx
```

---

## 4. Переменные окружения

Все переменные задаются в файле `.env` (на основе `.env.example`).

| Переменная       | Значение по умолчанию        | Описание                                                                 |
|------------------|------------------------------|--------------------------------------------------------------------------|
| `ADMIN_API_KEY`  | *(обязательно заполнить)*    | Секретный ключ для `/admin/items-with-images`                            |
| `PUBLIC_URL`     | `http://localhost:8080`      | Публичный URL API (используется в URL изображений)                      |
| `ALLOWED_ORIGINS`| `https://barchay.rudenko.tech` | CORS-источники (через запятую)                                         |
| `API_BASE`       | *(пусто)*                    | URL API для браузера. Пусто = same-origin через nginx-прокси            |
| `FRONTEND_PORT`  | `80`                         | Порт frontend-контейнера на хосте                                       |

**Переменные только для API-контейнера** (редко нужно менять):

| Переменная    | По умолчанию  | Описание                              |
|---------------|---------------|---------------------------------------|
| `PORT`        | `8080`        | Внутренний порт Go-сервера            |
| `DATA_DIR`    | `/app/data`   | Директория с `menu.json` и `quiz.json`|
| `UPLOADS_DIR` | `/app/uploads`| Директория загруженных изображений    |

---

## 5. Управление данными и обновление

### Начальные данные меню

При первом запуске директория `/app/data` пуста — API вернёт пустое меню.  
Загрузите начальный `menu.json` через эндпоинт `POST /admin/items-with-images`  
или скопируйте файл напрямую в Docker-volume:

```bash
# Скопировать локальный menu.json в контейнер
docker compose cp ./my-menu.json api:/app/data/menu.json

# Перезапустить API, чтобы подхватил файл
docker compose restart api
```

### Обновление сайта

```bash
cd /opt/barchaysite

# Получить последние изменения
git pull

# Пересобрать образы и перезапустить
docker compose up --build -d
```

Данные (меню, изображения) хранятся в Docker-volumes и **не теряются** при пересборке.

### Резервная копия данных

```bash
# Создать архив с данными и загрузками
docker run --rm \
  -v barchaysite_api_data:/data \
  -v barchaysite_api_uploads:/uploads \
  -v $(pwd):/backup \
  alpine tar czf /backup/barchay-backup-$(date +%Y%m%d).tar.gz /data /uploads
```

---

## 6. Проверка работоспособности

```bash
# Проверить статус контейнеров
docker compose ps

# Логи API
docker compose logs api

# Логи фронтенда (nginx)
docker compose logs frontend

# Тест API (замените домен на свой)
curl https://yourdomain.com/health

# Ожидаемый ответ:
# {"status":"ok"}

# Тест меню
curl https://yourdomain.com/menu | head -c 200
```

---

## 7. Частые проблемы

### Контейнер `api` не запускается

```bash
docker compose logs api
```

- **ADMIN_API_KEY не задан** — предупреждение в логах, сервер всё равно запускается. Задайте ключ в `.env`.
- **Ошибка чтения data/** — проверьте права Docker-volume.

### Сайт открывается, но меню пустое

1. Проверьте, что `/menu` отвечает: `curl http://localhost:8000/menu`
2. Убедитесь, что `menu.json` существует в volume (см. раздел 5).
3. Проверьте логи: `docker compose logs api`.

### CORS-ошибки в браузере

Установите `ALLOWED_ORIGINS` в `.env` равным домену вашего сайта (с протоколом):

```dotenv
ALLOWED_ORIGINS=https://yourdomain.com
```

Если домен сменился, перезапустите контейнеры: `docker compose restart`.

### Порт 80 занят

Измените `FRONTEND_PORT` в `.env` на свободный порт (например, `8000`) и обновите конфигурацию обратного прокси.

### Загрузка изображений не работает (413 Request Entity Too Large)

Nginx внутри frontend-контейнера уже имеет `client_max_body_size 55M`.  
Если перед ним стоит ещё один nginx или Caddy — добавьте туда аналогичное ограничение:

- **Caddy:** `request_body { max_size 55MB }` в блоке сайта
- **nginx:** `client_max_body_size 55M;` в блоке `server` или `location`
