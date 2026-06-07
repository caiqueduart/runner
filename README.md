# Runner

Runner é um projeto CLI para validação de assinaturas digitais, desenvolvido para a disciplina de Implementação e Integração de Software.

## Sobre

**Runner** é uma aplicação de linha de comando (CLI) chamada **Assinatura** que se comunica com uma aplicação Java `Assinador.jar` para realizar validação de assinaturas digitais.

## Referências

- **Especificações completas**: [github.com/kyriosdata/runner](https://github.com/kyriosdata/runner)
- **Plano de implementação**: [plano-de-implementação.md](./docs/plano-de-implementacao.md)
- **Decisões de Projeto**: [decisoes.md](./docs/decisoes.md)

## Como Usar

Documentação em desenvolvimento...

### Estrutura de diretórios

```
runner/
├── .github/
│   └── workflows/                  # Workflows do GitHub Actions (CI/CD)
│
├── docs/                           # Documentações
│
└── projects/
    ├── assinador/                  # Código-fonte da aplicação Assinador (Java)
    │   │
    │   ├── src/
    │   │   ├── services/           # Lógica de validação e simulação
    │   │   └── Main.java
    │   └── pom.xml                 # Arquivo de configuração do Maven
    │
    └── cli/                        # Código-fonte das aplicações Go (CLIs)
        ├── assinatura/             # CLI de assinatura
        │   ├── cmd/                # Comandos disponibilizados pelo CLI
        │   ├── internal/           # Códigos privados da cli
        │   └── main.go
        │
        ├── simulador/              # CLI do simulador
        │   ├── cmd/                # Comandos disponibilizados pelo CLI
        │   ├── internal/           # Códigos privados da cli
        │   └── main.go
        │
        ├── go.mod                  # Definição do módulo Go principal
        └── go.sum                  # Checksums das dependências Go
```
