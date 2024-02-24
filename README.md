# asd


Transacao para todas as querys de escrita



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