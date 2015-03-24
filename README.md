Solution Discussion
=================================

## Architecture

domain/
	common types and interfaces
interfaces/
	data layer
	rest interface
usecases/
	main business logic
usescases/lexer 
	parses HTML snippets

### Domain

Domain contains common types and interfaces that are shared across the entire 
project. The main shared interface being NamesRepository that defines how 
data (in this case name/URL pairs) are accessed.

### Interfaces

The interfaces layer handles serializing to/from wire formats and 
storage formats

#### Data

The current data layer (mapBasedNameRepositoryImpl) is simple and implemented
using a map and implements the NamesRepository interface.

A read/write mutex is used to ensure consistancy.  All read and write requests
are blocked when the map is being manipulated (written to).  This will cause
some contention but will enable the required consistency.  Since typical usage
patterns are more read intensivethis contention should not cause an issue.  If 
it is know that a specific implementation will have a heavy write pattern,
reevaluting whether to use a simple read/write mutex maybe needed.

The space complexity using a map is O(n).
The time complexity of insert and lookup is O(1).

Currently, the number of names/urls stored is limited by the memory of the computer.  If the number of name/url pairs exceeds available memory a new access layer can be built as long as it conforms to the NamesRepository interface.  Postgres, Redis or another persistance layer could be used, this
will affect the stated time and space complexity.

Also, if it is known that the name will be long and have a reasonable amount
of overlapping prefixes a prefix trie (http://en.wikipedia.org/wiki/Trie) can
be used to improve space complexity by sacrificing time complexity. 

#### Rest

GET		/names/{name:[A-Za-z0-9]+}
	Implemented by namesResourceImpl.RetrieveName()

	Complexity: O(1), underlying implementation uses a map, constant time lookup
	Space: O(n), underlying implementation uses a map.


PUT		/names/{name:[A-Za-z0-9]+}
	Implemented by namesResourceImpl.UpdateURLForName()

	Complexity: O(1), underlying implementation uses a map, constant time insert
	Space: O(n), underlying implementation uses a map.


DELETE	/names
	Implemented by namesResourceImpl.RemoveAllNames()

	Complexity: O(n), underlying implementation uses a map, need to iterate map to delete all records
	Space: O(n), underlying implementation uses a map.


POST	/annotate
	Implemented by annotateInteractorImpl.AnnotateHTML(), using a custom lexer

	Complexity: O(n), Lexer requires a single pass of the string
	Space: O(2n)->O(n), Using Go slices, no part of the string being lexed is
		copied. To build the final annotated string a string buffer is created
		to make appending more efficient and handle adding the additional
		hyperlinks. 



### Usecases

The usecases holds all business logic of the application (interactors).
Ideally all input validation is performed here.


#### Lexer

A basic Lexer is used to parse the HTML strings for annotation.  The
lexer is based off the Golang template parser implementation as
described by Rob Pike (https://www.youtube.com/watch?v=HxaD_trXwRE).

It parses the string in linear time and returns slices for processing
meaning that no memory is copied or allocated. There is a bit of overhead
with the internal channel that is used to communicate lexemes which
could be removed by using a buffer, but the channel implementation is
cleaner.  Additional memory is only allocated in a string buffer to
rebuild the the HTML snippet with the added hyperlinks.

The Lexer consists of a state machine which makes implementation
rather simple and makes extending it trivial.  In fact, I
initially read the problem description incorrectly and my first
implementation did not handle arbitrary HTML tags.  Adding an
additional state to correctly handle HTML tags was relatively easy.
This extensibility makes up for the initial complexity.