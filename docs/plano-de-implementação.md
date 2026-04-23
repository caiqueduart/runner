# Plano de Implementação

## Fase 1: Fundamentação e Planejamento (04/03 - 15/04)

- Entendimento do problema e escopo do projeto, análise detalhada dos requisitos.
- Ambientação com ferramentas e tecnologias, configuração do ambiente de desenvolvimento para Java e Go.
- Estudo da biblioteca Cobra (Go), frameworks para servidores HTTP em Java e o ecossistema Sigstore/Cosign para assinatura de artefatos, e Github Actions para CI/CD.
- Elaboração do plano de implementação, estruturação das etapas de desenvolvimento.
- Inicialização do projeto em Go utilizando Cobra CLI.

</br>

## Fase 2: Implementação da assinatura e configuração de CI/CD (desde 20/04)

- Implementação da assinatura CLI para realização de comandos para interagir com `assinador.jar`.
- Pipeline de CI/CD, configuração do GitHub Actions para compilação para Windows, Linux e macOS (amd64).
- Automação da geração de tags e disponibilização de binários no GitHub Releases.

</br>

## Fase 3: Implementação do Assinador (03/06 - 17/06)

- Criação do projeto base em Java do assinador.jar com validação de parâmetros.
- Implementação de endpoints HTTP no `assinador.jar` para suportar o modo servidor.
- Lógica para detecção e download automático do JDK 21 caso não esteja presente no sistema.
- Gestão de Ciclo de Vida com comandos para iniciar, parar e monitorar o status do processo servidor.

</br>

## Fase 4: Simulador e Testes (17/06 - 24/06)

- Desenvolvimento do simulador CLI para gerenciar o ciclo de vida do `simulador.jar`.
- Implementação da rotina para baixar a versão mais recente do simulador via GitHub Releases.
- Implementação dos testes unitários e de integração para todos os CLI e jar.

</br>

## Fase 5: Integridade Final (até 24/06)

- Integração do Cosign no pipeline para assinar os binários finais e gerar arquivos `.sig` e `.pem` para assinatura de artefatos.
- Elaboração do manual do usuário, guia de instalação e exemplos de uso.
