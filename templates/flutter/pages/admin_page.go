package pages

func AdminPageBytes() []byte {
	return []byte(
		`import 'package:flutter/material.dart';
import 'package:http/http.dart';
import 'package:supabase_flutter/supabase_flutter.dart';
import 'package:flutter_barcode_scanner/flutter_barcode_scanner.dart';
import 'dart:convert';
import '../backend.dart';
import '../models.dart';

class AdminPage extends StatefulWidget {
	const AdminPage({super.key});

	@override
	State<AdminPage> createState() => _AdminPageState();
}

class _AdminPageState extends State<AdminPage> {
	SupabaseClient client = Supabase.instance.client;
	bool sessionIsValid = false;
	Backend backend =
		Backend(client: Supabase.instance.client, deviceIp: '192.168.0.9');
	String? _backendResponse;
	String? _barcodeScanResponse;
	void getUserMetadata() {
		setState(() {
			_backendResponse =
				backend.client!.auth.currentUser!.userMetadata.toString();
		});
	}

	void testBackend() async {
		String response = await backend.GET(route: 'admin');
		setState(() => _backendResponse = response);
	}

	Future<void> scanBook() async {
		String isbn = await FlutterBarcodeScanner.scanBarcode(
			'#000000', 'Cancel Scanning', true, ScanMode.BARCODE);
		Map libraryData =
			await get(Uri.parse("https://openlibrary.org/isbn/$isbn.json"))
				.then((res) => jsonDecode(res.body));
		String title = libraryData['title'];
		setState(() {
			_barcodeScanResponse = isbn;
			_backendResponse = title;
		});
	}

	Future<void> getBookData() async {
		String isbn = _barcodeScanResponse!;
		Map libraryData =
			await get(Uri.parse("https://openlibrary.org/isbn/$isbn.json"))
				.then((res) => jsonDecode(res.body));
		String title = libraryData['title'];
		setState(() {
			_backendResponse = title;
		});
	}

	Future<void> saveBook() async {
		String seller = client.auth.currentUser!.email!;
		String response = await backend.POST(
			Book(isbn: _barcodeScanResponse!, seller: seller),
			route: 'save_book');
		setState(() {
			_backendResponse = response;
		});
	}

	@override
	void initState() {
		print(client.auth.currentUser!.userMetadata);
		if (client.auth.currentUser != null &&
			client.auth.currentSession != null &&
			client.auth.currentUser!.userMetadata?['admin'] == true) {
			sessionIsValid = true;
		}
		super.initState();
	}

	@override
	Widget build(BuildContext context) {
		return Scaffold(
			appBar: AppBar(title: Text(sessionIsValid ? "Admin" : 'Unauthorized!')),
			body: sessionIsValid
				? Column(
					mainAxisAlignment: MainAxisAlignment.center,
					children: [
						const Text('Congrats! You are an admin!'),
						FloatingActionButton(
						onPressed: () async {
							await scanBook();
						},
						child: const Text('scan book'),
						),
						FloatingActionButton(
							onPressed: getBookData, child: const Text('get book')),
						FloatingActionButton(
							onPressed: () async {
							await saveBook();
							},
							child: const Text('test admin backend')),
						_backendResponse != null
							? Text(_backendResponse!)
							: const Placeholder(),
						_barcodeScanResponse != null
							? Text(_barcodeScanResponse!)
							: const SizedBox(
								height: 1,
							)
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
