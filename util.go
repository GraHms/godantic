package godantic

func joinKeys(parent string, child string) string {
	if parent == "" {
		return "" + child
	}
	return parent + "." + child
}
