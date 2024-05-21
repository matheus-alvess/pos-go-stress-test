# Instruções

## Build e Run no Docker:

Para construir a imagem Docker, execute:

```sh
docker build -t stress-test .
```

Para executar o container Docker, use:

```sh
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```