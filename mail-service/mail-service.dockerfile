FROM alpine:latest

RUN mkdir /app

COPY mail-app /app

COPY templates /templates

CMD [ "/app/mail-app" ]
