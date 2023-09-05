package pages

func ErrorPageBytes() []byte {
	return []byte(
		`import 'package:flutter/material.dart';

class ErrorPage extends StatelessWidget {
	const ErrorPage({super.key, required this.error});

	final String error;

	@override
	Widget build(BuildContext context) {
		return Scaffold(
			appBar: AppBar(
			title: const Text("ERROR"),
			backgroundColor: Colors.red,
			),
			body: Text(error),
		);
	}
}
		`)
}
