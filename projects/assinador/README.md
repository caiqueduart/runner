# Assinador (Core Java)

O **Assinador** é o componente central do sistema Runner. Ele é responsável pela lógica de assinatura digital e validação, podendo operar tanto como um utilitário de linha de comando quanto como um servidor HTTP persistente.

## Funcionalidades

- **Assinatura Digital**: Gera códigos de assinatura baseados no conteúdo dos arquivos.
- **Validação**: Verifica se um arquivo corresponde a um código de assinatura fornecido.
- **Modo Servidor**: Expõe endpoints HTTP para integração com a CLI ou outros sistemas.
- **Auto-desligamento**: Encerra-se automaticamente após um período de inatividade configurável.

## Como Executar (Desenvolvimento)

```bash
# Modo CLI (Assinar)
java -cp bin App sign --file documento.txt

# Modo Servidor
java -cp bin App server --port 8080 --timeout 5
```

## Endpoints da API (Modo Servidor)

| Rota        | Método | Descrição                                                           |
| :---------- | :----- | :------------------------------------------------------------------ |
| `/sign`     | `POST` | Recebe o nome do arquivo no corpo e retorna o código de assinatura. |
| `/validate` | `POST` | Recebe o nome do arquivo e valida a assinatura existente.           |
| `/health`   | `GET`  | Retorna o status de saúde e tempo de atividade do servidor.         |
| `/stop`     | `POST` | Encerra o servidor.                                                 |
