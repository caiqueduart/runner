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
            
            if (cmd.contains("sign")) {
                SignatureService.sign(file);
            }
            
            if (cmd.contains("verify")) {
                SignatureService.validate(file);
            }

        } else {
            System.out.print(
                Tint.RED + "Erro: Nenhum argumento foi passado para o Assinador." + Tint.RESET);
        }
    }
}