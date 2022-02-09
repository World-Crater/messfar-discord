FROM superj80820/golang-process-image-base
WORKDIR /app
ADD . /app
COPY config.yaml.example /app/config.yaml
RUN go build -o app
RUN mkdir /app/images
ENTRYPOINT /app/app