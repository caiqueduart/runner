public class App {
    public static void main(String[] args) {
        System.out.println("Início do projeto runner!");

        if (args.length == 0) {
            System.out.println("Uso: java App <arquivo_para_validar>");
        } else {
            System.out.println("Arquivo a validar: " + args[0]);
        }
    }
}