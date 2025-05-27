import 'package:bizly_app/app_widget.dart';
import 'package:bizly_app/shared/services/auth_service.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

void main() {
  runApp(
    MultiProvider( // Se você tiver múltiplos providers
      providers: [
        ChangeNotifierProvider(create: (_) => AuthService()),
        // Adicione outros providers aqui (ex: AppointmentService)
      ],
      child: const AppWidget(),
    ),
  );
}