import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:bizly_app/shared/utils/secure_storage_util.dart'; // Você criará este util

class ApiService {
  // final String _baseUrl = 'http://10.0.2.2:8080/api/v1'; // Para emulador Android
  final String _baseUrl = 'http://localhost:8080/api/v1'; // Para emulador iOS ou web dev (ajuste se necessário)
  // Para dispositivo físico na mesma rede: use o IP da sua máquina na rede local
  // ex: final String _baseUrl = 'http://192.168.1.10:8080/api/v1';

  final SecureStorageUtil _secureStorage = SecureStorageUtil();

  Future<Map<String, String>> _getHeaders({bool requiresAuth = false}) async {
    Map<String, String> headers = {
      'Content-Type': 'application/json; charset=UTF-8',
    };
    if (requiresAuth) {
      final token = await _secureStorage.getToken();
      if (token != null) {
        headers['Authorization'] = 'Bearer $token';
      }
    }
    return headers;
  }

  // Exemplo de POST (para login, registro)
  Future<http.Response> post(String endpoint, Map<String, dynamic> body, {bool requiresAuth = false}) async {
    final url = Uri.parse('$_baseUrl/$endpoint');
    final headers = await _getHeaders(requiresAuth: requiresAuth);
    return http.post(
      url,
      headers: headers,
      body: jsonEncode(body),
    );
  }

  // Exemplo de GET (para buscar dados)
  Future<http.Response> get(String endpoint, {bool requiresAuth = true}) async {
    final url = Uri.parse('$_baseUrl/$endpoint');
    final headers = await _getHeaders(requiresAuth: requiresAuth);
    return http.get(url, headers: headers);
  }

  // Adicione métodos PUT, DELETE conforme necessário
}