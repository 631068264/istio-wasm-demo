go mod tidy
tinygo build -o plugin.wasm -scheduler=none -target=wasi  main.go

docker build -t xx.xx.xx:8080/istio/wasm_demo -f wasm-image.Dockerfile .

docker push xx.xx.xx:8080/istio/wasm_demo