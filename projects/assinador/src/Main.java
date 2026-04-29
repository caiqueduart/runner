package src;

public class Main {
    public static void main(String[] args) {

        String file = args[0];

        String assinaturaSimulada = "SIMULATED_SIG_" + file.hashCode();
        System.out.println("Resultado: " + assinaturaSimulada);
    }
}