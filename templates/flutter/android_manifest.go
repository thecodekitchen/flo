package flutter

func AndroidManifestBytes(scheme string) []byte {
	return []byte(
		`<manifest xmlns:android="http://schemas.android.com/apk/res/android">
	<application
		android:label="` + scheme + `_fe"
		android:name="${{applicationName}}"
		android:icon="@mipmap/ic_launcher">
		<activity
			android:name=".MainActivity"
			android:exported="true"
			android:launchMode="singleTop"
			android:theme="@style/LaunchTheme"
			android:configChanges="orientation|keyboardHidden|keyboard|screenSize|smallestScreenSize|locale|layoutDirection|fontScale|screenLayout|density|uiMode"
			android:hardwareAccelerated="true"
			android:windowSoftInputMode="adjustResize">
			<!-- Specifies an Android theme to apply to this Activity as soon as
					the Android process has started. This theme is visible to the user
					while the Flutter UI initializes. After that, this theme continues
					to determine the Window background behind the Flutter UI. -->
			<meta-data
				android:name="io.flutter.embedding.android.NormalTheme"
				android:resource="@style/NormalTheme"
			/>
			<intent-filter>
				<action android:name="android.intent.action.MAIN"/>
				<category android:name="android.intent.category.LAUNCHER"/>
			</intent-filter>
			<intent-filter>
				<action android:name="android.intent.action.VIEW" />
				<category android:name="android.intent.category.DEFAULT" />
				<category android:name="android.intent.category.BROWSABLE" />
				<data
				android:scheme="` + scheme + `"
				android:host="home" />
			</intent-filter>
		</activity>
		<!-- Don't delete the meta-data below.
				This is used by the Flutter tool to generate GeneratedPluginRegistrant.java -->
		<meta-data
			android:name="flutterEmbedding"
			android:value="2" />
	</application>
</manifest>`)
}
