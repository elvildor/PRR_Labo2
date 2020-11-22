module server

go 1.15

replace common => ./common

replace lamport => ./lamport

require (
	common v0.0.0-00010101000000-000000000000
	lamport v0.0.0-00010101000000-000000000000
)
