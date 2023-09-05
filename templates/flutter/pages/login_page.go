package pages

func LoginPageBytes() []byte {
	return []byte(
		`import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:supabase_auth_ui/supabase_auth_ui.dart';

class SignUpPage extends StatelessWidget {
	const SignUpPage({Key? key}) : super(key: key);

	@override
	Widget build(BuildContext context) {
		return Scaffold(
			appBar: AppBar(title: const Text('Sign Up / Sign In')),
			body: ListView(
				padding: const EdgeInsets.all(24.0),
				children: [
					SupaEmailAuth(
					redirectTo: kIsWeb ? null : 'bookup2323://home',
					onSignInComplete: (response) {
						Navigator.of(context).pushReplacementNamed('/home');
					},
					onSignUpComplete: (response) {
						Navigator.of(context).pushReplacementNamed('/home');
					},
					metadataFields: [
						MetaDataField(
						prefixIcon: const Icon(Icons.person),
						label: 'Username',
						key: 'username',
						validator: (val) {
							if (val == null || val.isEmpty) {
							return 'Please enter something';
							}
							return null;
						},
						),
					],
					),
					const Divider(),
					const Text('Logging in or signing up as an admin?'),
					FloatingActionButton(
						onPressed: () => Navigator.pushNamed(context, '/admin-login'),
						child: const Text('click here'))
				],
			),
		);
	}
}
		`)
}
