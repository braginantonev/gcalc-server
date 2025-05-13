# gcalc-server
A web service that accepts a complex expression in the form of JSON and returns the result. Available actions: +, -, *, /. And the precedence operators ().

## Installation
1. Install [Go](https://go.dev/doc/install)
2. Install [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
3. Clone repository: ```git clone https://github.com/braginantonev/gcalc-server```

## Usage
### Run Server
1. Go to installed repository
2. Install dependencies: ```go mod tidy```
3. Setup your server:
   - Create ```.env``` file: 
     * For Linux: ```touch .env```
     * For Windows: ```echo.  > .env``` (P.S. У меня нет винды, простите если не сработает. Создайте файл другим образом)
   - In ```.env``` write the following environments: 
``` .env
PORT=8080
GRPCPort=5000
COMPUTING_POWER=5
JWTSecretSignature="your_super_secret_signature"
```
4. For run http server: ```go run cmd/main.go``` (For Windows use Git Bash or WSL).
   * If you have error:
```
!!! Attention !!!
JWT signature in env JWTSecretSignature not found.
Please go to README.md - Usage, and follow the instruction!
```
   Run server with next command: ``` PORT=8080 GRPCPort=5000 COMPUTING_POWER=5 JWTSecretSignature="your_super_secret_signature" go run cmd/main.go ```

5. To stop server use ```Ctrl+C```

### Available request (examples)
1. Registration (```localhost:<your port>/api/v1/register```):
   * Curl: ```curl localhost:8080/api/v1/register -X POST --header 'Content-Type: application/json' --data '{"name": "Anton","password": "12345"}'```
   * Postman: ![](./docs/screens/regestration_example.png)
   * In result you get ```JWT``` token. Use him for calculate and get expressions
2. Login (```localhost:<your port>/api/v1/login```): 
   * Curl: ```curl localhost:8080/api/v1/login -X POST --header 'Content-Type: application/json' --data '{"name": "Anton","password": "12345"}'```
   * Postman: ![](./docs/screens/login_example.png)
   * In result you get ```JWT``` token. Use him for calculate and get expressions
3. Calculate (```localhost:<your port>/api/v1/calculate```):
   * If you authorized (have ```JWT``` token):
      - Curl: ```curl localhost:<your port>/api/v1/calculate -X POST -H 'Content-Type: application/json' -H "Authorization: <Your JWT token>" --data '<выражение>'```. 
      - Example: ```curl localhost:8080/api/v1/calculate -X POST -H 'Content-Type: application/json' -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcwNzk3MDcsImlhdCI6MTc0Njk5MzMwNywibmFtZSI6IkFudG9uIiwibmJmIjoxNzQ2OTkzMzA3fQ.TgJoztS5R0neg_Am6fKTVoXNEexz53ev1g36rZvjqgg" --data '{"expression": "1+1"}'```
      - Postman: ![](./docs/screens/calculate_jwt_example.png) ![](./docs/screens/calculate_example.png)
   * If you not authorized:
      - Curl: ```curl localhost:<your port>/api/v1/calculate -X POST -H 'Content-Type: application/json' --data '<выражение>'```
      - Example: ```curl localhost:8080/api/v1/calculate -X POST -H 'Content-Type: application/json' --data '{"expression": "1+1"}'```
      - Postman: ![](./docs/screens/calculate_example.png)
4. Get all your expressions (```localhost:<your port>/api/v1/expressions```):
   * If you authorized (have ```JWT``` token):
      - Curl: ```curl localhost:<your port>/api/v1/expressions -X GET -H "Authorization: <Your JWT token>"```
      - Example: ```curl localhost:8080/api/v1/expressions -X GET -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcwNzk3MDcsImlhdCI6MTc0Njk5MzMwNywibmFtZSI6IkFudG9uIiwibmJmIjoxNzQ2OTkzMzA3fQ.TgJoztS5R0neg_Am6fKTVoXNEexz53ev1g36rZvjqgg"```
      - Postman: ![](./docs/screens/expressions_jwt_example.png)
   * If you not authorized:
      - Curl: ```curl localhost:<your port>/api/v1/expressions -X GET```
      - Example: ```curl localhost:8080/api/v1/expressions -X GET```
      - Postman: ![](./docs/screens/expressions_example.png)
5. Get expression (```localhost:<your port>/api/v1/expressions/<id your expression>```):
   * If you authorized (have ```JWT``` token):
      - Curl: ```curl localhost:<your port>/api/v1/expressions/<id your expression> -X GET -H "Authorization: <Your JWT token>"```
      - Example: ```curl localhost:8080/api/v1/expressions/0 -X GET -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcwNzk3MDcsImlhdCI6MTc0Njk5MzMwNywibmFtZSI6IkFudG9uIiwibmJmIjoxNzQ2OTkzMzA3fQ.TgJoztS5R0neg_Am6fKTVoXNEexz53ev1g36rZvjqgg"```
      - Postman: ![](./docs/screens/expression_jwt_example.png)
   * If you not authorized:
      - Curl: ```curl localhost:<your port>/api/v1/expressions/<id your expression> -X GET```
      - Example: ```curl localhost:8080/api/v1/expressions/0 -X GET```
      - Postman: ![](./docs/screens/expression_example.png) 

## Documentation
For documentation on the library, usage examples and possible errors, follow the link: 
* [Русская](https://github.com/Antibrag/gcalc-server/blob/main/docs/ru/index.md)
* [English](https://github.com/Antibrag/gcalc-server/blob/main/docs/en/index.md)

