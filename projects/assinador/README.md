# Assinador (Core Java)

O **Assinador** é o componente central do sistema Runner. Ele é responsável pela lógica de assinatura digital e validação, podendo operar tanto como um utilitário de linha de comando quanto como um servidor HTTP persistente.

## Funcionalidades

- **Assinatura Digital**: Gera códigos de assinatura baseados no conteúdo dos arquivos.
- **Validação**: Verifica se um arquivo corresponde a um código de assinatura fornecido.
- **Modo Servidor**: Expõe endpoints HTTP para integração com a CLI ou outros sistemas.
- **Auto-desligamento**: Encerra-se automaticamente após um período de inatividade configurável.

## Como Executar (Desenvolvimento)

Pré-requisitos: JDK 21 e Maven 3.9 ou superior.

```bash
# Compilar e executar todos os testes
mvn verify

# O JAR executável será criado em target/assinador-v1.0.5.jar
java -jar target/assinador-v1.0.5.jar '{"command":"sign","file":"documento.txt"}'
```

O `App.java` exige que as requisições, mesmo via terminal, sejam feitas em formato **JSON**, garantindo um contrato unificado com a API HTTP.

```bash
# Modo CLI (Assinar) - Requer JSON
java -cp bin App '{"command": "sign", "file": "documento.txt"}'

# Modo Servidor
java -cp bin App server --port 8080 --timeout 5
```

## Contrato de Comunicação e API (Modo Servidor)

Tanto o modo CLI quanto as requisições HTTP (`/sign` e `/validate`) exigem **JSON** no corpo da requisição e retornam **JSON** estruturado.

### Exemplo de Requisição (POST /sign ou Argumento CLI):

```json
{
    "command": "sign",
    "file": "caminho/do/arquivo.txt",
    "flag": "--file"
}
```

### Endpoints Disponíveis

| Rota        | Método | Descrição                                                      | Exemplo de Resposta (Sucesso)                                                                                      |
| :---------- | :----- | :------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------- |
| `/sign`     | `POST` | Recebe o JSON com o caminho do arquivo e retorna a assinatura. | `{"message": "Arquivo assinado...", "fileName": "doc.txt", "code": "ABC", "signOutputPath": "...", "status": 200}` |
| `/validate` | `POST` | Recebe o JSON com o caminho do arquivo e valida a assinatura.  | `{"message": "Validação concluída.", "fileName": "doc.txt", "code": "ABC", "valid": true, "status": 200}`          |
| `/health`   | `GET`  | Retorna status de saúde e tempos de atividade em JSON.         | `{"status": "OK", "uptimeSeconds": 120, "code": 200}`                                                              |
| `/stop`     | `POST` | Encerra o servidor.                                            | `{"message": "Sinal de encerramento recebido.", "status": 200}`                                                    |

**Erros Estruturados:** Em caso de erro, o servidor retorna `400 Bad Request` (erro de usuário) ou `500 Internal Server Error` (erro de sistema) com o corpo:

```json
{
    "error": "Arquivo 'arquivo.txt' não encontrado.",
    "status": 400,
    "type": "user"
}
```

## Modo Desenvolvedor

Este projeto suporta um modo de execução para desenvolvimento que permite testar alterações no código Java sem a necessidade de gerar um novo JAR.

Para ativar:

1. Crie um arquivo `.env` na raiz do projeto.
2. Adicione a seguinte variável:

    ```env
    DEV_MODE=true
    ```

    Quando esta variável está ativa, a CLI de Assinatura irá compilar e executar os arquivos `.java` diretamente utilizando o comando `java -cp`, ignorando o JAR baixado do GitHub.
