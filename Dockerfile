FROM kong:3.6

USER root

COPY plugins /usr/local/share/lua/5.1/kong/plugins

ENV KONG_PLUGINS=bundled,firebase-token-introspect-plugin
ENV KONG_LUA_PACKAGE_PATH=/usr/local/share/lua/5.1/?.lua;;

COPY kong.yml /usr/local/kong/declarative/kong.yml

ENV KONG_DATABASE=off
ENV KONG_DECLARATIVE_CONFIG=/usr/local/kong/declarative/kong.yml
ENV KONG_PROXY_LISTEN=0.0.0.0:8000
ENV KONG_ADMIN_LISTEN=0.0.0.0:8001

USER kong
