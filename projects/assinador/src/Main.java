package src;

public class Main {
    public static void main(String[] args) {
        System.out.println("\nAssinatura gerada com sucesso, argumentos:\n");

        for(int i = 0; i < args.length; i++) {
            System.out.println(i + " > " + args[i]);
        }
    }
}