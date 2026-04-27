# Runner

Runner é um projeto CLI para validação de assinaturas digitais, desenvolvido para a disciplina de Implementação e Integração de Software.

## Sobre

**Runner** é uma aplicação de linha de comando (CLI) chamada **Assinatura** que se comunica com uma aplicação Java `Assinador.jar` para realizar validação de assinaturas digitais.

## Referências

- **Especificações completas**: [github.com/kyriosdata/runner](https://github.com/kyriosdata/runner)
- **Plano de implementação**: [plano-de-implementação.md](./docs/plano-de-implementação.md)

## Como Usar

Documentação em desenvolvimento...

### Estrutura de diretórios

```
runner/
├── .github/
│   └── workflows/              # Workflows do GitHub Actions (CI/CD)
│
├── docs/                       # Documentações
│
└── projects/
    ├── cli/                    # Código-fonte das aplicações Go (CLIs)
    │   ├── assinatura/         # CLI de assinatura
│   │   │   ├── cmd/            # Comandos disponibilizados pelo CLI
    │   │   └── main.go
    │   │
    │   ├── simulador/          # CLI do simulador
    │   │   ├── cmd/            # Comandos disponibilizados pelo CLI
    │   │   └── main.go
    │   │
    │   ├── go.mod              # Definição do módulo Go principal
    │   └── go.sum              # Checksums das dependências Go
    │
    └── assinador-java/         # Código-fonte da aplicação Assinador (Java)
       ├── src/
       │   ├── main/java/       # Lógica de validação e simulação
       │   └── test/java/       # Testes unitários e de integração
       ├── pom.xml              # Arquivo de configuração do Maven
       └── README.md            # Instruções específicas do projeto Java
```
