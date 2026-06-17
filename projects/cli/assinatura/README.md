# Assinatura CLI

A **CLI de Assinatura** é a interface principal do sistema Runner, permitindo que usuários interajam com o motor de assinatura sem conhecer o ecosistema Java.

## Instalação e Provisionamento

A CLI gerencia automaticamente suas dependências. Na primeira execução, ela baixará o JDK 21 e o `assinador.jar` compatível para a pasta `~/.hubsaude`.

## Comandos Principais

### 1. Assinar um arquivo

Realiza a assinatura de um documento. Se o servidor do assinador não estiver rodando, a CLI o iniciará em background automaticamente.

```bash
assinatura sign --file receita.json
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

## Observações

- **Execução**: Nos exemplos acima, o comando `assinatura` refere-se ao nome do arquivo binário baixado. Por exemplo, no Windows, o uso real seria: `.\assinatura-cli-v1.1.3-windows-amd64.exe sign --file receita.json`.
- **Modo Servidor**: A CLI prioriza a comunicação via HTTP (porta 8080) para maior performance em execuções repetitivas.
- **Rastreabilidade**: O PID do servidor em execução é armazenado em `~/.hubsaude/assinador.pid`.
