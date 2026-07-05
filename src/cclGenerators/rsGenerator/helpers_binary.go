package rsGenerator

func rustScalarByteSize(rustType string) string {
	switch rustType {
	case "i8", "u8":
		return "1"
	case "i16", "u16":
		return "2"
	case "i32", "u32", "f32":
		return "4"
	case "i64", "u64", "f64":
		return "8"
	default:
		return "4"
	}
}
