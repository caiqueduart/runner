import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import services.HttpServerService;
import services.SignatureService;
import services.Tint;

public class Main {
    public static void main(String[] args) throws UnsupportedEncodingException {
        /* Configurar UTF-8 */
        System.setOut(new PrintStream(System.out, true, "UTF-8"));
        System.setErr(new PrintStream(System.err, true, "UTF-8"));

        if (args.length == 0) {
            Tint.logFeedback("ASSINATURA", "Erro: Nenhum comando fornecido. Use 'sign' ou 'validate'.");
            System.exit(1);
        }

        String cmd = args[0];
        String fileName = "";
        int port = 8080;
        long timeout = 5; // Default de 5 minutos conforme CLI Go

        // Parsing dos argumentos dependendo do comando
        for (int i = 1; i < args.length; i++) {
            String arg = args[i];

            if (arg.equals("--file")) {
                if (i + 1 < args.length)
                    fileName = args[++i];
            } else if (arg.equals("--port") && cmd.equals("server")) {
                if (i + 1 < args.length)
                    port = Integer.parseInt(args[++i]);
            } else if (arg.equals("--timeout") && cmd.equals("server")) {
                if (i + 1 < args.length)
                    timeout = Long.parseLong(args[++i]);
            } else if (arg.startsWith("-")) {
                // Validação de flags desconhecidas
                String suggestion = (arg.equals("-f") || arg.equals("--f")) ? " Você quis dizer '--file'?" : "";
                Tint.logFeedback("ASSINATURA", "Erro: Flag '" + arg + "' não reconhecida." + suggestion);
                System.exit(1);
            } else if (fileName.isEmpty() && cmd.equals("validate")) {
                // Aceita posicional apenas no validate
                fileName = arg;
            }
        }

        // Validação de obrigatoriedade
        if (cmd.equals("sign")) {
            boolean usedFileFlag = false;
            for (String a : args)
                if (a.equals("--file"))
                    usedFileFlag = true;

            if (!usedFileFlag || fileName.isEmpty()) {
                Tint.logFeedback("ASSINATURA",
                        "Erro do usuário: O parâmetro '--file' é obrigatório para o comando sign.");
                System.exit(1);
            }
        } else if (cmd.equals("validate") && fileName.isEmpty()) {
            Tint.logFeedback("ASSINATURA", "Erro do usuário: Forneça o caminho do arquivo para validação.");
            System.exit(1);
        }

        try {
            switch (cmd) {
                case "server":
                    HttpServerService.start(port, timeout);
                    break;

                case "sign":
                    System.out.println(SignatureService.sign(fileName));
                    break;

                case "validate":
                    System.out.println(SignatureService.validate(fileName));
                    break;

                default:
                    Tint.logFeedback("ASSINATURA", "Erro: Comando '" + cmd + "' não reconhecido.");
                    System.exit(1);
            }
        } catch (Exception e) {
            Tint.logFeedback("ASSINATURA", "Erro do sistema: " + e.getMessage());
            System.exit(2);
        }
    }
}
