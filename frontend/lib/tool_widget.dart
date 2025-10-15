import 'package:flutter/material.dart';

class ToolWidget extends StatelessWidget {
  final String name;
  final IconData icon;
  final Color color;

  const ToolWidget({
    super.key,
    required this.name,
    required this.icon,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 3,
      child: InkWell(
        onTap: () {
          // TODO: Navigate to tool page when implemented
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('$name tool coming soon!'),
              duration: const Duration(seconds: 2),
            ),
          );
        },
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                icon,
                size: 48,
                color: color,
              ),
              const SizedBox(height: 12),
              Text(
                name,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// Available tools configuration
class ToolConfig {
  static const List<Map<String, dynamic>> availableTools = [
    {
      'id': 'todo',
      'name': 'To-Do',
      'icon': Icons.check_box,
      'color': Colors.blue,
    },
    {
      'id': 'chore-admin',
      'name': 'Chore Admin',
      'icon': Icons.cleaning_services,
      'color': Colors.green,
    },
    {
      'id': 'notes',
      'name': 'Notes',
      'icon': Icons.note,
      'color': Colors.orange,
    },
    {
      'id': 'calendar',
      'name': 'Calendar',
      'icon': Icons.calendar_month,
      'color': Colors.red,
    },
    {
      'id': 'budget',
      'name': 'Budget',
      'icon': Icons.attach_money,
      'color': Colors.purple,
    },
  ];

  static Map<String, dynamic>? getToolById(String id) {
    try {
      return availableTools.firstWhere((tool) => tool['id'] == id);
    } catch (e) {
      return null;
    }
  }
}
