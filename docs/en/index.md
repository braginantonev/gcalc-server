# Docs gcalc-server

## Contents
1. [Possible errors](./possible_errors.md)

## Introduction
gcalc-server uses the port from their environment variable ```PORT```. If it was not found, then the standard port ```8080``` is used.

The server accepts an expression from the request body in the ```JSON``` format.

Request body example:
``` JSON
{
    "expression": "1+1+(1+1)"
}
```

The server also transmits the result in the ```JSON``` format and return code ```200```:
```JSON
{
    "result": 4
}
```

If there was an error in the expression or the request body was entered incorrectly, the server will return an error:

Request:
```JSON
{
    "expression": "hi"
}
```

Response with code ```422```:
```JSON
{
    "error": "expression incorrect"
}
```


---
### Creating POST request
To create a POST request to the server, type this in the console
```Bash
curl localhost:<port>/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '<data>'
```
Example:
```Bash
curl localhost:8080/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '{"expression": "5+4"}'
```
Return:
```JSON
{
    "result":9
}
```

