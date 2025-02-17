package generate

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// dictionary base64字典：a-z, A-Z, 0-9, +, /
var dictionary = [64]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
	'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '+', '/',
}

// RandomBase64 生成随机Base64编码的字符串
//   - length: 长度
func RandomBase64(length int) string {
	// 以uint64为基准单位生成随机数，对应base64编码后为10位
	generateLength := length/10 + (9+(length%10))/10
	randomUint64 := make([]uint64, generateLength)
	for i := 0; i < generateLength; i++ {
		randomUint64[i] = rand.Uint64()
	}

	// 生成随机字符串
	builder := strings.Builder{}
	for i, j, k := 0, 0, 0; i < length && j < 10; i, j = i+1, j+1 {
		builder.WriteRune(dictionary[(randomUint64[k]>>(j*6))&0x3F])
		if j == 9 {
			k++
			j = -1
		}
	}

	return builder.String()
}

// RandomBase62 生成随机Base62编码的字符串
//   - length: 长度
func RandomBase62(length int) string {
	// 以uint64为基准单位生成随机数，对应base64编码后为10位
	generateLength := length/10 + (9+(length%10))/10
	randomUint64 := make([]uint64, generateLength)
	for i := 0; i < generateLength; i++ {
		randomUint64[i] = rand.Uint64()
	}

	builder := strings.Builder{}
	for i, j, k := 0, 0, 0; i < length && j < 10; i, j = i+1, j+1 {
		builder.WriteRune(dictionary[(randomUint64[k]>>(j*6))&0x3D])
		if j == 9 {
			k++
			j = -1
		}
	}

	return builder.String()
}

// RandomBase64WithPrefix 生成Base64编码的随机字符串，包含前缀
//   - prefix: 前缀
//   - totalLength: 总长度
func RandomBase64WithPrefix(prefix string, totalLength int) string {
	length := totalLength - len(prefix)
	generateLength := length/10 + (9+(length%10))/10
	randomUint64 := make([]uint64, generateLength)
	for i := 0; i < generateLength; i++ {
		randomUint64[i] = rand.Uint64()
	}

	builder := strings.Builder{}
	builder.WriteString(prefix)

	for i, j, k := 0, 0, 0; i < length && j < 10; i, j = i+1, j+1 {
		builder.WriteRune(dictionary[(randomUint64[k]>>(j*6))&0x3F])
		if j == 9 {
			k++
			j = -1
		}
	}

	return builder.String()
}

// RandomBase62WithPrefix 生成Base62编码的随机字符串，包含前缀
//   - prefix: 前缀
//   - totalLength: 总长度
func RandomBase62WithPrefix(prefix string, totalLength int) string {
	length := totalLength - len(prefix)
	generateLength := length/10 + (9+(length%10))/10
	randomUint64 := make([]uint64, generateLength)
	for i := 0; i < generateLength; i++ {
		randomUint64[i] = rand.Uint64()
	}

	builder := strings.Builder{}
	builder.WriteString(prefix)

	for i, j, k := 0, 0, 0; i < length && j < 10; i, j = i+1, j+1 {
		builder.WriteRune(dictionary[(randomUint64[k]>>(j*6))&0x3D])
		if j == 9 {
			k++
			j = -1
		}
	}

	return builder.String()
}

// RandomSixDigitNumberCode 生成6位随机数字验证码
// 示例输出： 000001, 000002, 000003...
func RandomSixDigitNumberCode() string {
	number := strconv.Itoa(rand.Intn(999999)) //nolint:gosec
	for len(number) < 6 {
		number = "0" + number
	}
	return number
}

// RandomFourDigitNumberCode 生成4位随机数字验证码
// 示例输出： 0001, 0002, 0003...
func RandomFourDigitNumberCode() string {
	number := strconv.Itoa(rand.Intn(9999)) //nolint:gosec
	for len(number) < 4 {
		number = "0" + number
	}
	return number
}

// RandomSixDigitNumberCodeWithPrefix 生成带前缀的6位随机数字验证码
//   - prefix: 前缀
//
// 示例输出： prefix000001, prefix000002, prefix000003...
func RandomSixDigitNumberCodeWithPrefix(prefix string) string {
	return prefix + RandomSixDigitNumberCode()
}

// RandomFourDigitNumberCodeWithPrefix 生成带前缀的4位随机数字验证码
//   - prefix: 前缀
//
// 示例输出： prefix0001, prefix0002, prefix0003...
func RandomFourDigitNumberCodeWithPrefix(prefix string) string {
	return prefix + RandomFourDigitNumberCode()
}

// UUID 生成UUID
func UUID() string {
	return uuid.NewString()
}

// TraceID 生成TraceID
func TraceID() string {
	return strings.ReplaceAll(UUID(), "-", "")
}
