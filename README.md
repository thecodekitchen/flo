# Flo

A framework for writing full stack applications with Flutter and Go.

# Installation

For the moment, the best way to install is to build from source. Just clone the repo into a directory on your gopath, run 'go build .' from its root where the go.mod file is, and move the resulting binary into an executable directory like /usr/local/bin on Unix systems. This assumes that you have both Flutter and Go installed on your system already. In order to use the deployment features in full, you will also need to have Docker and Kind installed on your system. The instructions for those features assume a basic knowledge of those tools.

# Alpha stage (active development)

The architecture of a Flo application, whether or not Go is chosen as the backend language, will consist of three essential components:

1. A Flutter frontend application that can compile to multiple platform-native code bases for use on a wide range of devices including Android, ios, Windows, Linux, Mac, and JS for web.

2. A backend API in a language of the user's choice*. 

3. A set of JSON models, the format for which is covered below.

* currently only Go, but Python should be coming soon. See [Flython](https://github.com/thecodekitchen/flython) for an idea of what this looks like with a Python-based CLI. I decided to move to Go for the main root CLI implementation of the integration pattern for all subsequent language models that need to be generated due to the modular nature of the architectures it enabled.

# Model Synchronization
Changes to the model schema should be made to a json file called, by default, "models.json" and applied with the 'flo sync' command like so:

```
flo sync
```

Add the filename of any alternative model files you might want to sync (for A/B testing or other purposes) as an extra argument.

Rule 1:

You can synchronize the model schema for the front and back end in this way with any JSON file, allowing for 
versioned backups of prior configurations. However, breaking changes can occur if either the front or back end 
model files are edited independently. The whole point is to never do this. If you need a unique class in the front
or back end logic, that's fine. Just don't extend them from the BaseModel class in either code base. Also, don't include them in the models file for either codebase or else they will be written over and deleted when the models file is re-generated.

Those class extensions operate as indicators that the front and back end need to synchronize those classes in a shared schema.

** To reiterate, any class that doesn't extend those base classes should be declared outside the models.py and models.dart files respectively.**

Don't worry, that was the most complicated rule.

Rule 2:

Any optional attribute names for your model should be prefixed with a '?' in the JSON specification like so:
```
{
   "User": {
      "name": "string",
      "?age": "int"
}
```

Rule 3:

Any list attributes should be suffixed with '[]' like so:
```
{
   "User": {
      "name": "string",
      "aliases[]": "string"
   }
}
```

Optional list attributes should combine the two syntaxes:
```
{
   "User": {
      "name": "string",
      "?aliases[]": "string"
   }
}
```
In the last case, the 'aliases' attribute of the User model will be treated in the generated model files as an optional List of strings on the User class which extends the the BaseModel class.

Rule 4:

Don't nest objects! Just create new ones and reference them. For instance, instead of this:

**Bad! Very bad!**
```
{
   "User": {
      "name",
      "documents[]": [
         {
            "name": "string",
            "content": "string"
         }
      ]
   }
}
```

do this:

**nice**
```
{
   "Document": {
      "name": "string",
      "content": "string"
   },
   "User": {
      "name": "string",
      "documents[]": "Document"
   }
}
```

flo will automatically interpret this as an array (list) of 'Document' objects. To avoid compatibility bugs, the only complex (structured) data types that should be declared as attribute types in models should be other declared models. The simple data types that are currently supported include "string", "int", "float", and "bool". More granular data type support is on the roadmap, but the goal of simple, readable models seems to be mostly served within these limitations.

Rule 5:

Models need to be referenced AFTER they are declared. The reason that 'Document' came first in the above example is because the 'User' model references 'Document' in one of its attributes. Only some languages require this, but it is applied as a general rule in order to ensure maximum compatibility.

Rule 6:

I know five was a nice round number, but I decided things would go a lot smoother if I added class extensions. Now we can say
```
{
   "Document": {
      "name": "string",
      "content": "string"
   },
   "User": {
      "name": "string",
      "documents[]": "Document"
   },
   "Project: {
      "name": "string"
   },
   "FloUser(User)": {
      "projects[]": "Project
   }
}
```
and the FloUser class will be an extension of the User class containing one extra property, projects, which is a list of Project instances.

# Supabase Authentication

For now, the Supabase Auth modules are baked into the core of the framework, but the backend middleware is turned off initially so that test calls (from external sources such as Postman or curl commands) can go through smoothly to your backend. Currently, in order to effectively login to the example app, you need to create a free Supabase project and enter its anon key and project url into the generated .env files in your front and back end directories (there is one in each of them!). The project name specified in your create command is used by default as the redirect url scheme, so you will also need to go to Auth/Url Configuration on the Supabase console and add 
```
<your-project-name>://home
``` 
to the list of allowed redirect urls.

Only then will the login functionality work on both mobile and web. Windows and Linux integrations are still in development.

The beta version will involve a base implementation with no auth baked in with the current implementation available through an optional flag, but it will require some refactoring. My goal was to build the more complicated case first and strip it down for the base implementation.

Other auth modules are in the roadmap including, but not limited to, Firebase, Auth0, and FusionAuth.

# Deployment

Flo automatically Dockerizes any backend codebase it generates. In addition, if you provide a valid dockerhub registry with the '-registry' flag in the create command, it will give you a Kubernetes manifest file for pulling the back end API's image from your registry and deploying it to a cluster. If you have installed Docker and Kind on your development system, you can run the backend in a local cluster alongside a TiKV-equipped SurrealDB instance for data persistence. This is automated in the form of the following command.
```
flo deploy -kind
```
The roadmap includes modules for performing this implementation with multiple cloud providers' Kubernetes hosting services via Terraform and Github Actions, the prototype for which can be found [here](https://github.com/thecodekitchen/terraform-k8s-example).

# Objectives

Overall, the purpose of these frameworks is to encourage standardized integration patterns that can be implemented across entire codebases simultaneously, eliminating the most tedious sorts of development drag. Schema mismatches don't need to happen anymore. We have the technology. This pattern should work well for programmers who are accustomed to MVC architectures with firmly established and synchronized data models that are fed to and from the controllers plugged into the user views.

The models are generated as established native data structures in both the back and front end codebases simultaneously so that the data can be translated effectively between the two languages. I've chosen Flutter/Dart as the primary frontend technology for its versatility to reach multiple platforms in their native paradigms, which is in line with my core motivation to provide idiomatic translations of a core architecture for as many platforms and languages as possible.

The future direction of development will likely be centered around understanding how the integration schemes employed in Flython and Flo, respectively, might work with the data structures available in other backend languages. A couple of promising targets on the horizon are Rust and Julia. More suggestions are extremely welcome. I'm excited to use this as a platform for expanding my knowledge base of backend languages.
