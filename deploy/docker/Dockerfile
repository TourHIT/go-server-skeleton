# docker build -t roandocker/gin-server-skeleton:0.0.1 Dockerfile .
# docker build --target builder -t roandocker/gin-server-skeleton:0.0.1 .

FROM golang:stretch as builder

WORKDIR /opt/app/

COPY server/config/Config.yaml /opt/app/config/Config.yaml
COPY server/web /opt/app/web/

RUN apt-get install git && \
	git clone https://github.com/codeandcode0x/go-server-skeleton.git && \
	cd gin-server-skeleton/ && \
    CGO_ENABLED=0 GOOS=linux go build -o /opt/app/server main.go


FROM ubuntu:16.04 as prod

WORKDIR /root/

COPY --from=0 /opt/app/server /root
COPY --from=0 /opt/app/config/Config.yaml  /root/config/Config.yaml
COPY --from=0 /opt/app/web /root/web/

RUN ls -all
RUN chmod 777 server

# health check
# HEALTHCHECK --interval=5s --timeout=3s --retries=3 \
#   CMD curl -fs http://localhost:8080 || exit 1

CMD ["/bin/sh", "-c", "/root/server"]