import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import services.HttpServerService;
import services.SignatureService;
import services.SignatureResult;
import services.ValidationResult;
import services.ServiceFactory;

public class App {
    public static void main(String[] args) throws UnsupportedEncodingException {
        int status = run(args);

        if (status != 0) {
            System.exit(status);
        }
    }

    public static int run(String[] args) throws UnsupportedEncodingException {
        /* configurar UTF-8 */
        System.setOut(new PrintStream(System.out, true, "UTF-8"));
        System.setErr(new PrintStream(System.err, true, "UTF-8"));

        if (args.length == 0) {
            printErrorJson("Nenhum comando fornecido. Use 'server' ou passe um JSON estruturado.", "user", 400);
            return 1;
        }

        String firstArg = args[0];

        // se for o comando de servidor (mantém compatibilidade de args tradicionais)
        if (firstArg.equals("server")) {
            int port = 8080;
            long timeout = 5;

            for (int i = 1; i < args.length; i++) {
                if (args[i].equals("--port") && i + 1 < args.length) {
                    port = Integer.parseInt(args[++i]);
                } else if (args[i].equals("--timeout") && i + 1 < args.length) {
                    timeout = Long.parseLong(args[++i]);
                }
            }

            HttpServerService.start(port, timeout);
            return 0;
        }

        // se for um comando de negócio (sign/validate), deve ser um JSON
        if (!firstArg.trim().startsWith("{")) {
            printErrorJson("O contrato de comunicação exige um objeto JSON como argumento.", "user", 400);
            return 1;
        }

        String jsonPayload = firstArg;
        String cmd = extractFromJson(jsonPayload, "command");
        String fileName = extractFromJson(jsonPayload, "file");

        if (cmd == null || cmd.isEmpty()) {
            printErrorJson("O parâmetro 'command' no JSON é obrigatório.", "user", 400);
            return 1;
        }
        
        if (fileName == null || fileName.isEmpty()) {
            printErrorJson("O parâmetro 'file' no JSON é obrigatório.", "user", 400);
            return 1;
        }

        try {
            SignatureService service = ServiceFactory.getSignatureService();

            switch (cmd) {
                case "sign":
                    SignatureResult signRes = service.sign(fileName);
                    System.out.println("""
                        {
                            "message": "Arquivo assinado com sucesso.",
                            "fileName": "%s",
                            "code": "%s",
                            "signOutputPath": "%s",
                            "status": 200
                        }
                        """.formatted(
                            signRes.fileName().replace("\\", "/"), 
                            signRes.code(), 
                            signRes.filePath().replace("\\", "/")
                        ));
                    break;

                case "validate":
                    ValidationResult valRes = service.validate(fileName);
                    System.out.println("""
                        {
                            "message": "Validação concluída.",
                            "fileName": "%s",
                            "code": "%s",
                            "valid": %b,
                            "status": 200
                        }
                        """.formatted(
                            valRes.fileName().replace("\\", "/"), 
                            valRes.code(), 
                            valRes.valid()
                        ));
                    break;

                default:
                    printErrorJson("Comando '" + cmd + "' não reconhecido.", "user", 400);
                    return 1;
            }
        } catch (Exception e) {
            boolean isUserError = e.getMessage().contains("não encontrado") || e.getMessage().contains("obrigatório") || e.getMessage().contains(".txt");
            if (isUserError) {
                printErrorJson(e.getMessage(), "user", 400);
                return 1;
            } else {
                printErrorJson(e.getMessage(), "system", 500);
                return 2;
            }
        }

        return 0;
    }

    private static String extractFromJson(String json, String key) {
        String pattern = "\"" + key + "\":\\s*\"([^\"]*)\"";
        java.util.regex.Matcher matcher = java.util.regex.Pattern.compile(pattern).matcher(json);

        if (matcher.find()) {
            return matcher.group(1);
        }
        
        return null;
    }

    private static void printErrorJson(String message, String type, int status) {
        System.out.println("""
            {
                "error": "%s",
                "status": %d,
                "type": "%s"
            }
            """.formatted(message.replace("\"", "\\\""), status, type));
    }
}
