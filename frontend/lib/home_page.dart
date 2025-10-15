import 'package:flutter/material.dart';
import 'api_service.dart';
import 'tool_widget.dart';
import 'tool_settings_dialog.dart';

class HomePage extends StatefulWidget {
  final String username;

  const HomePage({super.key, required this.username});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final ApiService _apiService = ApiService();
  List<String> _userTools = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadUserTools();
  }

  Future<void> _loadUserTools() async {
    setState(() {
      _isLoading = true;
    });

    final result = await _apiService.getUserTools(widget.username);

    if (mounted) {
      setState(() {
        _isLoading = false;
        if (result['success'] == true) {
          _userTools = List<String>.from(result['tools'] ?? []);
        }
      });
    }
  }

  Future<void> _openSettings() async {
    final result = await showDialog<List<String>>(
      context: context,
      builder: (context) => ToolSettingsDialog(
        username: widget.username,
        currentTools: _userTools,
      ),
    );

    if (result != null) {
      setState(() {
        _userTools = result;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('tools'),
        backgroundColor: Colors.blue,
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            tooltip: 'Settings',
            onPressed: _openSettings,
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            tooltip: 'Logout',
            onPressed: () {
              Navigator.of(context).popUntil((route) => route.isFirst);
            },
          ),
        ],
      ),
      body: _isLoading
          ? const Center(child: CircularProgressIndicator())
          : _userTools.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(
                        Icons.widgets_outlined,
                        size: 80,
                        color: Colors.grey,
                      ),
                      const SizedBox(height: 24),
                      Text(
                        'Welcome, ${widget.username}!',
                        style: const TextStyle(
                          fontSize: 28,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'No tools configured yet.',
                        style: TextStyle(
                          fontSize: 16,
                          color: Colors.grey[600],
                        ),
                      ),
                      const SizedBox(height: 24),
                      ElevatedButton.icon(
                        onPressed: _openSettings,
                        icon: const Icon(Icons.settings),
                        label: const Text('Configure Tools'),
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.blue,
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(
                            horizontal: 24,
                            vertical: 16,
                          ),
                        ),
                      ),
                    ],
                  ),
                )
              : Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Welcome, ${widget.username}!',
                        style: const TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 24),
                      Expanded(
                        child: GridView.builder(
                          gridDelegate: const SliverGridDelegateWithMaxCrossAxisExtent(
                            maxCrossAxisExtent: 200,
                            childAspectRatio: 1,
                            crossAxisSpacing: 16,
                            mainAxisSpacing: 16,
                          ),
                          itemCount: _userTools.length,
                          itemBuilder: (context, index) {
                            final toolId = _userTools[index];
                            final toolConfig = ToolConfig.getToolById(toolId);

                            if (toolConfig == null) {
                              return const SizedBox.shrink();
                            }

                            return ToolWidget(
                              name: toolConfig['name'] as String,
                              icon: toolConfig['icon'] as IconData,
                              color: toolConfig['color'] as Color,
                            );
                          },
                        ),
                      ),
                    ],
                  ),
                ),
    );
  }
}
