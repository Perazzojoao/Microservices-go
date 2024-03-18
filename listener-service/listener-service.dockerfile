FROM alpine:latest as production
RUN mkdir /app
COPY listenerApp /app
CMD [ "/app/listenerApp" ]