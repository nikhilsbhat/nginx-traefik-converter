### Description: Dockerfile for nginx-traefik-converter
FROM alpine:3.23

COPY nginx-traefik-converter /

# Starting
ENTRYPOINT [ "/nginx-traefik-converter" ]