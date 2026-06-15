package test.services;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import services.SignatureService;

import java.io.File;
import java.io.FileWriter;
import java.nio.file.Path;

import static org.junit.jupiter.api.Assertions.*;

class SignatureServiceTest {

    @TempDir
    Path tempDir;

    private File testFile;

    @BeforeEach
    void setUp() throws Exception {
        testFile = tempDir.resolve("test.txt").toFile();
        try (FileWriter writer = new FileWriter(testFile)) {
            writer.write("conteudo de teste");
        }
    }

    @Test
    void testSignSuccessful() throws Exception {
        String result = SignatureService.sign(testFile.getAbsolutePath());
        
        assertNotNull(result);
        assertTrue(result.contains("gerou o código de assinatura"));
        
        // Check if signature file was created
        String nameOnly = testFile.getName();
        String baseName = nameOnly.substring(0, nameOnly.lastIndexOf('.'));
        String extension = nameOnly.substring(nameOnly.lastIndexOf('.') + 1);
        String outputName = baseName + "-" + extension + "-assinatura.txt";
        
        File signatureFile = new File(System.getProperty("user.dir"), outputName);
        assertTrue(signatureFile.exists(), "Signature file should exist: " + signatureFile.getAbsolutePath());
        
        // Cleanup signature file
        signatureFile.delete();
    }

    @Test
    void testSignFileNotFound() {
        Exception exception = assertThrows(Exception.class, () -> {
            SignatureService.sign("non_existent_file.txt");
        });
        assertTrue(exception.getMessage().contains("não encontrado"));
    }

    @Test
    void testValidateSuccessful() throws Exception {
        // First sign to get a valid signature file
        SignatureService.sign(testFile.getAbsolutePath());
        
        String nameOnly = testFile.getName();
        String baseName = nameOnly.substring(0, nameOnly.lastIndexOf('.'));
        String extension = nameOnly.substring(nameOnly.lastIndexOf('.') + 1);
        String outputName = baseName + "-" + extension + "-assinatura.txt";
        File signatureFile = new File(System.getProperty("user.dir"), outputName);

        String result = SignatureService.validate(signatureFile.getAbsolutePath());
        
        assertNotNull(result);
        assertTrue(result.contains("está assinado sob o código"));
        
        // Cleanup
        signatureFile.delete();
    }

    @Test
    void testValidateFileNotFound() {
        Exception exception = assertThrows(Exception.class, () -> {
            SignatureService.validate("non_existent_signature.txt");
        });
        assertTrue(exception.getMessage().contains("não encontrado"));
    }

    @Test
    void testValidateInvalidExtension() throws Exception {
        File invalidFile = tempDir.resolve("invalid.pdf").toFile();
        invalidFile.createNewFile();
        
        Exception exception = assertThrows(Exception.class, () -> {
            SignatureService.validate(invalidFile.getAbsolutePath());
        });
        assertTrue(exception.getMessage().contains("deve ser um .txt"));
    }

    @Test
    void testMakeSimulatedSignCode() {
        String code1 = SignatureService.makeSimulatedSignCode("file1.txt");
        String code2 = SignatureService.makeSimulatedSignCode("file1.txt");
        String code3 = SignatureService.makeSimulatedSignCode("file2.txt");
        
        assertEquals(code1, code2);
        assertNotEquals(code1, code3);
        assertTrue(code1.startsWith("SIGN-"));
    }
}
