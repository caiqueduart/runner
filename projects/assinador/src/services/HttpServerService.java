package services;

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
    private static final long DEFAULT_TIMEOUT_MINUTES = 5;
    
    private static HttpServer server;
    private static ScheduledExecutorService scheduler;
    private static long startTime;
    private static long timeoutMillis;
    private static final AtomicLong lastRequestTime = new AtomicLong(System.currentTimeMillis());

    public static void start(int port, long timeoutMinutes) {
        int effectivePort = port > 0 ? port : DEFAULT_PORT;
        long effectiveTimeout = timeoutMinutes > 0 ? timeoutMinutes : DEFAULT_TIMEOUT_MINUTES;
        
        startTime = System.currentTimeMillis();
        timeoutMillis = TimeUnit.MINUTES.toMillis(effectiveTimeout);

        try {
            server = HttpServer.create(new InetSocketAddress(effectivePort), 0);
            server.setExecutor(Executors.newVirtualThreadPerTaskExecutor());

            server.createContext("/sign", new SignHandler());
            server.createContext("/validate", new ValidateHandler());
            server.createContext("/health", new HealthHandler());
            server.createContext("/stop", new StopHandler());

            server.start();
            logServerEvent("INFO", "Online na porta " + effectivePort);
            logServerEvent("INFO", "Auto-desligamento em " + effectiveTimeout + "m.");

            startTimeoutChecker(effectiveTimeout);

        } catch (IOException e) {
            logServerEvent("ERROR", "Erro ao iniciar: " + e.getMessage());
        }
    }

    private static void startTimeoutChecker(long timeoutMinutes) {
        scheduler = Executors.newSingleThreadScheduledExecutor();
        scheduler.scheduleAtFixedRate(() -> {
            long inactiveTime = System.currentTimeMillis() - lastRequestTime.get();

            if (inactiveTime >= timeoutMillis) {
                logServerEvent("INFO", "Encerrando por inatividade...");
                stopServer();
            }
        }, 5, 5, TimeUnit.SECONDS);
    }

    private static void updateLastRequestTime() {
        lastRequestTime.set(System.currentTimeMillis());
    }

    public static void stop() {
        stopServer();
    }

    private static void stopServer() {
        if (scheduler != null) {
            scheduler.shutdownNow();
        }
        if (server != null) {
            server.stop(0);
        }
        logServerEvent("INFO", "Encerrado.");
    }

    static class SignHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                exchange.sendResponseHeaders(405, -1);
                return;
            }

            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            
            try {
                // parse simplificado de JSON manual (para evitar dependências)
                String fileName = extractFromJson(body, "file");
                
                if (fileName == null || fileName.isEmpty()) {
                    sendJsonResponse(exchange, """
                        {
                            "error": "O parâmetro 'file' no JSON é obrigatório.",
                            "status": 400,
                            "type": "user"
                        }
                        """, 400);
                    return;
                }

                updateLastRequestTime();
                SignatureService service = ServiceFactory.getSignatureService();
                SignatureResult res = service.sign(fileName);
                
                String json = """
                    {
                        "message": "Arquivo assinado com sucesso.",
                        "fileName": "%s",
                        "code": "%s",
                        "signOutputPath": "%s",
                        "status": 200
                    }
                    """.formatted(
                        res.fileName().replace("\\", "/"), 
                        res.code(), 
                        res.filePath().replace("\\", "/")
                    );
                sendJsonResponse(exchange, json, 200);
            } catch (Exception e) {
                boolean isUserError = e.getMessage().contains("não encontrado") || e.getMessage().contains("obrigatório");
                int statusCode = isUserError ? 400 : 500;
                String type = isUserError ? "user" : "system";
                sendJsonResponse(exchange, """
                    {
                        "error": "%s",
                        "status": %d,
                        "type": "%s"
                    }
                    """.formatted(e.getMessage(), statusCode, type), statusCode);
            }
        }
    }

    static class ValidateHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                exchange.sendResponseHeaders(405, -1);
                return;
            }

            String body = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);

            try {
                String fileName = extractFromJson(body, "file");
                
                if (fileName == null || fileName.isEmpty()) {
                    sendJsonResponse(exchange, """
                        {
                            "error": "O parâmetro 'file' no JSON é obrigatório.",
                            "status": 400,
                            "type": "user"
                        }
                        """, 400);
                    return;
                }

                updateLastRequestTime();
                SignatureService service = ServiceFactory.getSignatureService();
                ValidationResult res = service.validate(fileName);
                
                String json = """
                    {
                        "message": "Validação concluída.",
                        "fileName": "%s",
                        "code": "%s",
                        "valid": %b,
                        "status": 200
                    }
                    """.formatted(
                        res.fileName().replace("\\", "\\\\"), 
                        res.code(), 
                        res.valid()
                    );
                sendJsonResponse(exchange, json, 200);
            } catch (Exception e) {
                boolean isUserError = e.getMessage().contains("não encontrado") || e.getMessage().contains("obrigatório") || e.getMessage().contains(".txt");
                int statusCode = isUserError ? 400 : 500;
                String type = isUserError ? "user" : "system";
                sendJsonResponse(exchange, """
                    {
                        "error": "%s",
                        "status": %d,
                        "type": "%s"
                    }
                    """.formatted(e.getMessage(), statusCode, type), statusCode);
            }
        }
    }

    private static String extractFromJson(String json, String key) {
        String pattern = "\"" + key + "\":\\s*\"([^\"]*)\"";
        java.util.regex.Matcher matcher = java.util.regex.Pattern.compile(pattern).matcher(json);
        if (matcher.find()) {
            return matcher.group(1);
        }
        return null;
    }

    static class HealthHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            long now = System.currentTimeMillis();
            long uptimeSeconds = (now - startTime) / 1000;
            long remainingMillis = timeoutMillis - (now - lastRequestTime.get());
            long remainingSeconds = Math.max(0, remainingMillis / 1000);
            
            String json = """
                {
                    "status": "OK",
                    "uptimeSeconds": %d,
                    "remainingSeconds": %d,
                    "code": 200
                }
                """.formatted(uptimeSeconds, remainingSeconds);
            
            sendJsonResponse(exchange, json, 200);
        }
    }

    static class StopHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
                exchange.sendResponseHeaders(405, -1);
                return;
            }
            sendJsonResponse(exchange, """
                {
                    "message": "Sinal de encerramento recebido.",
                    "status": 200
                }
                """, 200);
            new Thread(() -> {
                try {
                    Thread.sleep(500);
                    logServerEvent("INFO", "Encerrando...");
                    stopServer();
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }).start();
        }
    }

    private static void sendJsonResponse(HttpExchange exchange, String json, int statusCode) throws IOException {
        byte[] bytes = json.trim().getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=UTF-8");
        exchange.sendResponseHeaders(statusCode, bytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }

    private static void logServerEvent(String level, String message) {
        String json = """
            {
                "level": "%s",
                "component": "SERVER",
                "message": "%s",
                "timestamp": %d
            }
            """.formatted(level, message, System.currentTimeMillis());
        System.out.println(json);
    }
}
