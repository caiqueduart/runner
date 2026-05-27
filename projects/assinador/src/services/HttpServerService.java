package src.services;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicLong;

public class HttpServerService {
    private static final int DEFAULT_PORT = 8080;
    private static final long DEFAULT_TIMEOUT_MINUTES = 10;
    
    private static HttpServer server;
    private static ScheduledExecutorService scheduler;
    private static final AtomicLong lastRequestTime = new AtomicLong(System.currentTimeMillis());

    public static void start(int port, long timeoutMinutes) {
        int effectivePort = port > 0 ? port : DEFAULT_PORT;
        long effectiveTimeout = timeoutMinutes > 0 ? timeoutMinutes : DEFAULT_TIMEOUT_MINUTES;

        try {
            server = HttpServer.create(new InetSocketAddress(effectivePort), 0);
            
            // Usando Virtual Threads (Java 21) para o executor do servidor
            server.setExecutor(Executors.newVirtualThreadPerTaskExecutor());

            server.createContext("/sign", new SignHandler());
            server.createContext("/validate", new ValidateHandler());
            server.createContext("/health", new HealthHandler());
            server.createContext("/stop", new StopHandler());

            server.start();
            System.out.println(Tint.GREEN + "Servidor do Assinador iniciado na porta " + effectivePort + Tint.RESET);
            System.out.println(Tint.YELLOW + "Timeout de inatividade configurado para " + effectiveTimeout + " minutos." + Tint.RESET);

            startTimeoutChecker(effectiveTimeout);

        } catch (IOException e) {
            System.err.println(Tint.RED + "Erro ao iniciar o servidor: " + e.getMessage() + Tint.RESET);
        }
    }

    private static void startTimeoutChecker(long timeoutMinutes) {
        scheduler = Executors.newSingleThreadScheduledExecutor();
        scheduler.scheduleAtFixedRate(() -> {
            long inactiveTime = System.currentTimeMillis() - lastRequestTime.get();
            if (inactiveTime > TimeUnit.MINUTES.toMillis(timeoutMinutes)) {
                System.out.println(Tint.YELLOW + "Servidor encerrando por inatividade..." + Tint.RESET);
                stopServer();
            }
        }, 1, 1, TimeUnit.MINUTES);
    }

    private static void updateLastRequestTime() {
        lastRequestTime.set(System.currentTimeMillis());
    }

    private static void stopServer() {
        if (server != null) {
            server.stop(0);
        }
        if (scheduler != null) {
            scheduler.shutdownNow();
        }
        System.out.println(Tint.GREEN + "Servidor encerrado." + Tint.RESET);
        System.exit(0);
    }

    static class SignHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            updateLastRequestTime();
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                exchange.sendResponseHeaders(405, -1);
                return;
            }

            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            // Simples extração de parâmetro (em um cenário real usaríamos JSON)
            String fileName = body.trim();
            
            if (fileName.isEmpty()) {
                sendResponse(exchange, "Erro: Nome do arquivo não fornecido.", 400);
                return;
            }

            String signature = SignatureService.makeSimulatedSignCode(fileName);
            String response = "O Arquivo '" + fileName + "' gerou o código de assinatura '" + signature + "'.";
            sendResponse(exchange, response, 200);
        }
    }

    static class ValidateHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            updateLastRequestTime();
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                exchange.sendResponseHeaders(405, -1);
                return;
            }

            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            String fileName = body.trim();

            if (fileName.isEmpty()) {
                sendResponse(exchange, "Erro: Nome do arquivo não fornecido.", 400);
                return;
            }

            String signature = SignatureService.makeSimulatedSignCode(fileName);
            String response = "O Arquivo '" + fileName + "' está assinado sob o código '" + signature + "'.";
            sendResponse(exchange, response, 200);
        }
    }

    static class HealthHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            updateLastRequestTime();
            sendResponse(exchange, "OK", 200);
        }
    }

    static class StopHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            sendResponse(exchange, "Encerrando servidor...", 200);
            new Thread(() -> {
                try {
                    Thread.sleep(500);
                    stopServer();
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }).start();
        }
    }

    private static void sendResponse(HttpExchange exchange, String response, int statusCode) throws IOException {
        byte[] bytes = response.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "text/plain; charset=UTF-8");
        exchange.sendResponseHeaders(statusCode, bytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }
}
