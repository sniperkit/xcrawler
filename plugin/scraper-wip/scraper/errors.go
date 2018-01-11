package scraper

import (
	"errors"
)

var (
	// ErrInvalidConfig alert whenever you try to use something that is not a structure in the Save
	// function, or something that is not a pointer to a structure in the Load funciton
	ErrInvalidConfig = errors.New("etcetera: configuration must be a structure or a pointer to a structure")

	// ErrNotInitialized alert when you pass a structure to the Load function that has a map attribute
	// that is not initialized
	ErrNotInitialized = errors.New("etcetera: configuration has fields that are not initialized (map)")

	// ErrFieldNotMapped alert whenever you try to access a field that wasn't loaded in the client
	// structure. If we don't load the field before we cannot determinate the path or version
	ErrFieldNotMapped = errors.New("etcetera: trying to retrieve information of a field that wasn't previously loaded")

	// ErrFieldNotAddr is throw when a field that cannot be addressable is used in a place that we
	// need the pointer to identify the path related to the field
	ErrFieldNotAddr = errors.New("etcetera: field must be a pointer or an addressable value")
)

// https://github.com/coreos/etcd/blob/master/error/error.go
const (
	etcdErrorCodeKeyNotFound  etcdErrorCode = 100 // used in tests
	etcdErrorCodeNotFile      etcdErrorCode = 102 // used in tests
	etcdErrorCodeNodeExist    etcdErrorCode = 105
	etcdErrorCodeRaftInternal etcdErrorCode = 300 // used in tests
)

type etcdErrorCode int
