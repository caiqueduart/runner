package services;

import java.io.File;
import java.io.FileWriter;
import java.nio.file.Files;
import java.security.MessageDigest;
import java.util.Scanner;

public class SignatureService {

    public SignatureResult sign(String fileName) throws Exception {
        File file = new File(fileName);

        if (!file.exists()) {
            throw new Exception("Arquivo '" + fileName + "' não encontrado.");
        }

        String code = generateContentHash(file);

        String nameOnly = file.getName();
        String baseName = nameOnly.contains(".") ? nameOnly.substring(0, nameOnly.lastIndexOf('.')) : nameOnly;
        String extension = nameOnly.contains(".") ? nameOnly.substring(nameOnly.lastIndexOf('.') + 1) : "";
        String outputName = baseName + "-" + extension + "-assinatura.txt";

        File outputFile = new File(System.getProperty("user.dir"), outputName);

        try (FileWriter writer = new FileWriter(outputFile)) {
            writer.write(code);
        }

        return new SignatureResult(fileName, code, outputFile.getAbsolutePath());
    }

    public ValidationResult validate(String fileName) throws Exception {
        File file = new File(fileName);

        if (!file.exists()) {
            throw new Exception("Arquivo de assinatura '" + fileName + "' não encontrado.");
        }

        if (!fileName.endsWith(".txt")) {
            throw new Exception("O arquivo de validação deve ser um .txt gerado pela assinatura.");
        }

        String code;

        try (Scanner scanner = new Scanner(file)) {
            if (!scanner.hasNext()) throw new Exception("Arquivo de assinatura vazio.");
            code = scanner.next();
        }

        return new ValidationResult(fileName, code, true);
    }

    private String generateContentHash(File file) throws Exception {
        byte[] content = Files.readAllBytes(file.toPath());
        MessageDigest digest = MessageDigest.getInstance("SHA-256");
        byte[] hash = digest.digest(content);
        
        StringBuilder hexString = new StringBuilder();

        for (byte b : hash) {
            String hex = Integer.toHexString(0xff & b);
            if (hex.length() == 1) hexString.append('0');
            hexString.append(hex);
        }
        
        return hexString.toString().substring(0, 12).toUpperCase(); // 12 chars prefix for simulation
    }
}
