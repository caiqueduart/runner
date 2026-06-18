# Plano de Implementação

Este plano detalha as etapas de desenvolvimento do sistema **Runner**, integrando uma CLI em **Go** com um serviço de assinatura em **Java 21**, conforme as [especificações oficiais](https://github.com/kyriosdata/runner/blob/main/especificacao.md).

## Cronograma de Sprints

### Sprint 1: Fundamentação e Estudo

- [x] Estudo das linguagens: Java 21 (novas features) e Go (Cobra CLI).
- [x] Configuração dos repositórios e estrutura de pastas (`projects/assinador` e `projects/cli`).
- [x] Definição da estratégia de comunicação inicial (Modo Local).
- [x] Configuração inicial de CI/CD com GitHub Actions para releases básicas.

### Sprint 2: Core e Provisionamento

- [x] **Assinador (Java):** Implementação da lógica de simulação de assinatura e validação.
- [x] **CLI (Go):** Implementação dos comandos `sign`, `validate` e `version`.
- [x] **Auto-Provisionamento:** Lógica para baixar e configurar o JDK 21 automaticamente.
- [x] **Gestão de Artefatos:** Download automático do `assinador.jar` via GitHub Releases.
- [x] **Validação de Integridade:** Implementar verificação de SHA256 para arquivos baixados.

### Sprint 3: Modo Servidor

- [x] **Assinador Server:** Adicionar servidor HTTP ao projeto Java.
- [x] **Endpoints:** Criar rotas `POST /sign` e `POST /validate`.
- [x] **CLI Server Mode:**
    - [x] Lógica para detectar servidor ativo.
    - [x] Inicialização em background.
    - [x] Cliente HTTP em Go para comunicação com o servidor.
- [x] **Comando `stop`:** Implementar encerramento gracioso do servidor.
- [x] **Inatividade:** Implementar `--timeout` para auto-desligamento do Java.

### Sprint 4: Simulador e Integração PKCS#11

- [x] **Gestão do Simulador:** CLI deve baixar e gerenciar o ciclo de vida do `simulador.jar` (Implementado via `simulador-cli`).
- [ ] **Integração PKCS#11:** Foco redirecionado para a gestão de infraestrutura e prontidão para drivers (US-03).
- [x] **Configuração Dinâmica:** Passagem de parâmetros de configuração integrada ao fluxo de provisionamento.

### Sprint 5: Segurança Avançada e Refinamento

- [ ] **Assinaturas Digitais:** Integrar verificação de assinaturas **Cosign/Sigstore** na CLI.
- [ ] **Logs e Debug:** Implementar flags de verbosidade (`--verbose`) e logs estruturados.
- [x] **Tratamento de Erros:** Refinar mensagens de erro para serem aderentes à especificação.
- [x] **Documentação:** Atualizar README e documentação de uso da CLI.

### Sprint 6: Validação e Entrega

- [x] **Testes de Integração:** Suite de testes automatizados CLI <-> Java.
- [x] **Testes de Aceitação:** Validar todos os cenários descritos na especificação.
- [x] **Release Final:** Gerar versão final com todos os artefatos assinados.
