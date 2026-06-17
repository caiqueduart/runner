# Assinatura CLI

A **CLI de Assinatura** é a interface principal do sistema Runner, permitindo que usuários interajam com o motor de assinatura sem conhecer o ecosistema Java.

## Instalação e Provisionamento

A CLI gerencia automaticamente suas dependências. Na primeira execução, ela baixará o JDK 21 e o `assinador.jar` compatível para a pasta `~/.hubsaude`.

## Comandos Principais

O modo servidor HTTP é o padrão. Use `--local` para executar o JAR diretamente, `--port` para escolher a porta e `--timeout` para configurar os minutos de inatividade.

### 1. Assinar um arquivo

Realiza a assinatura de um documento. Se o servidor do assinador não estiver rodando, a CLI o iniciará em background automaticamente.

```bash
assinatura sign --file receita.json

# Execução local explícita
assinatura --local sign --file receita.json

# Servidor em porta e timeout personalizados
assinatura --port 9090 --timeout 10 sign --file receita.json
```

- **Retorno**: Mensagem confirmando a assinatura e o código gerado.
- **Arquivo**: Gera um arquivo `.txt` contendo a assinatura no mesmo diretório.

### 2. Validar uma assinatura

Valida se o conteúdo de um arquivo corresponde à assinatura registrada.

```bash
assinatura validate arquivo-assinatura.txt
```

- **Retorno**: Status de validade.

### 3. Parar o servidor

Encerra o processo do assinador que está rodando em background.

```bash
assinatura stop

# Para uma instância em outra porta
assinatura --port 9090 stop
```

### 4. Versão

Exibe a versão da CLI.

```bash
assinatura version
```

## Testes Automatizados

Para rodar os testes de integração (requer o Assinador disponível no GitHub ou localmente):

```bash
go test -v ./test/...
```

## Arquitetura de Integração

A CLI de Assinatura utiliza um contrato de comunicação estrito e bidirecional baseado em **JSON**.

- A CLI envia requisições em formato JSON (`{"command": "sign", "file": "..."}`) e recebe respostas em JSON. Isso se aplica tanto na comunicação via rede (Servidor HTTP na porta 8080) quanto na execução direta do `.jar` via subprocesso.
- O servidor Java atua puramente como um motor de dados de negócio. A CLI em Go lê os campos JSON de sucesso ou erro e aplica a formatação visual (Cores ANSI) para o usuário final.
- O Java classifica os erros estruturalmente no JSON (ex: `type: "user"` ou `type: "system"`). A CLI captura isso e traduz para Exit Codes de Unix padronizados:
    - `Exit Code 1`: Erros de validação ou entrada do usuário (ex: arquivo não encontrado).
    - `Exit Code 2`: Erros de execução do sistema.
- **Interoperabilidade**: Outras aplicações (como Frontends ou painéis Python) podem invocar a API HTTP local na porta 8080 recebendo dados limpos e perfeitamente tipados.

O contrato completo está documentado em [`docs/contrato-assinador.md`](../../../docs/contrato-assinador.md).
