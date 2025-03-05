## gcalc-server documentation

## Table of contents
1. [Possible Errors](./possible_errors.md)

## Introduction
gcalc-server uses the port from the ```PORT`` environment variable. If it is not found, the default port ```8080``` is used. The server only accepts POST requests.

### Adding an expression
The server adds an expression from the request body in ```JSON``` format.

Example request:
```Bash.
curl localhost:<port>/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '<expression>'
```

The server returns the `id` number of the expression:
``` JSON
{
    "id":"0"
}
```
Possible response statuses:
- 201 - Expression successfully added
- 422 - Expression is incorrectly added
- 500 - Internal server error

### Checking the status of the expression
After adding an example, you can check the status of its execution:
``` Bash
curl http://localhost:<port>/api/v1/expressions/<id>
```

``` JSON
{
    "id": "0",
    "status": "in progress"
}
```

If the example has already been solved, the response will contain the result of the expression:
``` JSON
{
    { "expression": {
        "id": "0",
        "status": "complete",
        "result":4
    }
}
```

Possible response statuses:
- 200 - Expression found
- 404 - No such expression
- 500 - Internal server error

### Obtaining a list of all expressions
To get a list of all entered expressions:
``` Bash
curl http://localhost:<port>/api/v1/expressions
```

```JSON
{
    "expressions": [
        {
            "id": "0",
            "status": "complete",
            "result":9
        },
        {
            "id": "1",
            "status": "complete",
            "result":3
        }
    ]
}
```

Possible response statuses:
- 200 - Successfully received the list of expressions
- 500 - Internal server error

---
### Creating a POST request

To create a POST request to the server, Type in the console:
``` Bash
curl localhost:<port>/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '<expression>'
```

Example:
``` Bash
curl localhost:8080/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '{“expression”: “5+4"}'
```
