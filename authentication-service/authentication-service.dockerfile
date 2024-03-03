FROM alpine:latest as production
RUN mkdir /app
COPY authApp /app
CMD [ "/app/authApp" ]