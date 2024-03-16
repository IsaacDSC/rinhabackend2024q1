# asd


Transacao para todas as querys de escrita

DB_PORT=5432 go run ./cmd/api/main.go



psql -U root -W rinha_backend -p 5432 -h localhost

DB_PORT=5432 go run ./cmd/main.go


## Conceitos
 - Cache
 - Processamento async -> queue para banco de dados
 - Errar rápido
   - Definir tempos limits para transação (ctx, cancel)


deploy.resources.limits.cpu 1.5 – uma unidade e meia de CPU distribuída entre todos os seus serviços
deploy.resources.limits.memory 550MB – 550 mega bytes de memória distribuídos entre todos os seus serviços

01-CPU 
  - 0.6 App
  - 0.17 Nginx
  - 0.13 Database
01-MEM
  - 200MB APP -> 55MB
  - 10MB Nginx
  - 140MB Database

02-CPU 
  - 0.6 App -> 14
  - 0.17 Nginx -> 7
  - 0.13 Database -> 14
02-MEM
  - 55MB APP -> 55
  - 55MB Nginx -> 15
  - 440MB Database -> 200

03-CPU
- 0.6 App -> 14
- 0.17 Nginx -> 7
- 0.13 Database -> 14
03-MEM
- 55MB APP -> 35
- 45MB Nginx
- 440MB Database -> 200

https://stackoverflow.com/questions/28265717/worker-connections-are-not-enough
Nginx Message logs -> 2024/02/22 12:01:17 [alert] 29#29: 1000 worker_connections are not enough
define -> 256

Postgres message logs -> 2024-02-22 12:00:46.538 UTC [1503] FATAL:  sorry, too many clients already
postgres config file -> max_conn_postgres 30

Adicionar logs críticos de erros na aplicação para coletar depois com as métricas



Aumentar o worker do nginx para o tanto que a aplicação aguenta


- [x] validaçoes mais rápidas para retornar o erro logo
- [ ] woker topando no topo connection -> 1000
- [ ] api com 35mb não aguenta e cai uma depois derruba as outras
  - Api esta ficando demorando para responder e deixando pools https connectados até 6000ms
    - adicionar um context timout para nao permitir isso
- [ ] database demorando para reponder possivelmente
   - Ajustar a performance das querys 
     - Adicionar index id clients
   - Economizar as transaçoes
   - As connections do database -> db.SetMaxOpenConns(83) db.SetMaxIdleConns(20)


### References
https://martinfowler.com/articles/lmax.html

### Analysis

Gatling
ALL Request OK -> 20206 requests
timer < 800ms -> 1309 requests
timer >= 800ms timer < 12000ms -> 0 requests
timer >= 1200ms -> 18835 requests

Database psql
ERROR:  Debit must be positive -> 1670 requests

Nginx Loggers
1000 worker_connections are not enough -> 1537 requests
live upstreams while connecting to upstream -> 2012 requests

/transacoes -> 5664 requests
/extrato -> 182 requests

App Loggers
app01 -> superfluous response.WriteHeader call from main.CreateTransaction -> 4898 requests
app01 -> superfluous response.WriteHeader call from main.GetTransactions -> 117 requests

app02 -> superfluous response.WriteHeader call from main.CreateTransaction ->  4878 request
app02 -> superfluous response.WriteHeader call from main.GetTransactions -> 137 requests

superfluous response.WriteHeader call from main.CreateTransaction -> 9776 requests
superfluous response.WriteHeader call from main.GetTransactions -> 254 requests

### Actions
OK -> Ao iniciar definir cache de (debit must be positive) Accounts (1,2,3,4,5)
OK -> Verificar status Code error (Debit must be positive) return 422

[//]: # (Ao fazer uma transação escrever cache account)
Verificar superfluous response.WriteHeader call Golang 
Ao fechar uma connection http matar o processo (ctx timout)
Mudar de http1.1 para http2
