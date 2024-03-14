FROM alpine:latest as production
RUN mkdir /app
COPY mailerApp /app
CMD [ "/app/mailerApp" ]