FROM debian

COPY ./app /app
COPY ./yml /yml
COPY ./configs /configs

ENTRYPOINT /app

