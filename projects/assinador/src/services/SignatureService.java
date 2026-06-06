package src.services;

public class SignatureService {
    SignatureService() {}

    public static String sign(String fileName) {
        String code = makeSimulatedSignCode(fileName);
        return formatSignatureMessage(fileName, code);
    }

    public static String validate(String fileName) {
        String code = makeSimulatedSignCode(fileName);
        return formatValidationMessage(fileName, code);
    }

    public static String formatSignatureMessage(String fileName, String signatureCode) {
        return "O Arquivo '" + Tint.GREEN + fileName + Tint.RESET + "' gerou o código de assinatura '" + Tint.GREEN + signatureCode + Tint.RESET + "'.";
    }

    public static String formatValidationMessage(String fileName, String signatureCode) {
        return "O Arquivo '" + Tint.GREEN + fileName + Tint.RESET + "' está assinado sob o código '" + Tint.GREEN + signatureCode + Tint.RESET + "'.";
    }

    public static String makeSimulatedSignCode(String fileName) {
        return "SIGN-" + Math.abs(fileName.hashCode());
    }
}
