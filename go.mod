module github.com/rocketlaunchr/unsafe

go 1.24.0

retract (
	v1.0.0 // Bug discovered.
	v1.0.1 // Bug discovered.
    v1.1.0 // Bug discovered.
	v1.2.0 // SetZeros should not have been exported.
	v1.2.1 // Too dangerous to use in production.
	v1.2.2 // Too dangerous to use in production.
	v1.3.0 // Too dangerous to use in production.
)