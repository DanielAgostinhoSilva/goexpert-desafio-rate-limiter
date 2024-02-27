# goexpert-desafio-tecnico

## Desafio Rate Limiter

`Objetivo`: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

`Descrição`: O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

1. `Endereço IP`: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
2. `Token de Acesso`: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
   API_KEY: <TOKEN>
3. `Endereço IP`: As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.


### Requisitos:

* O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
* O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
* O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
* As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
* Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
* O sistema deve responder adequadamente quando o limite é excedido:
    * Código HTTP: 429
    * Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame 
* Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
* Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
* A lógica do limiter deve estar separada do middleware.


### Exemplos:

1. `Limitação por IP`: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
2. `Limitação por Token`: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.
3. Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.


### `Dicas`:

* Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.


### Entrega:

* O código-fonte completo da implementação.
* Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
* Testes automatizados demonstrando o funcionamento.
* Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
* Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado.


# Resumo do Projeto

Este é um projeto Go que implementa o conceito de "Rate Limiter", que é uma técnica para limitar o número de requisições que um cliente pode fazer a um serviço em um determinado intervalo de tempo.

## Pacote `ratelimit`

Este pacote fornece uma implementação de um serviço de controle de Rate Limiter. Ele usa a biblioteca `golang.org/x/time/rate` para implementar a funcionalidade de controle de taxa.

## Estruturas

Há duas estruturas principais neste pacote: `RateLimiterParams` e `RateLimiterService`.

- `RateLimiterParams`: Esta estrutura armazena um limitador de Rate Limiter, juntamente com a entidade "Visitor" associada a esse limitador.

- `RateLimiterService`: Esta estrutura representa o serviço de Rate Limiter. Ele contém um mapa de controles de visitante que associa chaves de identificação de visitante a instâncias de `RateLimiterParams`, além de uma instância de `sync.Mutex` para lidar com a concorrência.

## Funções

A seguir, as principais funções definidas neste pacote e o que cada uma delas faz:

- `NewRateLimiterService`: Esta função é usada para criar uma nova instância do serviço RateLimiterService.

- `GetRateLimiterParams`: Esta função é usada para obter os parâmetros do RateLimiterService para uma determinada chave de API e endereço IP. Se um limitador de taxa para o IP especificado não existe, ele cria um novo.

- `isBlockedTimeExpired`: Esta função privada, verifica se o tempo de bloqueio para um visitante expirou.

- `NewRateLimiter`: Esta função é usada para criar um novo limitador de taxa para uma chave de API e um endereço IP.

- `CleanupVisitors`: Esta função é usada para limpar os visitantes cujo tempo de bloqueio expirou.

- `LockVisitor`: Esta função bloqueia o visitante com o IP especificado.


## Pacote `middleware`

Este pacote fornece um middleware de Rate Limiter. 

## Estrutura `RateLimiterMiddleware`

Essa estrutura representa o middleware Rate Limiter. Ela possui um campo `service`, que é uma instância de `ratelimit.RateLimiterService`.

## Funções

- `NewRateLimiterMiddleware`: Essa função é usada para criar uma nova instância do middleware Rate Limiter. Ela recebe um serviço de controle de taxa (Rate Limiter Service) como parâmetro.

- `Handler`: Essa função é um manipulador HTTP que implementa o Rate Limiter. Ela verifica a taxa de requisições de um IP e uma chave de API, caso a taxa limite seja atingida, ela bloqueia o visitante e retorna um erro HTTP com o status `http.StatusTooManyRequests`.

  Se a taxa limite não foi atingida, ela permite que a requisição prossiga para o próximo manipulador HTTP na cadeia de middleware.



## Subindo o docker-compose.yaml para executar a aplicação


## Passos

1. **Configurar as variaveis de ambiente**

   Acesse a pasta raiz do projeto e edite as variaveis de ambiente que controla os limite de requicioes e tempo de bloqueio.
   
   |     Variável                     | Descrição                                                          |
   |----------------------------------|--------------------------------------------------------------------|
   | `WEB_SERVER_PORT`                | A porta em que o servidor web estará escutando.                    |
   | `MAX_REQUEST_PER_SECOND_BY_TOKEN`| O número máximo de requisições por segundo permitido por token.     |
   | `MAX_REQUEST_PER_SECOND_BY_IP`   | O número máximo de requisições por segundo permitido por IP.        |
   | `BLOCKED_TIME_PER_SECOND`        | O tempo em segundos que um usuário será bloqueado após ultrapassar o limite de requisições. |
   | `REDIS_ADDR`                     | O endereço do servidor Redis.                                       |
   | `REDIS_PASSWORD`                 | A senha para autenticação no servidor Redis, se necessário.         |
   | `REDIS_DB`                       | O número do banco de dados Redis a ser usado.                       |


2. **Executar a aplicação**

   Acessa a pasta raiz do projeto e execute o seguinte comando para fazer o build e inicializar a aplicacao:

    ```shell
    make build
    ```


3. **Testar a aplicação**

   Use o seguinte comando para realizar uma consulta no sistema.
    ```shell
    curl http://localhost:8080
    ```

   Exemplo de codigo para testar varias requisições para simular um bloqueio por token 
   ```go
   package main
   
   import (
   "fmt"
   "net/http"
   )
   
   func main() {
   
       client := &http.Client{}
       for i := 0; i < 20; i++ {
           req, err := http.NewRequest("GET", "http://localhost:8080", nil)
   
           if err != nil {
               fmt.Println("Falha ao criar a solicitação:", err)
               return
           }
   
           req.Header.Add("API_KEY", "12345")
   
           resp, err := client.Do(req)
           if err != nil {
               fmt.Println("Erro ao fazer a solicitação:", err)
               return
           }
           defer resp.Body.Close()
           fmt.Println("Resposta recebida, código de status:", resp.StatusCode)
   
       }
   }
   ```
 