package src;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;
import src.services.Tint;

public class Main {
    public static void main(String[] args) throws UnsupportedEncodingException {
        System.setOut(new PrintStream(System.out, true, "UTF-8"));

        String file = args[0];

        String assinaturaSimulada = "SIMULATED_SIG_" + file.hashCode();
        System.out.println(
            "\n\"" + Tint.GREEN + file + Tint.RESET + "\"" + 
            " gerou código o de assinatura " 
            + "\"" + Tint.GREEN + assinaturaSimulada + Tint.RESET + "\"."
        );
    }
}