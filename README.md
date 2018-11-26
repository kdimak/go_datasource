# What is it

Different implementations of key-value stores. There are:
* `Database` - database store - the slowest store;
* `DistributedCache` - distributed cache store - working much faster than database store;
* `LocalCache` - local store - additional local layer between a user and distributed cache store, which uses LRU cache to
provide values for frequently-requested keys instantly. There is an assumption that `Database` is read-only.

# Implementation details
Original `Database` and `DistributedCache` structs are made safe for concurrent reads/writes. This reflects more realistically
the nature of databases / remote caches.
There are two local cache implementations:
1) `LocalCache` - first implementation, simply using LRU cache and making subsequent datasource calls: local cache -> remote cache -> database
2) `FastLocalCache` - the fastest (to satisfy the lowest possible latency requirement) but more complex implementation, we gain the speed by
a) parallel query of remote cache / database; b) by multiplexing cache waiting queries.

The latter implementation is used in main.go and perf tests.

# What is missed
1) Unit test of `FastLocalCache` - lack of time (but initial impl. `LocalCache` is covered).
 