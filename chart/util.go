package chart

// Sugar for internal consistency tests, which looks slightly
// nicer than if (!condition) { panic(...) }
func invariant(shouldBeTrue bool) {
	if (!shouldBeTrue) {
		panic("Internal consistency check failed.")
	}
}
