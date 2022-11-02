package cruncher

func Compare(old Domain, new Domain) (bool, error) {
	eq := true

	// Compare domain zone serial
	if old.Serial != new.Serial {
		eq = false
	}

	// Compare nameservers
	if !stringSlicesEqual(old.Nameservers, new.Nameservers) {
		eq = false
	}

	// Compare mailservers
	if !stringSlicesEqual(old.Mailservers, new.Mailservers) {
		eq = false
	}

	// done
	return eq, nil
}

// stringSlicesEqual function to check if slises are equal
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
