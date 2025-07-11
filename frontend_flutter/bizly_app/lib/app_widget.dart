import 'package:bizly_app/features/auth/screens/login_screen.dart';
import 'package:bizly_app/features/home/screens/home_screen.dart'; // Você criará esta
import 'package:bizly_app/shared/services/auth_service.dart';
import 'package:bizly_app/shared/services/appointment_service.dart'; // Adicionado
import 'package:bizly_app/shared/services/client_service.dart'; // Adicionado
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter_localizations/flutter_localizations.dart'; // Adicionado
import 'package:intl/intl.dart'; // Adicionado

class AppWidget extends StatelessWidget {
  const AppWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthService()),
        ChangeNotifierProvider(create: (_) => AppointmentService()),
        ChangeNotifierProvider(create: (_) => ClientService()), // Adicionado
      ],
      child: MaterialApp(
        title: 'Bizly App',
        theme: ThemeData(
          primarySwatch: Colors.blue, // Personalize seu tema
          // visualDensity: VisualDensity.adaptivePlatformDensity,
        ),
        localizationsDelegates: const [
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        supportedLocales: const [
          Locale('pt', 'BR'), // Português do Brasil
        ],
        locale: const Locale('pt', 'BR'), // Define o locale padrão
        home: Consumer<AuthService>(
          builder: (context, authService, _) {
            if (authService.isLoading) {
              return const Scaffold(body: Center(child: CircularProgressIndicator()));
            }
            if (authService.isAuthenticated) {
              return const HomeScreen(); // Tela principal após login
            } else {
              return const LoginScreen(); // Tela de login
            }
          },
        ),
        // Defina suas rotas aqui se for usar navegação nomeada
        // routes: {
        //   '/login': (context) => LoginScreen(),
        //   '/register': (context) => RegisterScreen(),
        //   '/home': (context) => HomeScreen(),
        // },
      ),
    );
  }
}