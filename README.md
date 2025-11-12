## Primeiro desafio técnico de conclusão
Rate limiter

## Texto descritivo do desafio
Objetivo: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.
Descrição: O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:
Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.
Requisitos:

O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
O sistema deve responder adequadamente quando o limite é excedido:
Código HTTP: 429  
Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
Todas as informações de "limiter" devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
A lógica do limiter deve estar separada do middleware.
Exemplos:  

Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.  
Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.  
Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.
Dicas:  

Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.  
Entrega:  

O código-fonte completo da implementação.  
Documentação explicando como o rate limiter funciona e como ele pode ser configurado.  
Testes automatizados demonstrando a eficácia e a robustez do rate limiter.  
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.  
O servidor web deve responder na porta 8080.  


## Como funciona neste projeto  
Limita o número de requisições, que cada IP ou Token pode fazer dentro de uma janela de tempo configurável.   
Se o limite for excedido, as próximas requisições são bloqueadas até que novo tempo de acesso seja reiniciado.  

## Configuração e Taxa utilizada  
Devem ser definidas no arquivo '.env'.  

### Variáveis de controle:  
- 'RATE_LIMITER_IP_MAX_REQUESTS': Número máximo de requisições por IP antes do bloqueio.  
- 'RATE_LIMITER_TOKEN_MAX_REQUESTS': Número máximo de requisições por token antes do bloqueio (padrão: 5)  
- 'RATE_LIMITER_TOKEN_BLOCK_TIME': Período de tempo (em segundos) para bloquear o token.  
- 'RATE_LIMITER_IP_BLOCK_TIME': Período de tempo (em segundos) para bloquear o IP.  
- 'IP_FAKE_TO_TESTER': Número IP, o qual será utilizado para testes. Caso não seja informado, será capturado o ip da conexão.  

## Como usar esse modelo sem o Docker

```
1. Subir um servidor de cache Redis com estes dados default  
- REDIS_HOST=redis  
- REDIS_PORT=6379  

2. Configure o arquivo `.env` conforme desejado "Variáveis de controle".

3. Subir o servior.  
- Entre na pasta do projeto;  
- Entre na pasta cmd/server e execute o server  
cd cmd/server  
go run main.go  

4. Em outro terminal execute um dos seguintes comandos:  
- Testar com o ip  
bash
curl http://localhost:8080/

- Testar com o token  
bash
curl --location 'http://localhost:8080' \
--header 'API_KEY: 125d588d848787we'  


5. Mensagens de retorno:
- Caso o acesso for válido a seguinte mensagem será enviada:
"Enpoint acessado com sucesso" com HTTP status 200.  

- No caso de acesso inválido a seguinte mensagem será enviada:
"Número máximo por tempo atingido." com HTTP status 429.  

```

## Como executar o projeto com Docker

```
1. Entre na pasta do projeto; 
bash

docker-compose up --build -d

2. Utilize as opções do passo 4 do item "Como usar esse modelo sem o Docker"
```

---