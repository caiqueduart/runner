package src.services;

public class Tint {
    public static final String RESET = "\u001B[0m";
    public static final String RED = "\u001B[31m";
    public static final String GREEN = "\u001B[32m";
    public static final String YELLOW = "\u001B[33m";
    public static final String BLUE = "\u001B[34m";
    public static final String CYAN = "\u001B[36m";

    public static void logInfo(String prefix, String message) {
        System.out.println(CYAN + "[" + prefix + "] " + RESET + message);
    }

    public static void logSuccess(String prefix, String message) {
        System.out.println(CYAN + "[" + prefix + "] " + GREEN + message + RESET);
    }

    public static void logWarn(String prefix, String message) {
        System.out.println(CYAN + "[" + prefix + "] " + YELLOW + message + RESET);
    }

    public static void logError(String prefix, String message) {
        System.err.println(CYAN + "[" + prefix + "] " + RED + message + RESET);
    }
}
