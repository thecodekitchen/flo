package flutter

func ModelsFileBytes() []byte {
	return []byte(
		`class BaseModel {
	String? type = 'BaseModel';
	BaseModel({this.type = 'BaseModel'});
	Map toJson() => {
		'type': 'BaseModel'
	};
}

class User extends BaseModel {
	String name;
	int age;

	User({
	super.type = 'User',
		required this.name,
		required this.age,
	});

	@override
	Map toJson() {
		return {
		'type': 'User',
		'name': name,
		'age': age,
		};
	}

	static User fromJson(Map json) {
		return User(
			type: 'User',
			name: json['name'],
			age: json['age'],
		);
	}            
}

class Book extends BaseModel {
	String isbn;
	String seller;

	Book({
	super.type = 'Book',
		required this.isbn,
		required this.seller,
	});

	@override
	Map toJson() {
		return {
		'type': 'Book',
		'isbn': isbn,
		'seller': seller,
		};
	}

	static Book fromJson(Map json) {
		return Book(
			type: 'Book',
			isbn: json['isbn'],
			seller: json['seller'],
		);
	}            
}

BaseModel parseModel(Map<dynamic, dynamic> modelMap) {
	switch (modelMap['type']) {
		case 'User':
			{
				return User.fromJson(modelMap);
			}
		case 'Book':
			{
				return Book.fromJson(modelMap);
			}
		default:
		throw const FormatException('invalid model map');
	}
}`)
}
