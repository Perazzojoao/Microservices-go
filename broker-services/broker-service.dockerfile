FROM alpine:latest as production
RUN mkdir /app
COPY brokerApp /app
CMD [ "/app/brokerApp" ]