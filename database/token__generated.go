package database

import (
	fmt "fmt"
	time "time"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
	github_com_johnnyeven_libtools_sqlx "github.com/johnnyeven/libtools/sqlx"
	github_com_johnnyeven_libtools_sqlx_builder "github.com/johnnyeven/libtools/sqlx/builder"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

var TokenTable *github_com_johnnyeven_libtools_sqlx_builder.Table

func init() {
	TokenTable = DBIn2Book.Register(&Token{})
}

func (token *Token) D() *github_com_johnnyeven_libtools_sqlx.Database {
	return DBIn2Book
}

func (token *Token) T() *github_com_johnnyeven_libtools_sqlx_builder.Table {
	return TokenTable
}

func (token *Token) TableName() string {
	return "t_token"
}

type TokenFields struct {
	ID           *github_com_johnnyeven_libtools_sqlx_builder.Column
	TokenID      *github_com_johnnyeven_libtools_sqlx_builder.Column
	UserID       *github_com_johnnyeven_libtools_sqlx_builder.Column
	ChannelID    *github_com_johnnyeven_libtools_sqlx_builder.Column
	AccessToken  *github_com_johnnyeven_libtools_sqlx_builder.Column
	TokenType    *github_com_johnnyeven_libtools_sqlx_builder.Column
	RefreshToken *github_com_johnnyeven_libtools_sqlx_builder.Column
	ExpiryTime   *github_com_johnnyeven_libtools_sqlx_builder.Column
	CreateTime   *github_com_johnnyeven_libtools_sqlx_builder.Column
	UpdateTime   *github_com_johnnyeven_libtools_sqlx_builder.Column
	Enabled      *github_com_johnnyeven_libtools_sqlx_builder.Column
}

var TokenField = struct {
	ID           string
	TokenID      string
	UserID       string
	ChannelID    string
	AccessToken  string
	TokenType    string
	RefreshToken string
	ExpiryTime   string
	CreateTime   string
	UpdateTime   string
	Enabled      string
}{
	ID:           "ID",
	TokenID:      "TokenID",
	UserID:       "UserID",
	ChannelID:    "ChannelID",
	AccessToken:  "AccessToken",
	TokenType:    "TokenType",
	RefreshToken: "RefreshToken",
	ExpiryTime:   "ExpiryTime",
	CreateTime:   "CreateTime",
	UpdateTime:   "UpdateTime",
	Enabled:      "Enabled",
}

func (token *Token) Fields() *TokenFields {
	table := token.T()

	return &TokenFields{
		ID:           table.F(TokenField.ID),
		TokenID:      table.F(TokenField.TokenID),
		UserID:       table.F(TokenField.UserID),
		ChannelID:    table.F(TokenField.ChannelID),
		AccessToken:  table.F(TokenField.AccessToken),
		TokenType:    table.F(TokenField.TokenType),
		RefreshToken: table.F(TokenField.RefreshToken),
		ExpiryTime:   table.F(TokenField.ExpiryTime),
		CreateTime:   table.F(TokenField.CreateTime),
		UpdateTime:   table.F(TokenField.UpdateTime),
		Enabled:      table.F(TokenField.Enabled),
	}
}

func (token *Token) IndexFieldNames() []string {
	return []string{"AccessToken", "ChannelID", "ID", "TokenID", "UserID"}
}

func (token *Token) ConditionByStruct() *github_com_johnnyeven_libtools_sqlx_builder.Condition {
	table := token.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(token)

	conditions := []*github_com_johnnyeven_libtools_sqlx_builder.Condition{}

	for _, fieldName := range token.IndexFieldNames() {
		if v, exists := fieldValues[fieldName]; exists {
			conditions = append(conditions, table.F(fieldName).Eq(v))
			delete(fieldValues, fieldName)
		}
	}

	if len(conditions) == 0 {
		panic(fmt.Errorf("at least one of field for indexes has value"))
	}

	for fieldName, v := range fieldValues {
		conditions = append(conditions, table.F(fieldName).Eq(v))
	}

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	return condition
}

func (token *Token) PrimaryKey() github_com_johnnyeven_libtools_sqlx.FieldNames {
	return github_com_johnnyeven_libtools_sqlx.FieldNames{"ID"}
}
func (token *Token) UniqueIndexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{
		"U_token":        github_com_johnnyeven_libtools_sqlx.FieldNames{"AccessToken", "Enabled"},
		"U_token_id":     github_com_johnnyeven_libtools_sqlx.FieldNames{"TokenID", "Enabled"},
		"U_user_channel": github_com_johnnyeven_libtools_sqlx.FieldNames{"UserID", "ChannelID", "Enabled"},
	}
}
func (token *Token) Comments() map[string]string {
	return map[string]string{
		"AccessToken":  "AccessToken is the token that authorizes and authenticates the requests.",
		"ChannelID":    "通道ID",
		"CreateTime":   "",
		"Enabled":      "",
		"ExpiryTime":   "Expiry is the optional expiration time of the access token.",
		"ID":           "",
		"RefreshToken": "RefreshToken is a token that's used by the application(as opposed to the user) to refresh the access token if it expires.",
		"TokenID":      "业务ID",
		"TokenType":    "TokenType is the type of token. The Type method returns either this or \"Bearer\", the default.",
		"UpdateTime":   "",
		"UserID":       "用户ID",
	}
}

func (token *Token) Create(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if token.CreateTime.IsZero() {
		token.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	token.UpdateTime = token.CreateTime

	stmt := token.D().
		Insert(token).
		Comment("Token.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		token.ID = uint64(lastInsertID)
	}

	return err
}

func (token *Token) DeleteByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (err error) {
	table := token.T()

	stmt := table.Delete().
		Comment("Token.DeleteByStruct").
		Where(token.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (token *Token) CreateOnDuplicateWithUpdateFields(db *github_com_johnnyeven_libtools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if token.CreateTime.IsZero() {
		token.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	token.UpdateTime = token.CreateTime

	table := token.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(token, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range token.UniqueIndexes() {
		for _, field := range fieldNames {
			delete(m, field)
		}
	}

	if len(m) == 0 {
		panic(fmt.Errorf("no fields for updates"))
	}

	for field := range fieldValues {
		if !m[field] {
			delete(fieldValues, field)
		}
	}

	stmt := table.
		Insert().Columns(cols).Values(vals...).
		OnDuplicateKeyUpdate(table.AssignsByFieldValues(fieldValues)...).
		Comment("Token.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (token *Token) FetchByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(token.ID),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) FetchByIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(token.ID),
			table.F("Enabled").Eq(token.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) DeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Delete().
		Comment("Token.DeleteByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(token.ID),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) UpdateByIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Token.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(token.ID),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return token.FetchByID(db)
	}
	return nil
}

func (token *Token) UpdateByIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(token, zeroFields...)
	return token.UpdateByIDWithMap(db, fieldValues)
}

func (token *Token) SoftDeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Token.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(token.ID),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return token.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (token *Token) FetchByTokenID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByTokenID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("TokenID").Eq(token.TokenID),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) FetchByTokenIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByTokenIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("TokenID").Eq(token.TokenID),
			table.F("Enabled").Eq(token.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) DeleteByTokenID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Delete().
		Comment("Token.DeleteByTokenID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("TokenID").Eq(token.TokenID),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) UpdateByTokenIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Token.UpdateByTokenIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("TokenID").Eq(token.TokenID),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return token.FetchByTokenID(db)
	}
	return nil
}

func (token *Token) UpdateByTokenIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(token, zeroFields...)
	return token.UpdateByTokenIDWithMap(db, fieldValues)
}

func (token *Token) SoftDeleteByTokenID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Token.SoftDeleteByTokenID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("TokenID").Eq(token.TokenID),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return token.DeleteByTokenID(db)
		}
		return err
	}
	return nil
}

func (token *Token) FetchByAccessToken(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByAccessToken").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("AccessToken").Eq(token.AccessToken),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) FetchByAccessTokenForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByAccessTokenForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("AccessToken").Eq(token.AccessToken),
			table.F("Enabled").Eq(token.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) DeleteByAccessToken(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Delete().
		Comment("Token.DeleteByAccessToken").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("AccessToken").Eq(token.AccessToken),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) UpdateByAccessTokenWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Token.UpdateByAccessTokenWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("AccessToken").Eq(token.AccessToken),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return token.FetchByAccessToken(db)
	}
	return nil
}

func (token *Token) UpdateByAccessTokenWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(token, zeroFields...)
	return token.UpdateByAccessTokenWithMap(db, fieldValues)
}

func (token *Token) SoftDeleteByAccessToken(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Token.SoftDeleteByAccessToken").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("AccessToken").Eq(token.AccessToken),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return token.DeleteByAccessToken(db)
		}
		return err
	}
	return nil
}

func (token *Token) FetchByUserIDAndChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByUserIDAndChannelID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(token.UserID),
			table.F("ChannelID").Eq(token.ChannelID),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) FetchByUserIDAndChannelIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Select().
		Comment("Token.FetchByUserIDAndChannelIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(token.UserID),
			table.F("ChannelID").Eq(token.ChannelID),
			table.F("Enabled").Eq(token.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) DeleteByUserIDAndChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()
	stmt := table.Delete().
		Comment("Token.DeleteByUserIDAndChannelID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(token.UserID),
			table.F("ChannelID").Eq(token.ChannelID),
			table.F("Enabled").Eq(token.Enabled),
		))

	return db.Do(stmt).Scan(token).Err()
}

func (token *Token) UpdateByUserIDAndChannelIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Token.UpdateByUserIDAndChannelIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(token.UserID),
			table.F("ChannelID").Eq(token.ChannelID),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return token.FetchByUserIDAndChannelID(db)
	}
	return nil
}

func (token *Token) UpdateByUserIDAndChannelIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(token, zeroFields...)
	return token.UpdateByUserIDAndChannelIDWithMap(db, fieldValues)
}

func (token *Token) SoftDeleteByUserIDAndChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	token.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := token.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Token.SoftDeleteByUserIDAndChannelID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("UserID").Eq(token.UserID),
			table.F("ChannelID").Eq(token.ChannelID),
			table.F("Enabled").Eq(token.Enabled),
		))

	dbRet := db.Do(stmt).Scan(token)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return token.DeleteByUserIDAndChannelID(db)
		}
		return err
	}
	return nil
}

type TokenList []Token

// deprecated
func (tokenList *TokenList) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (count int32, err error) {
	*tokenList, count, err = (&Token{}).FetchList(db, size, offset, conditions...)
	return
}

func (token *Token) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (tokenList TokenList, count int32, err error) {
	tokenList = TokenList{}

	table := token.T()

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(github_com_johnnyeven_libtools_sqlx_builder.Count(github_com_johnnyeven_libtools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

func (token *Token) List(db *github_com_johnnyeven_libtools_sqlx.DB, condition *github_com_johnnyeven_libtools_sqlx_builder.Condition) (tokenList TokenList, err error) {
	tokenList = TokenList{}

	table := token.T()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.List").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

func (token *Token) ListByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (tokenList TokenList, err error) {
	tokenList = TokenList{}

	table := token.T()

	condition := token.ConditionByStruct()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

// deprecated
func (tokenList *TokenList) BatchFetchByAccessTokenList(db *github_com_johnnyeven_libtools_sqlx.DB, accessTokenList []string) (err error) {
	*tokenList, err = (&Token{}).BatchFetchByAccessTokenList(db, accessTokenList)
	return
}

func (token *Token) BatchFetchByAccessTokenList(db *github_com_johnnyeven_libtools_sqlx.DB, accessTokenList []string) (tokenList TokenList, err error) {
	if len(accessTokenList) == 0 {
		return TokenList{}, nil
	}

	table := token.T()

	condition := table.F("AccessToken").In(accessTokenList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.BatchFetchByAccessTokenList").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

// deprecated
func (tokenList *TokenList) BatchFetchByChannelIDList(db *github_com_johnnyeven_libtools_sqlx.DB, channelIDList []uint64) (err error) {
	*tokenList, err = (&Token{}).BatchFetchByChannelIDList(db, channelIDList)
	return
}

func (token *Token) BatchFetchByChannelIDList(db *github_com_johnnyeven_libtools_sqlx.DB, channelIDList []uint64) (tokenList TokenList, err error) {
	if len(channelIDList) == 0 {
		return TokenList{}, nil
	}

	table := token.T()

	condition := table.F("ChannelID").In(channelIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.BatchFetchByChannelIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

// deprecated
func (tokenList *TokenList) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (err error) {
	*tokenList, err = (&Token{}).BatchFetchByIDList(db, idList)
	return
}

func (token *Token) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (tokenList TokenList, err error) {
	if len(idList) == 0 {
		return TokenList{}, nil
	}

	table := token.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

// deprecated
func (tokenList *TokenList) BatchFetchByTokenIDList(db *github_com_johnnyeven_libtools_sqlx.DB, tokenIDList []uint64) (err error) {
	*tokenList, err = (&Token{}).BatchFetchByTokenIDList(db, tokenIDList)
	return
}

func (token *Token) BatchFetchByTokenIDList(db *github_com_johnnyeven_libtools_sqlx.DB, tokenIDList []uint64) (tokenList TokenList, err error) {
	if len(tokenIDList) == 0 {
		return TokenList{}, nil
	}

	table := token.T()

	condition := table.F("TokenID").In(tokenIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.BatchFetchByTokenIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}

// deprecated
func (tokenList *TokenList) BatchFetchByUserIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userIDList []uint64) (err error) {
	*tokenList, err = (&Token{}).BatchFetchByUserIDList(db, userIDList)
	return
}

func (token *Token) BatchFetchByUserIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userIDList []uint64) (tokenList TokenList, err error) {
	if len(userIDList) == 0 {
		return TokenList{}, nil
	}

	table := token.T()

	condition := table.F("UserID").In(userIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Token.BatchFetchByUserIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&tokenList).Err()

	return
}
