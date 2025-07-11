import 'dart:convert';
import 'package:bizly_app/shared/models/user_model.dart';
import 'package:bizly_app/shared/services/api_service.dart';
import 'package:bizly_app/shared/utils/secure_storage_util.dart';
import 'package:flutter/foundation.dart';

class AuthService with ChangeNotifier {
  final ApiService _apiService = ApiService();
  final SecureStorageUtil _secureStorage = SecureStorageUtil();

  User? _currentUser;
  String? _token;
  bool _isLoading = false;
  String? _errorMessage; // Adicionado

  User? get currentUser => _currentUser;
  String? get token => _token;
  bool get isAuthenticated => _token != null && _currentUser != null;
  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage; // Adicionado

  AuthService() {
    _tryAutoLogin();
  }

  void _setLoading(bool value) {
    _isLoading = value;
    notifyListeners();
  }

  Future<void> _tryAutoLogin() async {
    _setLoading(true);
    final storedToken = await _secureStorage.getToken();
    if (storedToken != null) {
      final storedUserId = await _secureStorage.getUserId();
      if (storedUserId != null) {
        // TODO: Idealmente, faria uma chamada para /users/me para obter dados frescos do usuário
        // e confirmar a validade do token.
        // Por agora, vamos simular com os dados que temos.
        // Você precisaria buscar nome/email reais de algum lugar se quisesse popular completamente.
        // Esta é uma simplificação para o auto-login.
        _currentUser = User(id: storedUserId, name: "Usuário Carregado", email: "auto@login.com"); // Exemplo
        _token = storedToken;
      } else {
        // Se temos token mas não userId, algo está inconsistente. Logout.
        await logout();
      }
    }
    _setLoading(false);
  }

  Future<bool> login(String email, String password) async {
    _setLoading(true);
    try {
      final response = await _apiService.post('auth/login', {
        'email': email,
        'password': password,
      });

      if (response.statusCode == 200) {
        final responseData = jsonDecode(response.body);
        _token = responseData['token'] as String;
        _currentUser = User.fromJson(responseData['user'] as Map<String, dynamic>);

        await _secureStorage.saveToken(_token!);
        await _secureStorage.saveUserId(_currentUser!.id);

        _setLoading(false);
        _errorMessage = null; // Limpa qualquer erro anterior
        return true;
      } else {
        final errorData = jsonDecode(response.body);
        _errorMessage = errorData['error'] ?? 'Falha no login'; // Captura a mensagem de erro
        _setLoading(false);
        return false;
      }
    } catch (e) {
      _errorMessage = 'Erro de conexão: $e'; // Captura erro de exceção
      _setLoading(false);
      return false;
    }
  }

  Future<bool> register(String name, String email, String password) async {
    _setLoading(true);
    try {
      final response = await _apiService.post('users', {
        'name': name,
        'email': email,
        'password': password,
      });

      if (response.statusCode == 201) {
        _setLoading(false);
        _errorMessage = null; // Limpa qualquer erro anterior
        return true;
      } else {
        final errorData = jsonDecode(response.body);
        _errorMessage = errorData['error'] ?? 'Falha no cadastro'; // Captura a mensagem de erro
        _setLoading(false);
        return false;
      }
    } catch (e) {
      _errorMessage = 'Erro de conexão: $e'; // Captura erro de exceção
      _setLoading(false);
      return false;
    }
  }

  Future<void> logout() async {
    _setLoading(true);
    _token = null;
    _currentUser = null;
    await _secureStorage.deleteToken();
    await _secureStorage.deleteUserId(); // <<< CORRIGIDO AQUI
    _setLoading(false);
    // notifyListeners(); // _setLoading já notifica
  }
}