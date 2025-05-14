# Docs gcalc-server

### [Contents](./index.md)

## Possible errors
### 1.  Divide by zero
Occurs when a number in an expression is divided by zero.

Example:
```JSON
{
    "expression": "1/0"
}
```
or
```JSON
{
    "expression": "2/(3-3)"
}
```

Return (Code ```422```):
```JSON
{
    "id": 0,
    "error": "divide by zero"
}
```

### 2. Expression empty
Occurs when expression is empty.

Example:
```JSON
{
    "expression":""
}
```

Return (Code ```422```):
```JSON
{
    "id": 0,
    "error":"expression empty"
}
```

### 3. OperationWithoutValue
Occurs when the operator in the expression is missing a value.

Example:
```JSON
{
    "expression": "+1+1"
}
```

Return (Code ```422```):
```JSON
{
    "id": 0,
    "error": "operation don't have a value"
}
```

### 4. BracketsNotFound
Occurs when an expression contains opening and closing brackets, but lacks closing or opening brackets, respectively.

Example:
```JSON
{
    "expression": "5-(8-5"
}
```

Return (Code ```422```):
```JSON
{
    "id": 0,
    "error": "not found opened or closed bracket"
}
```

### 5. ExpressionIncorrect
Occurs when there are extra characters in an expression.

Example:
```JSON
{
    "expression": "5+4+&"
}
```

Return (Code ```422```):
```JSON
{
    "id": 0,
    "error": "expression incorrect"
}
```

### 6. RequestBodyEmpty
Occurs when request body is empty.

Example:
```Bash
curl localhost:8080/api/v1/calculate -X POST --header "application/json"
```

Return (Code ```400```):
```JSON
{
    "id": 0,
    "error": "Request body empty"
}
```

### 7. UnsupportedBodyType
Occurs when the request body is not written in JSON.

Example:
```Bash
curl localhost:8080/api/v1/calculate -X POST --header "application/json" --data "1488pashalco"
```

Return (Code ```415```):
```JSON
{
    "id": 0,
    "error": "Unsupported request body type"
}
```

### 8. InternalError
Occurs when an unexpected error occurs.

