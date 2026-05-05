package src;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import src.services.SignatureService;

public class Main {
    public static void main(String[] args) throws UnsupportedEncodingException {
        /** Configurar UTF-8 */
        System.setOut(new PrintStream(System.out, true, "UTF-8"));

        SignatureService.sign(args[0]);
    }
}