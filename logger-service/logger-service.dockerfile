FROM alpine:latest as production
RUN mkdir /app
COPY loggerServiceApp /app
CMD [ "/app/loggerServiceApp" ]