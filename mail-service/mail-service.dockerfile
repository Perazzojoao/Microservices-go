FROM alpine:latest as production
RUN mkdir /app
COPY mailerApp /app
COPY templates /templates
CMD [ "/app/mailerApp" ]