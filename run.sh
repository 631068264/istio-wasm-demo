tinygo build -o plugin.wasm -scheduler=none -target=wasi  main.go
envoy -c envoy.yaml --concurrency 2 --log-format '%v'

# curl 'http://127.0.0.1:18000/'