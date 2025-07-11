import 'package:bizly_app/shared/services/appointment_service.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:bizly_app/features/appointments/screens/create_appointment_screen.dart';

class CreateAppointmentScreen extends StatefulWidget {
  const CreateAppointmentScreen({super.key});

  @override
  State<CreateAppointmentScreen> createState() => _CreateAppointmentScreenState();
}

class _CreateAppointmentScreenState extends State<CreateAppointmentScreen> {
  final _formKey = GlobalKey<FormState>();
  final _clientNameController = TextEditingController();
  final _serviceDescController = TextEditingController();
  DateTime? _startTime;
  TimeOfDay? _selectedTime;

  @override
  void dispose() {
    _clientNameController.dispose();
    _serviceDescController.dispose();
    super.dispose();
  }

  Future<void> _selectDate() async {
    final pickedDate = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365)),
    );
    if (pickedDate != null) {
      setState(() {
        _startTime = pickedDate;
      });
    }
  }

  Future<void> _selectTime() async {
    final pickedTime = await showTimePicker(
      context: context,
      initialTime: TimeOfDay.now(),
      builder: (BuildContext context, Widget? child) {
        return MediaQuery(
          data: MediaQuery.of(context).copyWith(alwaysUse24HourFormat: true),
          child: child!,
        );
      },
    );
    if (pickedTime != null) {
      setState(() {
        _selectedTime = pickedTime;
      });
    }
  }

  Future<void> _submitForm() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }
    if (_startTime == null || _selectedTime == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Por favor, selecione data e hora.')),
      );
      return;
    }

    // Combina data e hora
    final finalStartTime = DateTime(
      _startTime!.year,
      _startTime!.month,
      _startTime!.day,
      _selectedTime!.hour,
      _selectedTime!.minute,
    ).toUtc(); // Converte para UTC
    // Exemplo: duração de 1 hora
    final finalEndTime = finalStartTime.add(const Duration(hours: 1));

    final appointmentData = {
      'clientName': _clientNameController.text,
      'serviceDescription': _serviceDescController.text,
      'startTime': finalStartTime.toIso8601String(), // Agora vai gerar com 'Z'
      'endTime': finalEndTime.toIso8601String(),   // Agora vai gerar com 'Z'
      // Adicione outros campos como clientEmail, phone, price, etc.
      'clientEmail': 'placeholder@email.com',
      'clientPhone': '000000000',
      'price': 0.0,
      'notes': '',
    };

    final success = await Provider.of<AppointmentService>(context, listen: false)
        .createAppointment(appointmentData);
    
    if (mounted) {
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Agendamento criado com sucesso!')),
        );
        Navigator.of(context).pop(); // Volta para a tela anterior
      } else {
        final errorMessage = Provider.of<AppointmentService>(context, listen: false).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro: ${errorMessage ?? "Falha ao criar agendamento"}')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    // Usaremos Consumer para mostrar o estado de loading
    return Consumer<AppointmentService>(
      builder: (context, appointmentService, child) {
        return Scaffold(
          appBar: AppBar(
            title: const Text('Novo Agendamento'),
          ),
          body: appointmentService.isLoading
              ? const Center(child: CircularProgressIndicator())
              : SingleChildScrollView(
                  padding: const EdgeInsets.all(16.0),
                  child: Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        TextFormField(
                          controller: _clientNameController,
                          decoration: const InputDecoration(labelText: 'Nome do Cliente'),
                          validator: (value) =>
                              value!.isEmpty ? 'Campo obrigatório' : null,
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _serviceDescController,
                          decoration: const InputDecoration(labelText: 'Descrição do Serviço'),
                           validator: (value) =>
                              value!.isEmpty ? 'Campo obrigatório' : null,
                        ),
                        const SizedBox(height: 16),
                        Row(
                          children: [
                            Expanded(
                              child: Text(
                                _startTime == null
                                    ? 'Nenhuma data selecionada'
                                    : 'Data: ${DateFormat('dd/MM/yyyy').format(_startTime!)}',
                              ),
                            ),
                            TextButton(
                              onPressed: _selectDate,
                              child: const Text('Selecionar Data'),
                            ),
                          ],
                        ),
                        Row(
                          children: [
                            Expanded(
                              child: Text(
                                _selectedTime == null
                                    ? 'Nenhuma hora selecionada'
                                    : 'Hora: ${_selectedTime!.format(context)}',
                              ),
                            ),
                            TextButton(
                              onPressed: _selectTime,
                              child: const Text('Selecionar Hora'),
                            ),
                          ],
                        ),
                        const SizedBox(height: 24),
                        ElevatedButton(
                          onPressed: _submitForm,
                          child: const Text('Salvar Agendamento'),
                        ),
                      ],
                    ),
                  ),
                ),
        );
      },
    );
  }
}