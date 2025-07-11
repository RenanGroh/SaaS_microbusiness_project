import 'package:flutter/foundation.dart';

class Client {
  final String id;
  final String userId;
  final String name;
  final String email;
  final String phone;
  final String notes;
  final DateTime createdAt;
  final DateTime updatedAt;

  Client({
    required this.id,
    required this.userId,
    required this.name,
    required this.email,
    required this.phone,
    required this.notes,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Client.fromJson(Map<String, dynamic> json) {
    return Client(
      id: json['id'] as String,
      userId: json['userId'] as String,
      name: json['name'] as String,
      email: json['email'] as String,
      phone: json['phone'] as String,
      notes: json['notes'] as String,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'userId': userId,
      'name': name,
      'email': email,
      'phone': phone,
      'notes': notes,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }
}