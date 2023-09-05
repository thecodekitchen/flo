package pages

func HomePageBytes() []byte {
	return []byte(
		`import 'package:flutter/material.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import '../backend.dart';

class HomePage extends StatefulWidget {
	const HomePage({super.key});

	@override
	State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
	SupabaseClient client = Supabase.instance.client;
	bool sessionIsValid = false;
	String? _backend_response;
	String? _username;
	Backend backend = Backend(client: Supabase.instance.client);

	void testBackend() async {
		String response = await backend.GET();
		setState(() => _backend_response = response);
	}

	@override
	void initState() {
		if (client.auth.currentUser != null && client.auth.currentSession != null) {
			sessionIsValid = true;
			_username = client.auth.currentUser!.userMetadata!['username'];
		}
		super.initState();
	}

	@override
	Widget build(BuildContext context) {
		return Scaffold(
			appBar: AppBar(
				title:
					Text(sessionIsValid ? "Welcome, $_username" : 'Unauthorized!')),
			body: sessionIsValid
				? Column(
						mainAxisAlignment: MainAxisAlignment.center,
						children: [
							const Text(
								'Click below to confirm the backend recognizes you.'),
							FloatingActionButton(
								onPressed: testBackend,
								child: const Text('test backend')),
							_backend_response != null
								? Text(_backend_response!)
								: const SizedBox(height: 0)
						],
					)
				: const Column(
					mainAxisAlignment: MainAxisAlignment.center,
					children: [
						Text("You don't have permission to view this page!")
					]));
	}
}
		`)
}
