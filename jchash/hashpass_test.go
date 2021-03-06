package jchash

import "testing"

type hashPair struct {
	password, hash512 string
}

// Expected hashes calculated at https://hash.online-convert.com/sha512-generator
var hashPairArray = []hashPair{
	hashPair{
		password: "angryMonkey",
		hash512:  "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
	},
	hashPair{
		password: "superSecret",
		hash512:  "pQaPqt7aC/CThmNsO8xnV+nkLfyJoyqpFzGzmvLivIpjmQXnvJqIULCUOpE+H1f3+p9laadfIkvAxMYZTAxnyQ==",
	},
}

var badPair = hashPair{
	password: "angryMonkey",
	hash512:  "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ==",
}

func TestHashPassword(t *testing.T) {
	for _, pair := range hashPairArray {
		result := HashPassword(pair.password)
		if pair.hash512 != result {
			t.Error("Got wrong hash for password: ", pair.password)
		}
	}
	// Negative case
	result := HashPassword(badPair.password)
	if badPair.hash512 == result {
		t.Error("Deliberate bad hash appears to match: ", badPair.password)
	}
}
