# go-musthave-group-diploma-tpl

получаю ошибку в тестах

```
Run gophermarttest \
=== RUN   TestGophermart
=== RUN   TestGophermart/TestEndToEnd
=== RUN   TestGophermart/TestEndToEnd/register_accrual_mechanic
    gophermart_e2e_test.go:50: 
        	Error Trace:	gophermart_e2e_test.go:50
        	            				suite.go:77
        	Error:      	Not equal: 
        	            	expected: 200
        	            	actual  : 400
        	Test:       	TestGophermart/TestEndToEnd/register_accrual_mechanic
        	Messages:   	Несоответствие статус кода ответа ожидаемому в хендлере 'POST http://localhost:39257/api/goods'
    gophermart_e2e_test.go:56: Оригинальный запрос:
        
        POST /api/goods HTTP/1.1
        Host: localhost:39257
        Accept: application/json
        Content-Type: application/json
        User-Agent: go-resty/2.6.0 (https://github.com/go-resty/resty)
        
        
        
        			{
        				"match": "IbRyCPtUCI0yN5",
        				"reward": 5,
        				"reward_type": "%"
        			}
```

локально

```
Status: 200 OK
Size: 30 Bytes
Time: 41 ms
Response
Headers5
Cookies
Results
Docs
{ }
Response Headers
Header
Value
content-encoding
gzip
content-type
text/plain; charset=utf-8
date
Sat, 19 Apr 2025 08:24:42 GMT
content-length
54
connection
close

Reward registered successfully
```