//go:build pro && enterprise

package main

func init() {
	features = append(features,
		"Enterprise Feature #1",
		"Enterprise Feature #2",
	)
}
