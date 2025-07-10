// lib/shared/services/api_service.dart
import 'dart:convert';
import 'dart:io'; // <<< ADICIONE ESTE IMPORT para 'Platform'

import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart'; // Para 'kIsWeb'
import 'package:bizly_app/shared/utils/secure_storage_util.dart';

class ApiService {
  // Usando um getter para definir a URL base dinamicamente
  String get _baseUrl {
    // kIsWeb é uma constante do Flutter que é true se o app estiver rodando na web.
    if (kIsWeb) {
      return 'http://localhost:8080/api/v1';
    }

    // Platform.isAndroid verifica se estamos rodando em um dispositivo/emulador Android.
    if (Platform.isAndroid) {
      // 10.0.2.2 é como o emulador Android se refere ao 'localhost' da máquina host.
      return 'http://10.0.2.2:8080/api/v1';
    }

    // Fallback para iOS (que pode usar localhost) ou outras plataformas.
    // NOTA: Se você for testar em um dispositivo FÍSICO Android conectado na mesma
    // rede Wi-Fi, você teria que usar o endereço IP da sua máquina na rede.
    // Ex: return 'http://192.168.1.5:8080/api/v1';
    return 'http://localhost:8080/api/v1';
  }

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

  Future<http.Response> post(String endpoint, Map<String, dynamic> body, {bool requiresAuth = false}) async {
    final url = Uri.parse('$_baseUrl/$endpoint');
    final headers = await _getHeaders(requiresAuth: requiresAuth);
    return http.post(
      url,
      headers: headers,
      body: jsonEncode(body),
    );
  }

  Future<http.Response> get(String endpoint, {bool requiresAuth = true}) async {
    final url = Uri.parse('$_baseUrl/$endpoint');
    final headers = await _getHeaders(requiresAuth: requiresAuth);
    return http.get(url, headers: headers);
  }

  // Adicione métodos PUT, DELETE conforme necessário
}