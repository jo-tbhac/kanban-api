package repository

const (
	// ErrorRecordNotFound is record not found error text.
	ErrorRecordNotFound string = "該当するレコードが見つかりませんでした"
	// ErrorInvalidSession is invalid session error text.
	ErrorInvalidSession string = "認証情報が不正です"
	// ErrorUserDoesNotExist is does not exist an user error text.
	ErrorUserDoesNotExist string = "該当するユーザーが存在しません"
	// ErrorInvalidPassword is invalid password error text.
	ErrorInvalidPassword string = "パスワードが一致しません"
	// ErrorAuthenticationFailed is an error text if an error occurs during authentication.
	ErrorAuthenticationFailed string = "認証中に問題が発生しました"
	// ErrorInvalidRequest is invalid request error text.
	ErrorInvalidRequest string = "リククエストが不正です"
	// ErrorUnavailableTestUser is when test user is unavailable error text.
	ErrorUnavailableTestUser string = "このテストユーザーは使用中です"
)
