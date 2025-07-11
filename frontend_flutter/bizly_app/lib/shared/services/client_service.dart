import 'dart:convert';
import 'package:bizly_app/shared/models/client_model.dart';
import 'package:bizly_app/shared/services/api_service.dart';
import 'package:flutter/foundation.dart';

class ClientService with ChangeNotifier {
  final ApiService _apiService = ApiService();

  List<Client> _clients = [];
  bool _isLoading = false;
  String? _errorMessage;

  List<Client> get clients => _clients;
  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage;

  void _setState({bool loading = false, String? error}) {
    _isLoading = loading;
    _errorMessage = error;
    notifyListeners();
  }

  Future<bool> createClient(Map<String, dynamic> data) async {
    _setState(loading: true);
    try {
      final response = await _apiService.createClient(data);
      if (response.statusCode == 201) {
        // Sucesso! Recarrega a lista para incluir o novo cliente.
        await fetchClients(); // Já chama notifyListeners() internamente
        _setState(loading: false); // Apenas para garantir que o loading termine
        return true;
      } else {
        final errorData = jsonDecode(response.body);
        _setState(loading: false, error: errorData['error'] ?? 'Falha ao criar cliente');
        return false;
      }
    } catch (e) {
      _setState(loading: false, error: 'Erro de conexão: $e');
      return false;
    }
  }

  Future<void> fetchClients() async {
    _setState(loading: true);
    try {
      final response = await _apiService.getClients();
      if (response.statusCode == 200) {
        final List<dynamic> responseData = jsonDecode(response.body);
        _clients = responseData.map((data) => Client.fromJson(data)).toList();
        _setState(loading: false);
      } else {
        final errorData = jsonDecode(response.body);
        _setState(loading: false, error: errorData['error'] ?? 'Falha ao buscar clientes');
      }
    } catch (e) {
      _setState(loading: false, error: 'Erro de conexão: $e');
    }
  }
}