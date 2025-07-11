import 'package:bizly_app/shared/models/appointment_model.dart'; // Adicionado
import 'package:bizly_app/shared/services/appointment_service.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';

class ManageAppointmentScreen extends StatefulWidget {
  final Appointment? appointment; // Agendamento opcional para edição

  const ManageAppointmentScreen({super.key, this.appointment});

  @override
  State<ManageAppointmentScreen> createState() => _ManageAppointmentScreenState();
}

class _ManageAppointmentScreenState extends State<ManageAppointmentScreen> {
  final _formKey = GlobalKey<FormState>();
  final _clientNameController = TextEditingController();
  final _serviceDescController = TextEditingController();
  final _clientEmailController = TextEditingController(); // Adicionado
  final _clientPhoneController = TextEditingController(); // Adicionado
  final _notesController = TextEditingController(); // Adicionado
  final _priceController = TextEditingController(); // Adicionado

  DateTime? _startTime;
  TimeOfDay? _selectedTime;

  bool get isEditing => widget.appointment != null;

  @override
  void initState() {
    super.initState();
    if (isEditing) {
      final appointment = widget.appointment!;
      _clientNameController.text = appointment.clientName;
      _serviceDescController.text = appointment.serviceDescription;
      _clientEmailController.text = appointment.clientEmail;
      _clientPhoneController.text = appointment.clientPhone;
      _notesController.text = appointment.notes;
      _priceController.text = appointment.price.toString();
      _startTime = appointment.startTime.toLocal();
      _selectedTime = TimeOfDay.fromDateTime(_startTime!); // Converte DateTime para TimeOfDay
    }
  }

  @override
  void dispose() {
    _clientNameController.dispose();
    _serviceDescController.dispose();
    _clientEmailController.dispose(); // Adicionado
    _clientPhoneController.dispose(); // Adicionado
    _notesController.dispose(); // Adicionado
    _priceController.dispose(); // Adicionado
    super.dispose();
  }

  Future<void> _selectDate() async {
    final pickedDate = await showDatePicker(
      context: context,
      initialDate: _startTime ?? DateTime.now(), // Usa data existente ou atual
      firstDate: DateTime.now().subtract(const Duration(days: 365 * 5)), // 5 anos atrás
      lastDate: DateTime.now().add(const Duration(days: 365 * 5)), // 5 anos para frente
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
      initialTime: _selectedTime ?? TimeOfDay.now(), // Usa hora existente ou atual
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
      'startTime': DateFormat("yyyy-MM-ddTHH:mm:ss'Z'").format(finalStartTime),
      'endTime': DateFormat("yyyy-MM-ddTHH:mm:ss'Z'").format(finalEndTime),
      'clientEmail': _clientEmailController.text,
      'clientPhone': _clientPhoneController.text,
      'price': double.tryParse(_priceController.text) ?? 0.0,
      'notes': _notesController.text,
    };

    bool success;
    if (isEditing) {
      success = await Provider.of<AppointmentService>(context, listen: false)
          .updateAppointment(widget.appointment!.id, appointmentData);
    } else {
      success = await Provider.of<AppointmentService>(context, listen: false)
          .createAppointment(appointmentData);
    }
    
    if (mounted) {
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(isEditing ? 'Agendamento atualizado com sucesso!' : 'Agendamento criado com sucesso!')),
        );
        Navigator.of(context).pop(); // Volta para a tela anterior
      } else {
        final errorMessage = Provider.of<AppointmentService>(context, listen: false).errorMessage;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro: ${errorMessage ?? (isEditing ? "Falha ao atualizar agendamento" : "Falha ao criar agendamento")}')),
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
            title: Text(isEditing ? 'Editar Agendamento' : 'Novo Agendamento'),
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
                        TextFormField(
                          controller: _clientEmailController,
                          decoration: const InputDecoration(labelText: 'Email do Cliente'),
                          keyboardType: TextInputType.emailAddress,
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _clientPhoneController,
                          decoration: const InputDecoration(labelText: 'Telefone do Cliente'),
                          keyboardType: TextInputType.phone,
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _priceController,
                          decoration: const InputDecoration(labelText: 'Preço'),
                          keyboardType: TextInputType.number,
                          validator: (value) {
                            if (value == null || value.isEmpty) return null;
                            if (double.tryParse(value) == null) {
                              return 'Preço inválido';
                            }
                            return null;
                          },
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _notesController,
                          decoration: const InputDecoration(labelText: 'Notas'),
                          maxLines: 3,
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
                          child: Text(isEditing ? 'Atualizar Agendamento' : 'Salvar Agendamento'),
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