package validator

import "fmt"

const (
	// ErrorAlreadyBeenTaken is duplicate unique key error text.
	ErrorAlreadyBeenTaken string = "レコードは既に登録済みです"
	// ErrorForeignKeyConstraintFailed is a foreign key constraint fails error text.
	ErrorForeignKeyConstraintFailed string = "関連するレコードが存在しません"
)

// ErrorRequired returns error text the filed must exist.
func ErrorRequired(field string) string {
	return fmt.Sprintf("%sは必須項目です", field)
}

// ErrorHexcolor returns error text that the field nust hexcolor.
func ErrorHexcolor(field string) string {
	return fmt.Sprintf("%sは16進数で入力してください", field)
}

// ErrorTooLong returns error text that the field is less than a param.
func ErrorTooLong(field, param string) string {
	return fmt.Sprintf("%sは%s文字以下で入力してください", field, param)
}

// ErrorTooShort returns error text that the field is more than a param.
func ErrorTooShort(field, param string) string {
	return fmt.Sprintf("%sは%s文字以上で入力してください", field, param)
}

// ErrorEqualField returns error text that the field is equal a param.
func ErrorEqualField(field, param string) string {
	return fmt.Sprintf("%sと%sの値は一致する必要があります", field, param)
}
