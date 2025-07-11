import 'package:flutter/foundation.dart';

class Appointment {
  final String id; // UUID
  final String userId; // UUID
  final String? clientId; // UUID Opcional
  final String clientName;
  final String clientEmail;
  final String clientPhone;
  final String serviceDescription;
  final DateTime startTime;
  final DateTime endTime;
  final String status;
  final String notes;
  final double price;
  final DateTime createdAt;
  final DateTime updatedAt;

  Appointment({
    required this.id,
    required this.userId,
    this.clientId,
    required this.clientName,
    required this.clientEmail,
    required this.clientPhone,
    required this.serviceDescription,
    required this.startTime,
    required this.endTime,
    required this.status,
    required this.notes,
    required this.price,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Appointment.fromJson(Map<String, dynamic> json) {
    return Appointment(
      id: json['id'] as String,
      userId: json['userId'] as String,
      clientId: json['clientId'] as String?, // Pode ser nulo
      clientName: json['clientName'] as String,
      clientEmail: json['clientEmail'] as String,
      clientPhone: json['clientPhone'] as String,
      serviceDescription: json['serviceDescription'] as String,
      startTime: DateTime.parse(json['startTime'] as String), // Parse da string ISO 8601
      endTime: DateTime.parse(json['endTime'] as String),
      status: json['status'] as String,
      notes: json['notes'] as String,
      price: (json['price'] as num).toDouble(), // Converte num para double
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }
}