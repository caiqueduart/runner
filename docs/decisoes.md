# Decisões de Projeto

Este documento registra as decisões arquiteturais e técnicas tomadas durante o desenvolvimento do sistema **Runner**, comparando-as com a especificação original e o [repositório de referência](https://github.com/kyriosdata/runner) do professor Dr. Fábio Nogueira.

---

**Decisão:** Uso de Servidor HTTP com.sun.net.httpserver.

- **Diferença:** Enquanto a especificação **US-02.4** exige os endpoints `POST /sign` e `POST /validate`, optei por não utilizar frameworks pesados (como Spring Boot), mantendo o uso de bibliotecas nativas do JDK.
- **Justificativa:** Garante um artefato (`assinador.jar`) mais leve.

---

**Decisão:** Sistema de Logs com Prefixos Semânticos Unificados. (extra)

- **Diferença:** Implementei um sistema de feedback com prefixos `[ASSINATURA]`, `[ASSINATURA CONFIG]` e `[ASSINATURA SERVIDOR]`, seguindo a recomendação de "mensagens úteis" da seção 8.2.
- **Justificativa:** Os prefixos permitem ao desenvolvedor e usuário distinguir instantaneamente a origem da mensagem.

---

**Decisão:** Ciclo de Vida Otimizado. (extra)

- **Diferença:** O servidor possui um cronômetro de auto-desligamento.
- **Justificativa:** Permite que o próprio usuário consulte o status do servidor em `http://localhost:8080/health`.

---

**Decisão:** Modularização da CLI Go. (extra)

- **Diferença:** O código Go das CLIs foi dividido em arquivos especializados (`assinador_manager.go`, `constants.go`, `utils.go`, etc) em vez de manter toda a lógica em um único arquivo.
- **Justificativa:** Melhora a manutenção a longo prazo e a legibilidade do projeto, separando as configurações globais das demais operações, e da lógica de gerenciamento de processos.

---

**Decisão:** Modo Desenvolvedor (DEV_MODE).

- **Diferença:** Implementação de uma variável de ambiente `DEV_MODE=true` em um arquivo `.env` que a CLI Assinatura usa para executar o código-fonte Java diretamente via `java -cp ... App.java`, em vez de baixar e executar o `.jar` da release.
- **Justificativa:** Agiliza o ciclo de desenvolvimento e testes de integração, permitindo validar mudanças no código Java instantaneamente sem a necessidade de gerar um novo artefato ou realizar um push para o GitHub.

---

**Decisão:** Uso de Java 21 e Virtual Threads (Project Loom).

- **Diferença:** O servidor Java utiliza `Executors.newVirtualThreadPerTaskExecutor()` para processar requisições HTTP.
- **Justificativa:** Permite que o servidor lide com alta concorrência de forma extremamente leve, sem o overhead de pools de threads tradicionais, aproveitando as funcionalidades mais modernas do JDK 21 para performance.

---

**Decisão:** Verificação de Integridade via SHA256.

- **Diferença:** Todos os downloads de artefatos realizados pela CLI são validados com um hash SHA256, e todos os arterfatos enviados ao repositório são assinados automaticamente pelo GitHub.
- **Justificativa:** Segurança e robustez. Garante que o binário baixado não foi corrompido ou alterado maliciosamente durante o transporte.
