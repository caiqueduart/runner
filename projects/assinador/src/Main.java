package src;
import java.io.PrintStream;
import java.io.UnsupportedEncodingException;

public class Main {
    // Definindo as cores como constantes para facilitar o uso
    private static String _ANSI_RESET = "\u001B[0m";
    private static String _ANSI_RED = "\u001B[31m";
    private static String _ANSI_GREEN = "\u001B[32m";
    private static String _ANSI_BLUE = "\u001B[34m";
    private static String _ANSI_YELLOW = "\033[33m";

    public static void main(String[] args) throws UnsupportedEncodingException {
        System.setOut(new PrintStream(System.out, true, "UTF-8"));

        String file = args[0];

        String assinaturaSimulada = "SIMULATED_SIG_" + file.hashCode();
        System.out.println(
            "\n\"" + _ANSI_GREEN + file + _ANSI_RESET + "\"" + 
            " gerou código o de assinatura " 
            + "\"" + _ANSI_GREEN + assinaturaSimulada + _ANSI_RESET + "\"."
        );

    }
}