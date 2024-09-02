## Parser challenges
A significant portion of our time and effort in building the interpreter was devoted to constructing the parser. We have to carefully consider how to represent each node for our language features, and considerable energy was spent on devising strategies to determine which nodes to create by peeking ahead at tokens.

One of the main reasons for this extensive effort was our lack of prior experience in parsing source text into an Abstract Syntax Tree (AST). We had to allocate time to learning and understanding the concepts, such as the expression parsing algorithm, which took us five days to fully grasp. Consequently, someone with more experience in the syntactic phase of language implementation could likely build the parser more efficiently. To illustrate, out of seven pull requests, five were dedicated to creating the parser.

For readers developing a domain-specific language (DSL) where the primary focus is on the language's behaviour rather than its syntax, especially if the syntax is similar to existing languages, it may be more efficient to leverage existing open-source parsers. This approach could abstract away much of the syntactic phase, allowing you to focus on the development of the behaviour of the language.


talk about the benefit of tests