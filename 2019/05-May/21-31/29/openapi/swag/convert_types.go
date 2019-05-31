package swag

// String returns a pointer to of the string value passed in.
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is null.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
