# Anagrab

[Anagrab](https://kucicm.github.io/anagrab) finds single or multi words anagrams.

## Build

``` bash
tinygo build -o public/main.wasm -target wasm ./src/
```

## Local Start

Step 1. Build docker image with nginx.

``` bash
docker build -t anagrab .
```

Step 2. Start the container and add volume path.

``` bash
docker run -d -p 8080:80 --name anagrab -v "$(pwd)":/usr/share/nginx/html anagrab
```

Step 3. Any changes to files will be reflected on the [hosted pages](http://localhost:8080).

