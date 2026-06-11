package services;

public class SignatureService {
    SignatureService() {}

    public static String sign(String fileName) throws Exception {
        java.io.File file = new java.io.File(fileName);
        if (!file.exists()) {
            throw new Exception("Arquivo '" + fileName + "' não encontrado.");
        }

        String code = makeSimulatedSignCode(fileName);

        String nameOnly = file.getName();
        String baseName = nameOnly.contains(".") ? nameOnly.substring(0, nameOnly.lastIndexOf('.')) : nameOnly;
        String extension = nameOnly.contains(".") ? nameOnly.substring(nameOnly.lastIndexOf('.') + 1) : "";
        String outputName = baseName + "-" + extension + "-assinatura.txt";

        // Usamos o diretório de execução atual do processo Java
        java.io.File outputFile = new java.io.File(System.getProperty("user.dir"), outputName);
        try (java.io.FileWriter writer = new java.io.FileWriter(outputFile)) {
            writer.write(code);
        }

        return formatSignatureMessage(fileName, code) + "\n" +
               Tint.CYAN + "[ASSINATURA] " + Tint.RESET + "Arquivo gerado em: " + Tint.GREEN + outputFile.getAbsolutePath() + Tint.RESET;
    }


    public static String validate(String fileName) throws Exception {
        java.io.File file = new java.io.File(fileName);
        if (!file.exists()) {
            throw new Exception("Arquivo de assinatura '" + fileName + "' não encontrado.");
        }

        if (!fileName.endsWith(".txt")) {
            throw new Exception("O arquivo de validação deve ser um .txt gerado pela assinatura.");
        }

        String code;
        try (java.util.Scanner scanner = new java.util.Scanner(file)) {
            if (!scanner.hasNext()) throw new Exception("Arquivo de assinatura vazio.");
            code = scanner.next();
        }

        return formatValidationMessage(fileName, code);
    }

    public static String formatSignatureMessage(String fileName, String signatureCode) {
        return Tint.CYAN + "[ASSINATURA] " + Tint.RESET + "O Arquivo '" + Tint.GREEN + fileName + Tint.RESET + "' gerou o código de assinatura '" + Tint.GREEN + signatureCode + Tint.RESET + "'.";
    }

    public static String formatValidationMessage(String fileName, String signatureCode) {
        return Tint.CYAN + "[ASSINATURA] " + Tint.RESET + "O Arquivo '" + Tint.GREEN + fileName + Tint.RESET + "' está assinado sob o código '" + Tint.GREEN + signatureCode + Tint.RESET + "'.";
    }

    public static String makeSimulatedSignCode(String fileName) {
        return "SIGN-" + Math.abs(fileName.hashCode());
    }
}
