import 'package:bizly_app/features/appointments/screens/manage_appointment_screen.dart'; // Adicionado
import 'package:bizly_app/shared/models/appointment_model.dart';
import 'package:bizly_app/shared/services/appointment_service.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';

class AppointmentDetailScreen extends StatelessWidget {
  final Appointment appointment;

  const AppointmentDetailScreen({super.key, required this.appointment});

  @override
  Widget build(BuildContext context) {
    final appointmentService = Provider.of<AppointmentService>(context, listen: false);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Detalhes do Agendamento'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildDetailRow('Cliente:', appointment.clientName),
            _buildDetailRow('Serviço:', appointment.serviceDescription),
            _buildDetailRow('Data/Hora:', DateFormat('dd/MM/yyyy HH:mm').format(appointment.startTime.toLocal())),
            _buildDetailRow('Status:', appointment.status),
            _buildDetailRow('Preço:', 'R\$ ${appointment.price.toStringAsFixed(2)}'),
            _buildDetailRow('Notas:', appointment.notes.isEmpty ? 'Nenhuma' : appointment.notes),
            const SizedBox(height: 20),
            // Botões de Ação
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (_) => ManageAppointmentScreen(appointment: appointment),
                  ),
                ).then((_) => appointmentService.fetchAppointments()); // Recarrega a lista ao voltar
              },
              child: const Text('Editar Agendamento'),
            ),
            const SizedBox(height: 10),
            ElevatedButton(
              onPressed: () => _confirmCancel(context, appointmentService),
              style: ElevatedButton.styleFrom(backgroundColor: Colors.orange),
              child: const Text('Cancelar Agendamento'),
            ),
            const SizedBox(height: 10),
            ElevatedButton(
              onPressed: () => _confirmComplete(context, appointmentService),
              style: ElevatedButton.styleFrom(backgroundColor: Colors.green),
              child: const Text('Marcar como Concluído'),
            ),
            const SizedBox(height: 10),
            ElevatedButton(
              onPressed: () => _confirmDelete(context, appointmentService),
              style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
              child: const Text('Excluir Agendamento'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontSize: 16),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _confirmCancel(BuildContext context, AppointmentService service) async {
    final bool? confirm = await showDialog<bool>(
      context: context,
      builder: (BuildContext dialogContext) {
        return AlertDialog(
          title: const Text('Confirmar Cancelamento'),
          content: const Text('Tem certeza que deseja cancelar este agendamento?'),
          actions: <Widget>[
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(false),
              child: const Text('Não'),
            ),
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(true),
              child: const Text('Sim'),
            ),
          ],
        );
      },
    );

    if (confirm == true) {
      final success = await service.cancelAppointment(appointment.id);
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Agendamento cancelado com sucesso!')),
        );
        Navigator.of(context).pop(); // Volta para a lista
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro ao cancelar: ${service.errorMessage}')),
        );
      }
    }
  }

  Future<void> _confirmComplete(BuildContext context, AppointmentService service) async {
    final bool? confirm = await showDialog<bool>(
      context: context,
      builder: (BuildContext dialogContext) {
        return AlertDialog(
          title: const Text('Confirmar Conclusão'),
          content: const Text('Tem certeza que deseja marcar este agendamento como concluído?'),
          actions: <Widget>[
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(false),
              child: const Text('Não'),
            ),
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(true),
              child: const Text('Sim'),
            ),
          ],
        );
      },
    );

    if (confirm == true) {
      final success = await service.updateAppointmentStatus(appointment.id, 'COMPLETED');
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Agendamento marcado como concluído!')),
        );
        Navigator.of(context).pop(); // Volta para a lista
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro ao concluir: ${service.errorMessage}')),
        );
      }
    }
  }

  Future<void> _confirmDelete(BuildContext context, AppointmentService service) async {
    final bool? confirm = await showDialog<bool>(
      context: context,
      builder: (BuildContext dialogContext) {
        return AlertDialog(
          title: const Text('Confirmar Exclusão'),
          content: const Text('Tem certeza que deseja excluir este agendamento permanentemente?'),
          actions: <Widget>[
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(false),
              child: const Text('Não'),
            ),
            TextButton(
              onPressed: () => Navigator.of(dialogContext).pop(true),
              child: const Text('Sim'),
            ),
          ],
        );
      },
    );

    if (confirm == true) {
      final success = await service.deleteAppointment(appointment.id);
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Agendamento excluído com sucesso!')),
        );
        Navigator.of(context).pop(); // Volta para a lista
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Erro ao excluir: ${service.errorMessage}')),
        );
      }
    }
  }
}
