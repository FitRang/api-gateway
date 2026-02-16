FROM kong:3.6

COPY kong.yml /usr/local/kong/kong.yml

COPY plugins/token-introspect-plugin /usr/local/share/lua/5.1/kong/plugins/token-introspect-plugin
COPY plugins/token-refresh-plugin /usr/local/share/lua/5.1/kong/plugins/token-refresh-plugin

ENV KONG_DATABASE=off \
    KONG_DECLARATIVE_CONFIG=/usr/local/kong/kong.yml \
    KONG_PROXY_ACCESS_LOG=/dev/stdout \
    KONG_ADMIN_ACCESS_LOG=/dev/stdout \
    KONG_PROXY_ERROR_LOG=/dev/stderr \
    KONG_ADMIN_ERROR_LOG=/dev/stderr \
    KONG_ADMIN_LISTEN=0.0.0.0:8001 \
    KONG_PLUGINS=bundled,token-introspect-plugin,token-refresh-plugin

HEALTHCHECK CMD curl -f http://localhost:8001/status || exit 1
