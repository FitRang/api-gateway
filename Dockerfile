FROM traefik:v3.0

USER root

RUN mkdir -p /plugins-local

COPY traefik.yml /etc/traefik/traefik.yml
COPY dynamic.yml /etc/traefik/dynamic.yml

COPY plugins /plugins-local

RUN chown -R traefik:traefik /plugins-local

USER traefik
