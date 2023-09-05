package flutter

func MainFileBytes() []byte {
	return []byte(
		`import 'package:flutter/material.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';

import './pages/error_page.dart';
import './pages/login_page.dart';
import 'pages/home_page.dart';
import './pages/admin_page.dart';
import './pages/admin_login_page.dart';

Future<void> main() async {
	await dotenv.load();
	String supabaseUrl = dotenv.get('SUPABASE_URL');
	String supabaseApiKey = dotenv.get('SUPABASE_API_KEY');
	await Supabase.initialize(url: supabaseUrl, anonKey: supabaseApiKey);
	runApp(const MyApp());
}

class MyApp extends StatefulWidget {
	const MyApp({super.key});

	@override
	State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
	SupabaseClient client = Supabase.instance.client;
	@override
	Widget build(BuildContext context) {
		return MaterialApp(
			title: 'Flo Demo',
			theme: ThemeData(
			colorScheme: ColorScheme.fromSeed(seedColor: Colors.indigo),
			useMaterial3: true,
			),
			initialRoute: '/',
			routes: {
			'/': (context) => const SignUpPage(),
			'/error': (context) => const ErrorPage(error: ''),
			'/home': (context) => const HomePage(),
			'/admin': (context) => const AdminPage(),
			'/admin-login': (context) => const AdminSignUpPage()
			},
		);
	}
}
		`)
}
