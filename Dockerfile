FROM debian

COPY ./app /app
COPY ./yml /yml

ENTRYPOINT /app

