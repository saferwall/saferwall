package gib

// ngrams returns a list of all n-grams of length n in a given string s.
func ngrams(s string, n int) []string {

	var ngrams []string = make([]string, 0, len(s))

	for i := 0; i < len(s)-n+1; i++ {
		ngrams = append(ngrams, s[i:i+n])
	}

	return ngrams
}
