FROM alpine:latest

RUN mkdir /app

COPY limestoneApp /app

CMD [ "/app/limestoneApp"]