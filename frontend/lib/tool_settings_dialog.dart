import 'package:flutter/material.dart';
import 'tool_widget.dart';
import 'api_service.dart';

class ToolSettingsDialog extends StatefulWidget {
  final String username;
  final List<String> currentTools;

  const ToolSettingsDialog({
    super.key,
    required this.username,
    required this.currentTools,
  });

  @override
  State<ToolSettingsDialog> createState() => _ToolSettingsDialogState();
}

class _ToolSettingsDialogState extends State<ToolSettingsDialog> {
  late Set<String> _selectedTools;
  final ApiService _apiService = ApiService();
  bool _isSaving = false;

  @override
  void initState() {
    super.initState();
    _selectedTools = Set.from(widget.currentTools);
  }

  Future<void> _saveSettings() async {
    setState(() {
      _isSaving = true;
    });

    final result = await _apiService.saveUserTools(
      widget.username,
      _selectedTools.toList(),
    );

    if (!mounted) return;

    if (result['success'] == true) {
      Navigator.pop(context, _selectedTools.toList());
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Tool preferences saved successfully!'),
          backgroundColor: Colors.green,
        ),
      );
    } else {
      setState(() {
        _isSaving = false;
      });
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(result['message'] ?? 'Failed to save preferences'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('Configure Tools'),
      content: SizedBox(
        width: 400,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Select which tools you want to see on your home page:',
                style: TextStyle(fontSize: 14),
              ),
              const SizedBox(height: 16),
              ...ToolConfig.availableTools.map((tool) {
                final toolId = tool['id'] as String;
                final toolName = tool['name'] as String;
                final toolIcon = tool['icon'] as IconData;
                final toolColor = tool['color'] as Color;

                return CheckboxListTile(
                  title: Row(
                    children: [
                      Icon(toolIcon, size: 24, color: toolColor),
                      const SizedBox(width: 12),
                      Text(toolName),
                    ],
                  ),
                  value: _selectedTools.contains(toolId),
                  onChanged: _isSaving
                      ? null
                      : (bool? value) {
                          setState(() {
                            if (value == true) {
                              _selectedTools.add(toolId);
                            } else {
                              _selectedTools.remove(toolId);
                            }
                          });
                        },
                );
              }).toList(),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: _isSaving ? null : () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        ElevatedButton(
          onPressed: _isSaving ? null : _saveSettings,
          style: ElevatedButton.styleFrom(
            backgroundColor: Colors.blue,
            foregroundColor: Colors.white,
          ),
          child: _isSaving
              ? const SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: Colors.white,
                  ),
                )
              : const Text('Save'),
        ),
      ],
    );
  }
}
