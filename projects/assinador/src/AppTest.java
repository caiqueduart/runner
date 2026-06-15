import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.nio.file.Path;
import java.nio.file.Files;
import java.io.File;

class AppTest {

    @TempDir
    Path tempDir;

    @Test
    void testRunNoArgs() throws Exception {
        int status = App.run(new String[]{});
        assertEquals(1, status, "Deve retornar erro se não houver argumentos");
    }

    @Test
    void testRunUnknownCommand() throws Exception {
        int status = App.run(new String[]{"invalid"});
        assertEquals(1, status, "Deve retornar erro para comando desconhecido");
    }

    @Test
    void testRunSignMissingFile() throws Exception {
        int status = App.run(new String[]{"sign"});
        assertEquals(1, status, "Deve retornar erro se --file estiver ausente no sign");
    }

    @Test
    void testRunSignSuccess() throws Exception {
        Path testFile = tempDir.resolve("app-test.txt");
        Files.writeString(testFile, "app test content");

        int status = App.run(new String[]{"sign", "--file", testFile.toAbsolutePath().toString()});
        assertEquals(0, status, "Deve retornar sucesso para comando sign válido");

        // Cleanup
        String signFileName = "app-test-txt-assinatura.txt";
        File signFile = new File(System.getProperty("user.dir"), signFileName);
        if (signFile.exists()) signFile.delete();
    }

    @Test
    void testRunValidateSuccess() throws Exception {
        Path testFile = tempDir.resolve("app-validate.txt");
        Files.writeString(testFile, "app validate content");

        // Gera a assinatura primeiro
        App.run(new String[]{"sign", "--file", testFile.toAbsolutePath().toString()});
        
        String signFileName = "app-validate-txt-assinatura.txt";
        File signFile = new File(System.getProperty("user.dir"), signFileName);

        try {
            int status = App.run(new String[]{"validate", signFile.getAbsolutePath()});
            assertEquals(0, status, "Deve retornar sucesso para comando validate válido");
        } finally {
            if (signFile.exists()) signFile.delete();
        }
    }

    @Test
    void testRunInvalidFlag() throws Exception {
        int status = App.run(new String[]{"sign", "-f", "somefile.txt"});
        assertEquals(1, status, "Deve retornar erro para flag inválida -f");
    }
}
