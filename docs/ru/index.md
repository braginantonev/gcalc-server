# Документация gcalc-server

## Оглавление
1. [Возможные ошибки](./possible_errors.md)

## Введение
gcalc-server использует порт из переменной окружения ```PORT```. Если он не найдена, то используется стандартный порт ```8080```.

Сервер принимает выражение из тела запроса в формате ```JSON```.

Пример тела запроса:
``` JSON
{
    "expression": "1+1+(1+1)"
}
```

Сервер также возвращает результат вычисления в формате ```JSON``` и возвращает код ```200```:
```JSON
{
    "result": 4
}
```

Если в выражении были обнаруженны ошибки или тело запроса было составленно некорректно, то сервер возвращает ошибку:

Запрос:
```JSON
{
    "expression": "hi"
}
```

Ответ с кодом ```422```:
```JSON
{
    "error": "expression incorrect"
}
```


---
### Создание POST запроса
Для создания POST запроса на сервер, Введите в консоль:
```Bash
curl localhost:<порт>/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '<выражение>'
```
Пример:
```Bash
curl localhost:8080/api/v1/calculate -X POST --header 'Content-Type: application/json' --data '{"expression": "5+4"}'
```
Ответ:
```JSON
{
    "result":9
}
```

