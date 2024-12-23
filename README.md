# gcalc-server
A web service that accepts a complex expression in the form of JSON and returns the result. Available actions: +, -, *, /. And the precedence operators ().

## Installation
1. Install [Go](https://go.dev/doc/install)
2. Install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
3. Clone repository: ```git clone https://github.com/Antibrag/gcalc-server```

## Usage
1. Go to installed repository
2. For run http server: ```go run cmd/main.go``` (For Windows use Git Bash or WSL)
3. If you want start server on anouther port:
   1. For windows: ```sets PORT=<your_port>```
   2. For linux/mac: ```export PORT=<your_port>```
5. To stop server use ```Ctrl+C```

## Documentation
For documentation on the library, usage examples and possible errors, follow the link: 
* [Русская](https://github.com/Antibrag/gcalc-server/blob/main/docs/ru/index.md)
* [English](https://github.com/Antibrag/gcalc-server/blob/main/docs/en/index.md)

