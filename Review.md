# Major Changes

**Better Input Handling** - Switched from `fmt.Scan()` to `bufio.Reader` to properly handle user input

**Type Safety** - Added a Provider type with constants instead of comparing strings

**Improved Validation** - Simplified phone number validation by working directly with strings instead of converting to int slices

**Constants** - Defined constants for magic numbers and URLs for easier maintenance