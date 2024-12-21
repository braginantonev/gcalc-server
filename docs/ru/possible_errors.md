# Документация gcalc-server

### [Оглавление](./index.md)

## Возможные ошибки
### 1.  Divide by zero
Возникает если в выражении значение делится на ноль.

Пример:
```JSON
{
    "expression": "1/0"
}
```
или
```JSON
{
    "expression": "2/(3-3)"
}
```

Ответ (Код ```422```):
```JSON
{
    "error": "divide by zero"
}
```

---
### 2. Expression empty
Возникает если выражение пусто.

Пример:
```JSON
{
    "expression":""
}
```

Ответ (Код ```422```):
```JSON
{
    "error":"expression empty"
}
```

---
### 3. OperationWithoutValue
Возникает, когда в выражение присутствует оператор, но не указанно для него значение.

Пример:
```JSON
{
    "expression": "+1+1"
}
```

Ответ (Код ```422```):
```JSON
{
    "error": "operation don't have a value"
}
```

---
### 4. BracketsNotFound
Возникает, когда в выражении присутствуют открывающая и закрывающая скобка, но отсутствуют закрывающая или открывающая скобка, соответственно.

Пример:
```JSON
{
    "expression": "5-(8-5"
}
```

Ответ (Код ```422```):
```JSON
{
    "error": "not found opened or closed bracket"
}
```

---
### 5. ExpressionIncorrect
Возникает, когда в выражении присутствуют посторонние символы.

Пример:
```JSON
{
    "expression": "5+4+&"
}
```

Ответ (Код ```422```):
```JSON
{
    "error": "expression incorrect"
}
```

---
### 6. RequestBodyEmpty
Возникает, когда тело запроса отсутствует.

Пример:
```Bash
curl localhost:8080/api/v1/calculate -X POST --header "application/json"
```

Ответ (Код ```400```):
```JSON
{
    "error":"Request body empty"
}
```

---
### 7. UnsupportedBodyType
Возникает, когда тело запроса не в формате JSON.

Пример:
```Bash
curl localhost:8080/api/v1/calculate -X POST --header "application/json" --data "1488pashalco"
```

Ответ (Код ```415```):
```JSON
{
    "error":"Unsupported request body type"
}
```

### 8. InternalError
Возникает в необработанных случаях и внутренних ошибках.

