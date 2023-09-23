gvecd
=====

A novel vector database for K-dimensional similarity searches on arbitrarily
large vector spaces.

Instead of using trees, gvecd sorts subdivisions of the vector space in memory
by euclidean distance and performs a simple binary search to find nearest
neighbors. It serializes blocks of vector data in a binary format which allows
binary search to be performed with seek operations, so parallel searches on
disk or in memory can be performed on many blocks simultaneously.

Is this good? Does it perform as well as other solutions? It seems fast.

I think it has potential.
