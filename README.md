# go-wasm-runner

A simple Go server that executes a simple wasm program upon http request.
In this example, the wasm program is built with *tinygo* (see testdata subfolder).

# Running

`go run .`

Default server port is 8080, you can change it by passing `-port` parameter.

# Testing

Once started, go browser and type `http://localhost:8080/run?a=1&b=2`. Here, parameters `a` and `b` are numbers that needs to be added. The expected output: `result: 3`.

