# Flo

A framework for writing full stack applications with Flutter and Go.

# Prototype (actively developing)

Right now, all this will do is generate an example API with example Auth middleware for the Supabase integration that will eventually be auto-generatable in the front end, then generate a basic Flutter app with no framework-specific modifications. The roadmap is to initially gain parity with the functionality of [Flython](https://github.com/thecodekitchen/flython). Once that is achieved, I intend to expand into other auto-generatable modules for a variety of common web application tasks that involve coordinating front and back end development for both Flython and Flo. 

Overall, the purpose of these frameworks is to encourage standardized integration patterns that can be implemented across entire codebases simultaneously, eliminating the most tedious sorts of development drag. Schema mismatches don't need to happen anymore. We have the technology.

# Instructions

Clone this project to a directory in your GOPATH.
Then run 
```
go build .
```
from the cloned directory and try it out by running
```
./flo create flo_test
```
NOTE: Must be a valid Dart package name, so no hyphens or spaces, only underscores.

For now, it's creating the app directory in the parent folder of your cloned directory. This is so that the generated project isn't interpreted as a module within the Flo package. This is provisional and will likely change as the strategy develops.