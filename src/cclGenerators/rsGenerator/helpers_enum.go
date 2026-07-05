package rsGenerator

func rustSerdeIntegerMethod(baseType string) string {
	switch baseType {
	case "i8":
		return "i8"
	case "i16":
		return "i16"
	case "i32":
		return "i32"
	case "i64":
		return "i64"
	case "u8":
		return "u8"
	case "u16":
		return "u16"
	case "u32":
		return "u32"
	case "u64":
		return "u64"
	default:
		return "i32"
	}
}
