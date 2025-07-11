import 'package:bizly_app/shared/services/auth_service.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:bizly_app/features/appointments/screens/appointment_list_screen.dart';
import 'package:bizly_app/features/clients/screens/client_registration_screen.dart'; // Adicionado

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final authService = Provider.of<AuthService>(context);
    final user = authService.currentUser;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Bizly Home'),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () {
              authService.logout();
              // A navegação para LoginScreen será feita pelo Consumer no AppWidget
            },
          ),
        ],
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('Bem-vindo(a), ${user?.name ?? 'Usuário'}!'),
            Text('Email: ${user?.email ?? ''}'),
            Text('ID: ${user?.id ?? ''}'),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () {
                // Navegar para a tela de listar agendamentos
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (_) => const AppointmentListScreen(),
                  ),
                );
              },
              child: const Text('Ver Meus Agendamentos'),
            ),
            const SizedBox(height: 10),
            ElevatedButton(
              onPressed: () {
                // Navegar para a tela de cadastro de cliente
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (_) => const ClientRegistrationScreen(),
                  ),
                );
              },
              child: const Text('Cadastrar Novo Cliente'),
            ),
          ],
        ),
      ),
    );
  }
}
