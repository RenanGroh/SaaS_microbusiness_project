import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class SecureStorageUtil {
  final _storage = const FlutterSecureStorage();
  static const _tokenKey = 'jwt_token';
  static const _userIdKey = 'user_id'; // Opcional

  Future<void> saveToken(String token) async {
    await _storage.write(key: _tokenKey, value: token);
  }

  Future<String?> getToken() async {
    return await _storage.read(key: _tokenKey);
  }

  Future<void> deleteToken() async {
    await _storage.delete(key: _tokenKey);
  }

  // Opcional: salvar outros dados do usuário
  Future<void> saveUserId(String userId) async {
     await _storage.write(key: _userIdKey, value: userId);
  }

  Future<String?> getUserId() async {
     return await _storage.read(key: _userIdKey);
  }

  Future<void> deleteUserId() async {
    await _storage.delete(key: _userIdKey);
  }
}