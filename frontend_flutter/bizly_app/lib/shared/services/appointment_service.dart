import 'dart:convert';
import 'package:bizly_app/shared/models/appointment_model.dart';
import 'package:bizly_app/shared/services/api_service.dart';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http; // Adicionado

class AppointmentService with ChangeNotifier {
  final ApiService _apiService = ApiService();

  List<Appointment> _appointments = [];
  bool _isLoading = false;
  String? _errorMessage;

  List<Appointment> get appointments => _appointments;
  bool get isLoading => _isLoading;
  String? get errorMessage => _errorMessage;

  void _setState({bool loading = false, String? error}) {
    _isLoading = loading;
    _errorMessage = error;
    notifyListeners();
  }

  Future<void> fetchAppointments() async {
    _setState(loading: true);
    try {
      final response = await _apiService.getAppointments();
      if (response.statusCode == 200) {
        final List<dynamic> responseData = jsonDecode(response.body);
        _appointments = responseData.map((data) => Appointment.fromJson(data)).toList();
        _setState(loading: false);
      } else {
        final errorData = jsonDecode(response.body);
        _setState(loading: false, error: errorData['error'] ?? 'Falha ao buscar agendamentos');
      }
    } catch (e) {
      _setState(loading: false, error: 'Erro de conexão: $e');
    }
  }

  Future<bool> _performAppointmentAction(
      Future<http.Response> Function() apiCall, int successStatusCode, String actionName) async {
    _setState(loading: true);
    try {
      final response = await apiCall();
      if (response.statusCode == successStatusCode) {
        await fetchAppointments();
        _setState(loading: false);
        return true;
      } else {
        final errorData = jsonDecode(response.body);
        _setState(loading: false, error: errorData['error'] ?? 'Falha ao $actionName agendamento');
        return false;
      }
    } catch (e) {
      _setState(loading: false, error: 'Erro de conexão: $e');
      return false;
    }
  }

  Future<bool> createAppointment(Map<String, dynamic> data) async {
    _setState(loading: true);
    try {
      final response = await _apiService.createAppointment(data);
       if (response.statusCode == 201) {
        // Sucesso! Recarrega a lista para incluir o novo agendamento.
        await fetchAppointments(); // Já chama notifyListeners() internamente
        _setState(loading: false); // Apenas para garantir que o loading termine
        return true;
      } else {
        final errorData = jsonDecode(response.body);
        _setState(loading: false, error: errorData['error'] ?? 'Falha ao criar agendamento');
        return false;
      }
    } catch (e) {
      _setState(loading: false, error: 'Erro de conexão: $e');
      return false;
    }
  }

  Future<bool> cancelAppointment(String appointmentId) async {
    return _performAppointmentAction(
        () => _apiService.cancelAppointment(appointmentId), 200, 'cancelar');
  }

  Future<bool> updateAppointmentStatus(String appointmentId, String newStatus) async {
    return _performAppointmentAction(
        () => _apiService.updateAppointment(appointmentId, {'status': newStatus}), 200, 'atualizar status');
  }

  Future<bool> updateAppointment(String appointmentId, Map<String, dynamic> data) async {
    return _performAppointmentAction(
        () => _apiService.updateAppointment(appointmentId, data), 200, 'atualizar');
  }

  Future<bool> deleteAppointment(String appointmentId) async {
    return _performAppointmentAction(
        () => _apiService.deleteAppointment(appointmentId), 204, 'excluir');
  }
}