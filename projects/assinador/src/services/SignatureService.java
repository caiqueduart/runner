package src.services;

public class SignatureService {
    SignatureService() {}

    public static void sign(String file) {
        String assinaturaSimulada = "SIMULATED_SIG_" + file.hashCode();

        System.out.println(
            "\n\"" + Tint.GREEN + file + Tint.RESET + "\"" + 
            " gerou código o de assinatura " 
            + "\"" + Tint.GREEN + assinaturaSimulada + Tint.RESET + "\"."
        );
    }

    public static void validate(String file) {}
}
