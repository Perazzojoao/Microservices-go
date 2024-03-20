FROM alpine:latest as production
RUN mkdir /app
COPY frontEndApp /app
CMD [ "/app/frontEndApp" ]