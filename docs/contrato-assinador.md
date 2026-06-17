# Contrato CLI ↔ Assinador

Este documento define a API entre a CLI Go `assinatura` e o Assinador Java. A mesma operação e os mesmos objetos JSON são usados no modo HTTP e no modo local.

## Modos de execução

- **HTTP (padrão):** a CLI reutiliza ou inicia o servidor e chama `POST /sign` ou `POST /validate`.
- **Local:** ativado explicitamente com `--local`; a CLI executa `java -jar assinador.jar <json>`.

Opções globais da CLI:

```text
--local          executa o JAR diretamente
--port           porta HTTP, padrão 8080
--timeout        inatividade em minutos, padrão: 5
```

## Requisição

```json
{
    "command": "sign",
    "file": "caminho/do/arquivo.txt",
    "flag": "--file"
}
```

`command` aceita `sign` ou `validate`. O JAR é a autoridade para validar campos e regras de negócio.

## Respostas

Sucesso de assinatura:

```json
{
    "message": "Arquivo assinado com sucesso.",
    "fileName": "documento.txt",
    "code": "ABC123",
    "signOutputPath": "documento-txt-assinatura.txt",
    "status": 200
}
```

Sucesso de validação:

```json
{
    "message": "Validação concluída.",
    "fileName": "documento-txt-assinatura.txt",
    "code": "ABC123",
    "valid": true,
    "status": 200
}
```

Erro:

```json
{
    "error": "Descrição do problema e como corrigi-lo.",
    "status": 400,
    "type": "user"
}
```

`type=user` representa entrada inválida e resulta em exit code `1`. `type=system` representa falha operacional e resulta em exit code `2`.

## HTTP

| Rota        | Método | Resultado                           |
| ----------- | ------ | ----------------------------------- |
| `/health`   | `GET`  | Saúde e tempo restante da instância |
| `/sign`     | `POST` | Assinatura simulada                 |
| `/validate` | `POST` | Validação simulada                  |
| `/stop`     | `POST` | Encerramento controlado             |

Todas as respostas usam `Content-Type: application/json; charset=UTF-8`. A CLI impõe timeout de 10 segundos por requisição e não muda automaticamente para o modo local quando o servidor falha.

## Streams e códigos de saída

- `stdout`: resultado normal destinado ao usuário ou consumidor.
- `stderr`: diagnóstico e logs operacionais.
- `0`: sucesso.
- `1`: erro de entrada do usuário.
- `2`: erro de sistema, transporte ou dependência.
