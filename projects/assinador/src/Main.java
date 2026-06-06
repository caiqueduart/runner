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
        
        if (args.length > 0) {

            String cmd = args[0];
            
            switch (cmd) {
                case "server":
                    int port = 8080;
                    long timeout = 10;
                    
                    for (int i = 1; i < args.length; i++) {
                        if (args[i].equals("--port") && i + 1 < args.length) {
                            port = Integer.parseInt(args[++i]);
                        } else if (args[i].equals("--timeout") && i + 1 < args.length) {
                            timeout = Long.parseLong(args[++i]);
                        }
                    }
                    HttpServerService.start(port, timeout);
                    break;

                case "sign":
                    if (args.length < 2) {
                        System.out.print(Tint.RED + "Erro: Nome do arquivo não fornecido para o comando sign." + Tint.RESET);
                        return;
                    }
                    System.out.println(SignatureService.sign(args[1]));
                    break;

                case "validate":
                    if (args.length < 2) {
                        System.out.print(Tint.RED + "Erro: Nome do arquivo não fornecido para o comando validate." + Tint.RESET);
                        return;
                    }
                    System.out.println(SignatureService.validate(args[1]));
                    break;

                default:
                    String formatedArgs = String.join(" ", args);
                    System.out.print(Tint.RED + "Erro: O argumento '" + formatedArgs + "' não foi reconhecido pelo Assinador." + Tint.RESET);   
            }
           
        } else {
            System.out.print(Tint.RED + "Erro: Nenhum argumento foi passado para o Assinador." + Tint.RESET);
        }
    }
}