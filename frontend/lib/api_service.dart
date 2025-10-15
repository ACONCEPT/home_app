import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  // The Flutter web app runs in the browser, so it needs to use localhost
  // to access the backend API which is exposed on port 8080
  static const String baseUrl = 'http://localhost:8080/api';

  Future<Map<String, dynamic>> login(String username, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
        }),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        final body = jsonDecode(response.body);
        return {
          'success': false,
          'message': body['message'] ?? 'Login failed',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'message': 'Error connecting to server: $e',
      };
    }
  }

  Future<Map<String, dynamic>> signup(String username, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/signup'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
        }),
      );

      if (response.statusCode == 201) {
        return jsonDecode(response.body);
      } else {
        final body = jsonDecode(response.body);
        return {
          'success': false,
          'message': body['message'] ?? 'Signup failed',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'message': 'Error connecting to server: $e',
      };
    }
  }

  Future<Map<String, dynamic>> getUserTools(String username) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl/user/tools?username=$username'),
        headers: {'Content-Type': 'application/json'},
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        final body = jsonDecode(response.body);
        return {
          'success': false,
          'message': body['message'] ?? 'Failed to get tools',
          'tools': [],
        };
      }
    } catch (e) {
      return {
        'success': false,
        'message': 'Error connecting to server: $e',
        'tools': [],
      };
    }
  }

  Future<Map<String, dynamic>> saveUserTools(String username, List<String> tools) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/user/tools'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'tools': tools,
        }),
      );

      if (response.statusCode == 200) {
        return jsonDecode(response.body);
      } else {
        final body = jsonDecode(response.body);
        return {
          'success': false,
          'message': body['message'] ?? 'Failed to save tools',
        };
      }
    } catch (e) {
      return {
        'success': false,
        'message': 'Error connecting to server: $e',
      };
    }
  }
}
