FROM golang:alpine AS build
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}
ENV HOME /app
RUN apk add git && mkdir /build
ADD . /build/
WORKDIR /build
RUN chown -R 777 /build && go build -o checkbuild main.go

FROM alpine:latest AS runtime
ARG APP_VERSION
ENV APP_VERSION ${APP_VERSION}
ENV HOME /app
RUN mkdir -p /app \
  && adduser -S -D -H -h /app appuser
COPY --from=build /build/checkbuild ./app/
RUN chown -R appuser /app
WORKDIR /app
USER appuser
EXPOSE 8080
ENTRYPOINT ["./checkbuild"]
