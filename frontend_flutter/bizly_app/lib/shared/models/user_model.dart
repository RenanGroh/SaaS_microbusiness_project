// Para @required ou late final

class User {
  final String id; // UUID como String
  final String name;
  final String email;
  // final DateTime? createdAt; // Opcional
  // final DateTime? updatedAt; // Opcional

  User({
    required this.id,
    required this.name,
    required this.email,
    // this.createdAt,
    // this.updatedAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      name: json['name'] as String,
      email: json['email'] as String,
      // createdAt: json['createdAt'] != null ? DateTime.parse(json['createdAt'] as String) : null,
      // updatedAt: json['updatedAt'] != null ? DateTime.parse(json['updatedAt'] as String) : null,
    );
  }

  Map<String, dynamic> toJson() { // Para enviar dados (ex: no cadastro, se necessário)
    return {
      'id': id,
      'name': name,
      'email': email,
      // 'createdAt': createdAt?.toIso8601String(),
      // 'updatedAt': updatedAt?.toIso8601String(),
    };
  }
}