package src.services;

public class SignatureService {
    SignatureService() {}

    public static void sign(String fileName) {
        String assinaturaSimulada = "SIGN-" + Math.abs(fileName.hashCode());

        System.out.println(
            "O Arquivo \'" + Tint.GREEN + fileName + Tint.RESET + "\' " + 
            "gerou o código de assinatura" 
            + " \'" + Tint.GREEN + assinaturaSimulada + Tint.RESET + "\'."
        );
    }

    public static void validate(String fileName) {}
}
