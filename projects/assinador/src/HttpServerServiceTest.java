import static org.junit.jupiter.api.Assertions.*;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.io.TempDir;
import services.HttpServerService;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.file.Path;
import java.nio.file.Files;
import java.io.File;

class HttpServerServiceTest {

    private static final int TEST_PORT = 8081;
    private static HttpClient client = HttpClient.newHttpClient();

    @TempDir
    static Path tempDir;

    @BeforeAll
    static void startServer() {
        // Inicia o servidor em uma porta diferente da padrão para evitar conflitos
        new Thread(() -> HttpServerService.start(TEST_PORT, 1)).start();
        // Pequena espera para o servidor subir
        try { Thread.sleep(1000); } catch (InterruptedException e) {}
    }

    @AfterAll
    static void stopServer() {
        HttpServerService.stop();
    }

    @Test
    void testHealthEndpoint() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/health"))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        assertEquals(200, response.statusCode());
        assertTrue(response.body().contains("Status: OK"));
    }

    @Test
    void testSignEndpoint() throws Exception {
        Path testFile = tempDir.resolve("server-test.txt");
        Files.writeString(testFile, "conteudo para o servidor");

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/sign"))
                .POST(HttpRequest.BodyPublishers.ofString(testFile.toAbsolutePath().toString()))
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        assertEquals(200, response.statusCode());
        assertTrue(response.body().contains("gerou o código de assinatura"));

        // Cleanup the generated signature file
        String expectedSignFileName = "server-test-txt-assinatura.txt";
        File signFile = new File(System.getProperty("user.dir"), expectedSignFileName);
        if (signFile.exists()) signFile.delete();
    }

    @Test
    void testSignEndpointEmptyBody() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/sign"))
                .POST(HttpRequest.BodyPublishers.ofString(""))
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        assertEquals(400, response.statusCode());
        assertTrue(response.body().contains("Erro do usuário"));
    }

    @Test
    void testValidateEndpoint() throws Exception {
        Path testFile = tempDir.resolve("validate-test.txt");
        Files.writeString(testFile, "conteudo para validar");

        // Primeiro assina para gerar o arquivo
        HttpRequest signRequest = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/sign"))
                .POST(HttpRequest.BodyPublishers.ofString(testFile.toAbsolutePath().toString()))
                .build();
        client.send(signRequest, HttpResponse.BodyHandlers.ofString());

        String signFileName = "validate-test-txt-assinatura.txt";
        File signFile = new File(System.getProperty("user.dir"), signFileName);

        try {
            HttpRequest validateRequest = HttpRequest.newBuilder()
                    .uri(URI.create("http://localhost:" + TEST_PORT + "/validate"))
                    .POST(HttpRequest.BodyPublishers.ofString(signFile.getAbsolutePath()))
                    .build();

            HttpResponse<String> response = client.send(validateRequest, HttpResponse.BodyHandlers.ofString());
            assertEquals(200, response.statusCode());
            assertTrue(response.body().contains("está assinado sob o código"));
        } finally {
            if (signFile.exists()) signFile.delete();
        }
    }

    @Test
    void testStopEndpoint() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/stop"))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        assertEquals(200, response.statusCode());
        assertTrue(response.body().contains("Sinal de encerramento recebido"));
    }
}
