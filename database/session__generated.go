package database

import (
	fmt "fmt"
	time "time"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
	github_com_johnnyeven_libtools_sqlx "github.com/johnnyeven/libtools/sqlx"
	github_com_johnnyeven_libtools_sqlx_builder "github.com/johnnyeven/libtools/sqlx/builder"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

var SessionTable *github_com_johnnyeven_libtools_sqlx_builder.Table

func init() {
	SessionTable = DBIn2Book.Register(&Session{})
}

func (session *Session) D() *github_com_johnnyeven_libtools_sqlx.Database {
	return DBIn2Book
}

func (session *Session) T() *github_com_johnnyeven_libtools_sqlx_builder.Table {
	return SessionTable
}

func (session *Session) TableName() string {
	return "t_session"
}

type SessionFields struct {
	ID         *github_com_johnnyeven_libtools_sqlx_builder.Column
	SessionID  *github_com_johnnyeven_libtools_sqlx_builder.Column
	UserID     *github_com_johnnyeven_libtools_sqlx_builder.Column
	ExpireTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	CreateTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	UpdateTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	Enabled    *github_com_johnnyeven_libtools_sqlx_builder.Column
}

var SessionField = struct {
	ID         string
	SessionID  string
	UserID     string
	ExpireTime string
	CreateTime string
	UpdateTime string
	Enabled    string
}{
	ID:         "ID",
	SessionID:  "SessionID",
	UserID:     "UserID",
	ExpireTime: "ExpireTime",
	CreateTime: "CreateTime",
	UpdateTime: "UpdateTime",
	Enabled:    "Enabled",
}

func (session *Session) Fields() *SessionFields {
	table := session.T()

	return &SessionFields{
		ID:         table.F(SessionField.ID),
		SessionID:  table.F(SessionField.SessionID),
		UserID:     table.F(SessionField.UserID),
		ExpireTime: table.F(SessionField.ExpireTime),
		CreateTime: table.F(SessionField.CreateTime),
		UpdateTime: table.F(SessionField.UpdateTime),
		Enabled:    table.F(SessionField.Enabled),
	}
}

func (session *Session) IndexFieldNames() []string {
	return []string{"ID", "SessionID", "UserID"}
}

func (session *Session) ConditionByStruct() *github_com_johnnyeven_libtools_sqlx_builder.Condition {
	table := session.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(session)

	conditions := []*github_com_johnnyeven_libtools_sqlx_builder.Condition{}

	for _, fieldName := range session.IndexFieldNames() {
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

func (session *Session) PrimaryKey() github_com_johnnyeven_libtools_sqlx.FieldNames {
	return github_com_johnnyeven_libtools_sqlx.FieldNames{"ID"}
}
func (session *Session) Indexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{"I_user_id": github_com_johnnyeven_libtools_sqlx.FieldNames{"UserID"}}
}
func (session *Session) UniqueIndexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{"U_session_id": github_com_johnnyeven_libtools_sqlx.FieldNames{"SessionID", "Enabled"}}
}
func (session *Session) Comments() map[string]string {
	return map[string]string{
		"CreateTime": "",
		"Enabled":    "",
		"ExpireTime": "过期时间",
		"ID":         "",
		"SessionID":  "业务ID",
		"UpdateTime": "",
		"UserID":     "用户ID",
	}
}

func (session *Session) Create(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if session.CreateTime.IsZero() {
		session.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	session.UpdateTime = session.CreateTime

	stmt := session.D().
		Insert(session).
		Comment("Session.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		session.ID = uint64(lastInsertID)
	}

	return err
}

func (session *Session) DeleteByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (err error) {
	table := session.T()

	stmt := table.Delete().
		Comment("Session.DeleteByStruct").
		Where(session.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (session *Session) CreateOnDuplicateWithUpdateFields(db *github_com_johnnyeven_libtools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if session.CreateTime.IsZero() {
		session.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	session.UpdateTime = session.CreateTime

	table := session.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(session, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range session.UniqueIndexes() {
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
		Comment("Session.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (session *Session) FetchByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()
	stmt := table.Select().
		Comment("Session.FetchByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(session.ID),
			table.F("Enabled").Eq(session.Enabled),
		))

	return db.Do(stmt).Scan(session).Err()
}

func (session *Session) FetchByIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()
	stmt := table.Select().
		Comment("Session.FetchByIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(session.ID),
			table.F("Enabled").Eq(session.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(session).Err()
}

func (session *Session) DeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()
	stmt := table.Delete().
		Comment("Session.DeleteByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(session.ID),
			table.F("Enabled").Eq(session.Enabled),
		))

	return db.Do(stmt).Scan(session).Err()
}

func (session *Session) UpdateByIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Session.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(session.ID),
			table.F("Enabled").Eq(session.Enabled),
		))

	dbRet := db.Do(stmt).Scan(session)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return session.FetchByID(db)
	}
	return nil
}

func (session *Session) UpdateByIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(session, zeroFields...)
	return session.UpdateByIDWithMap(db, fieldValues)
}

func (session *Session) SoftDeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Session.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(session.ID),
			table.F("Enabled").Eq(session.Enabled),
		))

	dbRet := db.Do(stmt).Scan(session)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return session.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (session *Session) FetchBySessionID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()
	stmt := table.Select().
		Comment("Session.FetchBySessionID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("SessionID").Eq(session.SessionID),
			table.F("Enabled").Eq(session.Enabled),
		))

	return db.Do(stmt).Scan(session).Err()
}

func (session *Session) FetchBySessionIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()
	stmt := table.Select().
		Comment("Session.FetchBySessionIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("SessionID").Eq(session.SessionID),
			table.F("Enabled").Eq(session.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(session).Err()
}

func (session *Session) DeleteBySessionID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()
	stmt := table.Delete().
		Comment("Session.DeleteBySessionID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("SessionID").Eq(session.SessionID),
			table.F("Enabled").Eq(session.Enabled),
		))

	return db.Do(stmt).Scan(session).Err()
}

func (session *Session) UpdateBySessionIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Session.UpdateBySessionIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("SessionID").Eq(session.SessionID),
			table.F("Enabled").Eq(session.Enabled),
		))

	dbRet := db.Do(stmt).Scan(session)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return session.FetchBySessionID(db)
	}
	return nil
}

func (session *Session) UpdateBySessionIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(session, zeroFields...)
	return session.UpdateBySessionIDWithMap(db, fieldValues)
}

func (session *Session) SoftDeleteBySessionID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	session.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := session.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Session.SoftDeleteBySessionID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("SessionID").Eq(session.SessionID),
			table.F("Enabled").Eq(session.Enabled),
		))

	dbRet := db.Do(stmt).Scan(session)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return session.DeleteBySessionID(db)
		}
		return err
	}
	return nil
}

type SessionList []Session

// deprecated
func (sessionList *SessionList) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (count int32, err error) {
	*sessionList, count, err = (&Session{}).FetchList(db, size, offset, conditions...)
	return
}

func (session *Session) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (sessionList SessionList, count int32, err error) {
	sessionList = SessionList{}

	table := session.T()

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Session.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(github_com_johnnyeven_libtools_sqlx_builder.Count(github_com_johnnyeven_libtools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&sessionList).Err()

	return
}

func (session *Session) List(db *github_com_johnnyeven_libtools_sqlx.DB, condition *github_com_johnnyeven_libtools_sqlx_builder.Condition) (sessionList SessionList, err error) {
	sessionList = SessionList{}

	table := session.T()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Session.List").
		Where(condition)

	err = db.Do(stmt).Scan(&sessionList).Err()

	return
}

func (session *Session) ListByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (sessionList SessionList, err error) {
	sessionList = SessionList{}

	table := session.T()

	condition := session.ConditionByStruct()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Session.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&sessionList).Err()

	return
}

// deprecated
func (sessionList *SessionList) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (err error) {
	*sessionList, err = (&Session{}).BatchFetchByIDList(db, idList)
	return
}

func (session *Session) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (sessionList SessionList, err error) {
	if len(idList) == 0 {
		return SessionList{}, nil
	}

	table := session.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Session.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&sessionList).Err()

	return
}

// deprecated
func (sessionList *SessionList) BatchFetchBySessionIDList(db *github_com_johnnyeven_libtools_sqlx.DB, sessionIDList []string) (err error) {
	*sessionList, err = (&Session{}).BatchFetchBySessionIDList(db, sessionIDList)
	return
}

func (session *Session) BatchFetchBySessionIDList(db *github_com_johnnyeven_libtools_sqlx.DB, sessionIDList []string) (sessionList SessionList, err error) {
	if len(sessionIDList) == 0 {
		return SessionList{}, nil
	}

	table := session.T()

	condition := table.F("SessionID").In(sessionIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Session.BatchFetchBySessionIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&sessionList).Err()

	return
}

// deprecated
func (sessionList *SessionList) BatchFetchByUserIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userIDList []uint64) (err error) {
	*sessionList, err = (&Session{}).BatchFetchByUserIDList(db, userIDList)
	return
}

func (session *Session) BatchFetchByUserIDList(db *github_com_johnnyeven_libtools_sqlx.DB, userIDList []uint64) (sessionList SessionList, err error) {
	if len(userIDList) == 0 {
		return SessionList{}, nil
	}

	table := session.T()

	condition := table.F("UserID").In(userIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Session.BatchFetchByUserIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&sessionList).Err()

	return
}
