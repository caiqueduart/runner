package services;

public class ServiceFactory {
    private static SignatureService signatureService;

    public static synchronized SignatureService getSignatureService() {
        if (signatureService == null) {
            signatureService = new SignatureService();
        }
        return signatureService;
    }
}
