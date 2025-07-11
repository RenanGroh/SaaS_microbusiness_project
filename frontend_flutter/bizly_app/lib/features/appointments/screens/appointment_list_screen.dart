// lib/features/appointments/screens/appointment_list_screen.dart
import 'package:bizly_app/shared/models/appointment_model.dart';
import 'package:bizly_app/shared/services/appointment_service.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';

// --- LINHA DE IMPORT ADICIONADA AQUI ---
import 'package:bizly_app/features/appointments/screens/create_appointment_screen.dart';
// ----------------------------------------

class AppointmentListScreen extends StatefulWidget {
  const AppointmentListScreen({super.key});

  @override
  State<AppointmentListScreen> createState() => _AppointmentListScreenState();
}

class _AppointmentListScreenState extends State<AppointmentListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<AppointmentService>(context, listen: false).fetchAppointments();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Meus Agendamentos'),
      ),
      body: Consumer<AppointmentService>(
        builder: (context, appointmentService, child) {
          if (appointmentService.isLoading && appointmentService.appointments.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (appointmentService.errorMessage != null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('Erro: ${appointmentService.errorMessage}'),
                  const SizedBox(height: 10),
                  ElevatedButton(
                    onPressed: () => appointmentService.fetchAppointments(),
                    child: const Text('Tentar Novamente'),
                  )
                ],
              ),
            );
          }

          if (appointmentService.appointments.isEmpty) {
            return const Center(
              child: Text('Nenhum agendamento encontrado.\nClique no botão + para adicionar um.'),
            );
          }

          return ListView.builder(
            itemCount: appointmentService.appointments.length,
            itemBuilder: (context, index) {
              final appointment = appointmentService.appointments[index];
              return Card(
                margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: ListTile(
                  title: Text(appointment.serviceDescription, style: const TextStyle(fontWeight: FontWeight.bold)),
                  subtitle: Text(
                    'Cliente: ${appointment.clientName}\n'
                    'Data: ${DateFormat('dd/MM/yyyy HH:mm').format(appointment.startTime.toLocal())}',
                  ),
                  trailing: Text(
                    appointment.status,
                    style: TextStyle(
                      color: _getStatusColor(appointment.status),
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  onTap: () {
                    // TODO: Navegar para a tela de detalhes do agendamento
                  },
                ),
              );
            },
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.of(context).push(
            // Agora que o import foi adicionado, esta linha não dará mais erro.
            MaterialPageRoute(builder: (_) => const CreateAppointmentScreen()),
          );
        },
        tooltip: 'Novo Agendamento',
        child: const Icon(Icons.add),
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'CONFIRMED':
        return Colors.green.shade700;
      case 'PENDING':
        return Colors.orange.shade700;
      case 'CANCELLED':
        return Colors.red.shade700;
      case 'COMPLETED':
        return Colors.blue.shade700;
      default:
        return Colors.grey.shade600;
    }
  }
}