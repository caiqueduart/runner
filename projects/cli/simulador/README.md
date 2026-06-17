# Simulador CLI

A **CLI do Simulador** é uma ferramenta dedicada à gestão do ambiente **HubSaúde**. Ela facilita o download e gestão do ciclo de vida do simulador.

## Funcionalidades

- **Provisionamento Automático**: Baixa o JDK 21 e o `simulador.jar` (v0.1.7) automaticamente.
- **Gestão via HTTPS**: Comunica-se com o simulador através do protocolo seguro na porta 8443.
- **Verificação de Porta**: Garante que a porta 8443 esteja livre antes de tentar iniciar o serviço.

## Comandos

### 1. Iniciar o Simulador

Inicia o simulador em background.

```bash
simulador start
```

- **Obs**: Verifica se o simulador já está rodando ou se a porta está ocupada por outro processo.

### 2. Verificar Status

Verifica se o simulador está online e exibe informações de versão da API.

```bash
simulador status
```

- **Retorno**: Detalhes em JSON retornados pelo simulador (ex: `{"version":"0.1.7","name":"HubSaúde Simulador"}`).

### 3. Parar o Simulador

Encerra o simulador usando o endpoint `/shutdown`.

```bash
simulador stop
```

### 4. Versão

Exibe a versão da CLI e a versão do JAR compatível.

```bash
simulador version
```

## Observações Técnicas

- **Execução**: Nos exemplos acima, o comando `simulador` refere-se ao nome do arquivo binário baixado. Por exemplo, no Windows, o uso real seria: `.\simulador-cli-v1.0.0-windows-amd64.exe start`.
