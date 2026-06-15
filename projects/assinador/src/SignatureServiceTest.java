import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.junit.jupiter.api.BeforeEach;
import services.SignatureService;
import java.io.File;
import java.nio.file.Path;
import java.nio.file.Files;

class SignatureServiceTest {

    @TempDir
    Path tempDir;

    private File testFile;

    @BeforeEach
    void setUp() throws Exception {
        testFile = tempDir.resolve("test.txt").toFile();
        Files.writeString(testFile.toPath(), "conteudo de teste");
    }

    @Test
    void testSignSuccess() throws Exception {
        String result = SignatureService.sign(testFile.getAbsolutePath());
        assertNotNull(result);
        assertTrue(result.contains("gerou o código de assinatura"));
        
        // Verifica se o arquivo de assinatura foi criado
        String expectedSignFileName = "test-txt-assinatura.txt";
        File signFile = new File(System.getProperty("user.dir"), expectedSignFileName);
        assertTrue(signFile.exists(), "Arquivo de assinatura deve existir");
        
        // Cleanup
        signFile.delete();
    }

    @Test
    void testSignFileNotFound() {
        assertThrows(Exception.class, () -> {
            SignatureService.sign("arquivo_inexistente.txt");
        });
    }

    @Test
    void testValidateSuccess() throws Exception {
        SignatureService.sign(testFile.getAbsolutePath());
        String signFileName = "test-txt-assinatura.txt";
        File signFile = new File(System.getProperty("user.dir"), signFileName);
        
        try {
            String result = SignatureService.validate(signFile.getAbsolutePath());
            assertNotNull(result);
            assertTrue(result.contains("está assinado sob o código"));
        } finally {
            signFile.delete();
        }
    }

    @Test
    void testValidateFileNotFound() {
        assertThrows(Exception.class, () -> {
            SignatureService.validate("assinatura_inexistente.txt");
        });
    }

    @Test
    void testValidateWrongExtension() throws Exception {
        File wrongFile = tempDir.resolve("wrong.pdf").toFile();
        Files.writeString(wrongFile.toPath(), "not a signature");
        
        assertThrows(Exception.class, () -> {
            SignatureService.validate(wrongFile.getAbsolutePath());
        }, "Deve falhar se não for .txt");
    }
}
