import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.junit.jupiter.api.BeforeEach;
import services.SignatureService;
import services.SignatureResult;
import services.ValidationResult;
import services.ServiceFactory;
import java.io.File;
import java.nio.file.Path;
import java.nio.file.Files;

class SignatureServiceTest {

    @TempDir
    Path tempDir;

    private File testFile;
    private SignatureService service;

    @BeforeEach
    void setUp() throws Exception {
        testFile = tempDir.resolve("test.txt").toFile();
        Files.writeString(testFile.toPath(), "conteudo de teste");
        service = ServiceFactory.getSignatureService();
    }

    @Test
    void testSignSuccess() throws Exception {
        SignatureResult result = service.sign(testFile.getAbsolutePath());
        assertNotNull(result);
        assertNotNull(result.code());
        
        // Verifica se o arquivo de assinatura foi criado
        File signFile = new File(result.filePath());
        assertTrue(signFile.exists(), "Arquivo de assinatura deve existir");
        
        // Cleanup
        signFile.delete();
    }

    @Test
    void testSignFileNotFound() {
        assertThrows(Exception.class, () -> {
            service.sign("arquivo_inexistente.txt");
        });
    }

    @Test
    void testValidateSuccess() throws Exception {
        SignatureResult signResult = service.sign(testFile.getAbsolutePath());
        File signFile = new File(signResult.filePath());
        
        try {
            ValidationResult result = service.validate(signFile.getAbsolutePath());
            assertNotNull(result);
            assertTrue(result.valid());
            assertEquals(signResult.code(), result.code());
        } finally {
            signFile.delete();
        }
    }

    @Test
    void testValidateFileNotFound() {
        assertThrows(Exception.class, () -> {
            service.validate("assinatura_inexistente.txt");
        });
    }

    @Test
    void testValidateWrongExtension() throws Exception {
        File wrongFile = tempDir.resolve("wrong.pdf").toFile();
        Files.writeString(wrongFile.toPath(), "not a signature");
        
        assertThrows(Exception.class, () -> {
            service.validate(wrongFile.getAbsolutePath());
        }, "Deve falhar se não for .txt");
    }
}
