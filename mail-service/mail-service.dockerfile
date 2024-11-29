FROM alpine:latest

RUN mkdir /app

COPY mail-app /app

CMD [ "/app/mail-app" ]
