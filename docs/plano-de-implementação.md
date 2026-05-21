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
- [ ] **Validação de Integridade:** Implementar verificação de SHA256 para arquivos baixados.

### Sprint 3: Modo Servidor e Performance

- [ ] **Assinador Server:** Adicionar servidor HTTP ao projeto Java.
- [ ] **Endpoints:** Criar rotas `POST /sign` e `POST /validate`.
- [ ] **CLI Server Mode:**
    - [ ] Lógica para detectar servidor ativo.
    - [ ] Inicialização em background.
    - [ ] Cliente HTTP em Go para comunicação com o servidor.
- [ ] **Comando `stop`:** Implementar encerramento gracioso do servidor.
- [ ] **Inatividade:** Implementar `--timeout` para auto-desligamento do Java.

### Sprint 4: Simulador e Integração PKCS#11

- [ ] **Gestão do Simulador:** CLI deve baixar e gerenciar o ciclo de vida do `simulador.jar`.
- [ ] **Integração PKCS#11:** Assinador configurado para usar o driver do simulador.
- [ ] **Configuração Dinâmica:** Passagem de parâmetros de configuração do token via CLI.

### Sprint 5: Segurança Avançada e Refinamento

- [ ] **Assinaturas Digitais:** Integrar verificação de assinaturas **Cosign/Sigstore** na CLI.
- [ ] **Logs e Debug:** Implementar flags de verbosidade (`--verbose`) e logs estruturados.
- [ ] **Tratamento de Erros:** Refinar mensagens de erro para serem 100% aderentes à especificação.
- [ ] **Documentação:** Atualizar README e documentação de uso da CLI.

### Sprint 6: Validação e Entrega

- [ ] **Testes de Integração:** Suite de testes automatizados CLI <-> Java.
- [ ] **Testes de Aceitação:** Validar todos os cenários descritos na especificação.
- [ ] **Release Final:** Gerar versão final com todos os artefatos assinados.
