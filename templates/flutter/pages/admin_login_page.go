package pages

func AdminLoginPageBytes() []byte {
	return []byte(
		`import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:supabase_auth_ui/supabase_auth_ui.dart';

class AdminSignUpPage extends StatelessWidget {
	const AdminSignUpPage({Key? key}) : super(key: key);

	@override
	Widget build(BuildContext context) {
		return Scaffold(
			appBar: AppBar(title: const Text('Sign Up / Sign In')),
			body: ListView(
				padding: const EdgeInsets.all(24.0),
				children: [
					SupaEmailAuth(
					redirectTo: kIsWeb ? null : 'bookup2323://projects',
					onSignInComplete: (response) {
						if (response.user?.userMetadata?['admin'] != true) {
						showDialog<void>(
							context: context,
							builder: (BuildContext context) {
							return const AlertDialog(
								title: Text('AlertDialog Title'),
								content: SingleChildScrollView(
								child: ListBody(
									children: <Widget>[
									Text('You are not an authorized admin.'),
									],
								),
								),
								actions: [],
							);
							},
						);
						}
						Navigator.of(context).pushReplacementNamed('/admin');
					},
					onSignUpComplete: (response) async {
						await Supabase.instance.client.auth
							.updateUser(UserAttributes(data: {'admin': true}));
						print(Supabase.instance.client.auth.currentUser!.userMetadata
							.toString());
						Navigator.of(context).pushReplacementNamed('/admin');
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
					const Text('Logging in or signing up as a regular user?'),
					FloatingActionButton(
						onPressed: () => Navigator.pushNamed(context, '/'),
						child: const Text('click here'))
				],
			),
		);
	}
}
		`)
}
