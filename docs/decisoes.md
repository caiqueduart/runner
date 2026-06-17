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

**Decisão:** Separação de UI e Contrato de Integração (JSON).

- **Diferença:** O servidor Java não envia strings via HTTP. Toda a comunicação entre a CLI e o servidor é feita através de objetos JSON.
- **Justificativa:** Melhora a interoperabilidade (outros clientes podem consumir a API). A responsabilidade de "como exibir os dados" passa a ser exclusivamente do cliente (CLI Go), enquanto o servidor foca apenas em "quais dados retornar".

---

**Decisão:** Uso de Java 21 Records e Text Blocks.

- **Diferença:** Uso de `record` para modelos de dados e `Text Blocks` para geração manual de JSON sem dependências externas.
- **Justificativa:** Aproveita as funcionalidades modernas do Java 21 para manter o código limpo, seguro e o artefato JAR extremamente leve, sem necessidade de bibliotecas como Jackson ou Gson.
