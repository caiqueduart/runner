package src;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import src.services.HttpServerService;
import src.services.SignatureService;
import src.services.Tint;

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

        // Parsing manual dos argumentos
        for (int i = 1; i < args.length; i++) {
            if (args[i].equals("--file")) {
                if (i + 1 < args.length) {
                    fileName = args[++i];
                }
            } else if (args[i].startsWith("-")) {
                String flag = args[i];
                String suggestion = flag.equals("-f") ? " Você quis dizer '--file'?" : "";
                Tint.logFeedback("ASSINATURA", "Erro: Flag '" + flag + "' não reconhecida." + suggestion);
                System.exit(1);
            } else if (fileName.isEmpty()) {
                // Se não é uma flag e fileName está vazio, assumimos como argumento posicional
                fileName = args[i];
            }
        }

        // Validação exclusiva no JAR
        if (cmd.equals("sign")) {
            /*  No caso do sign, o usuário deve ter usado --file (ou passamos via loop acima)
                Mas para garantir a regra, verificamos se o fileName foi preenchido.
                Se o usuário digitar 'sign arquivo.txt', o loop acima pegará o fileName.
                No entanto, se queremos OBRIGAR a flag --file no sign: */
            boolean usedFileFlag = false;
            for(String a : args) if(a.equals("--file")) usedFileFlag = true;

            if (!usedFileFlag || fileName.isEmpty()) {
                Tint.logFeedback("ASSINATURA", "Erro do usuário: O parâmetro '--file' é obrigatório para o comando sign.");
                System.exit(1);
            }
        } else if (cmd.equals("validate")) {
            if (fileName.isEmpty()) {
                Tint.logFeedback("ASSINATURA", "Erro do usuário: Forneça o caminho do arquivo para validação.");
                System.exit(1);
            }
        }
        try {
            switch (cmd) {
                case "server":
                    int port = 8080;
                    long timeout = 10;
                    for (int i = 1; i < args.length; i++) {
                        if (args[i].equals("--port") && i + 1 < args.length) port = Integer.parseInt(args[++i]);
                        if (args[i].equals("--timeout") && i + 1 < args.length) timeout = Long.parseLong(args[++i]);
                    }
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
