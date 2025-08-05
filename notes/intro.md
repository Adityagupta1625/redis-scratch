Redis is the most popular in-memory key-value store, used primarily for caching as no storage is faster than memory. Caching servers are inevitable because it’s the easiest way to scale. A cache can take slow DB queries out of the equation.

A cache server is a map<string,string> over a network. But in Redis, the “value” part is not just a string, it can be any of the data structures: hash, list, or sorted set, enabling complex use cases such as ranking, list pagination, etc. This is why Redis is called a data structure server.

The first version of Redis was written in the TCL scripting language. But that’s only a toy version of Redis; the final product is in C.
