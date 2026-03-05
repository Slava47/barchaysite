# ── Stage 1: nothing to compile — the frontend is pure static files ──────────
# We use a plain nginx image and copy all static assets into it.

FROM nginx:1.27-alpine

# Remove the default nginx config and replace with ours
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy all frontend static files
COPY index.html manifest.json sw.js CNAME config.js /usr/share/nginx/html/
COPY css/   /usr/share/nginx/html/css/
COPY js/    /usr/share/nginx/html/js/
COPY icons/ /usr/share/nginx/html/icons/
COPY picture/ /usr/share/nginx/html/picture/

# Entrypoint generates config.js from env vars at container startup
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
