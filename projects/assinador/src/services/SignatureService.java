package src.services;

public class SignatureService {
    SignatureService() {}

    public static void sign(String fileName) {
        String simulatetdSign = SignatureService._makeSimulatedSignCode(fileName);

        System.out.println(
            "O Arquivo \'" + Tint.GREEN + fileName + Tint.RESET + "\' " + 
            "gerou o código de assinatura" 
            + " \'" + Tint.GREEN + simulatetdSign + Tint.RESET + "\'."
        );
    }

    public static void validate(String fileName) {
        String simulatetdSign = SignatureService._makeSimulatedSignCode(fileName);

        System.out.println(
            "O Arquivo \'" + Tint.GREEN + fileName + Tint.RESET + "\' " + 
            "está assinado sob o código" 
            + " \'" + Tint.GREEN + simulatetdSign + Tint.RESET + "\'."
        );
    }

    private static String _makeSimulatedSignCode(String fileName) {
        return "SIGN-" + Math.abs(fileName.hashCode());
    }
}
