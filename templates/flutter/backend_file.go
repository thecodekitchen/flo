package flutter

func BackendFileBytes() []byte {
	return []byte(
		`import 'dart:convert';

		import 'package:device_info_plus/device_info_plus.dart';
		import 'package:flutter/foundation.dart';
		import 'package:http/http.dart';
		import 'package:supabase_flutter/supabase_flutter.dart';
		
		import './models.dart';
		
		class Backend {
		  SupabaseClient? client;
		  String? deviceIp;
		  String? deploymentUrl;
		
		  Backend({this.client, this.deviceIp, this.deploymentUrl});
		
		  Future<bool> detectEmulator() async {
			final DeviceInfoPlugin deviceInfo = DeviceInfoPlugin();
			bool? isEmulator;
		
			if (defaultTargetPlatform == TargetPlatform.android) {
			  final androidInfo = await deviceInfo.androidInfo;
			  isEmulator = !androidInfo.isPhysicalDevice;
			} else if (defaultTargetPlatform == TargetPlatform.iOS) {
			  final iosInfo = await deviceInfo.iosInfo;
			  isEmulator = !iosInfo.isPhysicalDevice;
			} else {
			  isEmulator = false;
			}
			return isEmulator;
		  }
		
		  Future<String> baseUrl() async {
			if (deploymentUrl != null) {
			  return deploymentUrl!;
			}
		
			final bool isEmulator = await detectEmulator();
			if (isEmulator) {
			  if (defaultTargetPlatform == TargetPlatform.android) {
				return 'http://10.0.2.2:8000';
			  }
			  if (defaultTargetPlatform == TargetPlatform.iOS) {
				return 'http://127.0.0.1:8000';
			  }
			}
			if (kIsWeb) {
			  return 'http://localhost:8000';
			}
			// if deviceIp is null at this point in execution,
			// platform is likely desktop, so localhost still works.
			return 'http://${deviceIp ?? 'localhost'}:8000';
		  }
		
		  Future<String> POST(BaseModel model,
			  {String? route, Map<String, String>? headers}) async {
			String baseUrlString = await baseUrl();
			Uri base = Uri.parse(baseUrlString);
		
			Map<String, String> postHeaders =
				Map.fromEntries([const MapEntry('Content-Type', 'application/json')]);
			if (client != null) {
			  try {
				final String accessToken = client!.auth.currentSession!.accessToken;
				postHeaders
					.addEntries([MapEntry('Authorization', 'Bearer $accessToken')]);
			  } catch (err) {
				print(err);
			  }
			}
			if (headers != null) {
			  postHeaders.addAll(headers);
			}
			if (route == null) {
			  Response response =
				  await post(base, headers: postHeaders, body: jsonEncode(model));
			  return response.body;
			} else {
			  Uri url = Uri.parse('$baseUrlString/$route');
			  Response response =
				  await post(url, headers: postHeaders, body: jsonEncode(model));
			  return response.body;
			}
		  }
		
		  Future<String> GET({String? route, Map<String, String>? headers}) async {
			String baseUrlString = await baseUrl();
			Uri base = Uri.parse(baseUrlString);
			Map<String, String> getHeaders = Map.fromEntries([]);
			if (client != null) {
			  try {
				final String accessToken = client!.auth.currentSession!.accessToken;
				getHeaders
					.addEntries([MapEntry('Authorization', 'Bearer $accessToken')]);
			  } catch (err) {
				print(err);
			  }
			}
			if (headers != null) {
			  getHeaders.addAll(headers);
			}
			if (route == null) {
			  Response response = await get(base, headers: getHeaders);
			  return response.body;
			} else {
			  Uri url = Uri.parse('$baseUrlString/$route');
			  Response response = await get(url, headers: getHeaders);
			  return response.body;
			}
		  }
		
		  Future<String> DELETE({String? route, Map<String, String>? headers}) async {
			String baseUrlString = await baseUrl();
			Uri base = Uri.parse(baseUrlString);
			Map<String, String> deleteHeaders = Map.fromEntries([]);
			if (client != null) {
			  try {
				final String accessToken = client!.auth.currentSession!.accessToken;
				deleteHeaders
					.addEntries([MapEntry('Authorization', 'Bearer $accessToken')]);
			  } catch (err) {
				print(err);
			  }
			}
			if (headers != null) {
			  deleteHeaders.addAll(headers);
			}
			if (route == null) {
			  Response response = await delete(base, headers: deleteHeaders);
			  return response.body;
			} else {
			  Uri url = Uri.parse('$baseUrlString/$route');
			  Response response = await delete(url, headers: deleteHeaders);
			  return response.body;
			}
		  }
		
		  Future<String> PUT(BaseModel model,
			  {String? route, Map<String, String>? headers}) async {
			String baseUrlString = await baseUrl();
			Uri base = Uri.parse(baseUrlString);
			Map<String, String> putHeaders =
				Map.fromEntries([const MapEntry('Content-Type', 'application/json')]);
			if (client != null) {
			  try {
				final String accessToken = client!.auth.currentSession!.accessToken;
				putHeaders
					.addEntries([MapEntry('Authorization', 'Bearer $accessToken')]);
			  } catch (err) {
				print(err);
			  }
			}
			if (headers != null) {
			  putHeaders.addAll(headers);
			}
			if (route == null) {
			  Response response =
				  await put(base, headers: putHeaders, body: jsonEncode(model));
			  return response.body;
			} else {
			  Uri url = Uri.parse('$baseUrlString/$route');
			  Response response =
				  await post(url, headers: putHeaders, body: jsonEncode(model));
			  return response.body;
			}
		  }
		}`)
}
