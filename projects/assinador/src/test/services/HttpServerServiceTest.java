package test.services;

import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import services.HttpServerService;

import java.io.File;
import java.io.FileWriter;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.file.Path;

import static org.junit.jupiter.api.Assertions.*;

class HttpServerServiceTest {

    private static final int TEST_PORT = 8888;
    private static HttpClient client;

    @TempDir
    static Path tempDir;

    private static File testFile;

    @BeforeAll
    static void startServer() throws Exception {
        testFile = tempDir.resolve("server-test.txt").toFile();
        try (FileWriter writer = new FileWriter(testFile)) {
            writer.write("conteudo para teste de servidor");
        }

        client = HttpClient.newHttpClient();
        
        // Start server in a separate thread
        new Thread(() -> HttpServerService.start(TEST_PORT, 1)).start();
        
        // Wait a bit for server to start
        Thread.sleep(1000);
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
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/sign"))
                .POST(HttpRequest.BodyPublishers.ofString(testFile.getAbsolutePath()))
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        
        assertEquals(200, response.statusCode());
        assertTrue(response.body().contains("gerou o código de assinatura"));
        
        // Cleanup generated signature file
        String nameOnly = testFile.getName();
        String baseName = nameOnly.substring(0, nameOnly.lastIndexOf('.'));
        String extension = nameOnly.substring(nameOnly.lastIndexOf('.') + 1);
        String outputName = baseName + "-" + extension + "-assinatura.txt";
        new File(System.getProperty("user.dir"), outputName).delete();
    }

    @Test
    void testSignEndpointMissingFile() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/sign"))
                .POST(HttpRequest.BodyPublishers.ofString(""))
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        
        assertEquals(400, response.statusCode());
        assertTrue(response.body().contains("O parâmetro '--file' é obrigatório"));
    }

    @Test
    void testInvalidMethod() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("http://localhost:" + TEST_PORT + "/sign"))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        
        assertEquals(405, response.statusCode());
    }
}
