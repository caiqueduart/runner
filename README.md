# Runner

Runner é um sistema para validação de assinaturas digitais, desenvolvido para a disciplina de Implementação e Integração de Software. O sistema é composto por um serviço core em Java e duas interfaces de linha de comando (CLI) em Go.

## Componentes do Projeto

O projeto está dividido em três componentes principais:

1.  **[Assinador (Java)](./projects/assinador/README.md)**: O motor de assinatura e servidor HTTP que realiza o processamento de assinaturas.
2.  **[CLI de Assinatura (Go)](./projects/cli/assinatura/README.md)**: Interface principal para usuários realizarem assinaturas e validações de arquivos.
3.  **[CLI do Simulador (Go)](./projects/cli/simulador/README.md)**: Ferramenta de gestão do ciclo de vida do Simulador HubSaúde (US-03).

## Referências

- **Especificações oficiais**: [github.com/kyriosdata/runner](https://github.com/kyriosdata/runner)
- **Plano de implementação**: [plano-de-implementacao.md](./docs/plano-de-implementacao.md)
- **Decisões de Projeto**: [decisoes.md](./docs/decisoes.md)

## Como Começar

As CLIs do Runner possuem **auto-provisionamento**. Isso significa que você não precisa instalar o Java manualmente.

1.  **Baixe o executável** da CLI desejada na seção de [Releases](https://github.com/caiqueduart/runner/releases).
2.  **Execute um comando**: Ao rodar algum comando de Assinatura ou Simulador pela primeira vez, as CLIs irão:
    - Baixar o JDK 21 automaticamente para uma pasta local (`.hubsaude/jdk`).
    - Baixar os arquivos JAR necessários.
    - Configurar o ambiente de execução.

---

### Estrutura de diretórios

```
runner/
├── .github/workflows/          # CI/CD (Releases automatizadas)
├── docs/                       # Documentação técnica e plano de sprints
└── projects/
    ├── assinador/              # Core em Java (Servidor de Assinatura)
    └── cli/
        ├── assinatura/         # CLI principal (Usuário Final)
        └── simulador/          # CLI de infraestrutura (Ambiente de Teste)
```
