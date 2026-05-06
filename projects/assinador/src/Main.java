package src;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import src.services.SignatureService;
import src.services.Tint;

public class Main {
    public static void main(String[] args) throws UnsupportedEncodingException {
        /* Configurar UTF-8 */
        System.setOut(new PrintStream(System.out, true, "UTF-8"));
        
        if (args.length > 0) {

            String cmd = args[0];
            String file = args[1];
            
            switch (cmd) {
                case "sign":
                    SignatureService.sign(file);
                    break;

                case "validate":
                    SignatureService.validate(file);
                    break;
                    
                default:
                    System.out.print(Tint.RED + "Erro: O comando passado para o Assinador não foi reconhecido." + Tint.RESET);   
            }
           
        } else {
            System.out.print(Tint.RED + "Erro: Nenhum argumento foi passado para o Assinador." + Tint.RESET);
        }
    }
}