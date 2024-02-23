# asd


Transacao para todas as querys de escrita



psql -U root -W rinha_backend -p 5432 -h localhost

DB_PORT=5432 go run ./cmd/main.go


## Conceitos
 - Errar rápido
   - Definir tempos limits para transação (ctx, cancel)
 -


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


nginx                                        1.80%     30.72MiB / 55MiB      55.85%    44.7MB / 39.5MB   12.8MB / 1.12MB   2
b92b2995e34f   api02                         12.74%    54.52MiB / 55MiB      99.13%    15.2MB / 11.3MB   122MB / 112MB     10
2abcf864d9a9   app01                         1.54%     54.09MiB / 55MiB      98.35%    14.8MB / 11.2MB   164MB / 154MB     10
fd70c8513dd9   rinhabackend2024_q1-db-1      12.91%    254.2MiB / 440MiB     57


- [ ] validaçoes mais rápidas para retornar o erro logo 
- [ ] woker topando no topo connection -> 1000
- [ ] api com 35mb não aguenta e cai uma depois derruba as outras
  - Api esta ficando demorando para responder e deixando pools https connectados até 6000ms
    - adicionar um context timout para nao permitir isso
- [ ] database demorando para reponder possivelmente
   - Ajustar a performance das querys 
     - Adicionar index id clients
   - Economizar as transaçoes
   - As connections do database -> db.SetMaxOpenConns(83) db.SetMaxIdleConns(20)
