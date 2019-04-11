package database

import (
	fmt "fmt"
	time "time"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
	github_com_johnnyeven_libtools_sqlx "github.com/johnnyeven/libtools/sqlx"
	github_com_johnnyeven_libtools_sqlx_builder "github.com/johnnyeven/libtools/sqlx/builder"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

var ChannelTable *github_com_johnnyeven_libtools_sqlx_builder.Table

func init() {
	ChannelTable = DBIn2Book.Register(&Channel{})
}

func (channel *Channel) D() *github_com_johnnyeven_libtools_sqlx.Database {
	return DBIn2Book
}

func (channel *Channel) T() *github_com_johnnyeven_libtools_sqlx_builder.Table {
	return ChannelTable
}

func (channel *Channel) TableName() string {
	return "t_channel"
}

type ChannelFields struct {
	ID           *github_com_johnnyeven_libtools_sqlx_builder.Column
	ChannelID    *github_com_johnnyeven_libtools_sqlx_builder.Column
	Name         *github_com_johnnyeven_libtools_sqlx_builder.Column
	ClientID     *github_com_johnnyeven_libtools_sqlx_builder.Column
	ClientSecret *github_com_johnnyeven_libtools_sqlx_builder.Column
	AuthURL      *github_com_johnnyeven_libtools_sqlx_builder.Column
	TokenURL     *github_com_johnnyeven_libtools_sqlx_builder.Column
	CreateTime   *github_com_johnnyeven_libtools_sqlx_builder.Column
	UpdateTime   *github_com_johnnyeven_libtools_sqlx_builder.Column
	Enabled      *github_com_johnnyeven_libtools_sqlx_builder.Column
}

var ChannelField = struct {
	ID           string
	ChannelID    string
	Name         string
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	CreateTime   string
	UpdateTime   string
	Enabled      string
}{
	ID:           "ID",
	ChannelID:    "ChannelID",
	Name:         "Name",
	ClientID:     "ClientID",
	ClientSecret: "ClientSecret",
	AuthURL:      "AuthURL",
	TokenURL:     "TokenURL",
	CreateTime:   "CreateTime",
	UpdateTime:   "UpdateTime",
	Enabled:      "Enabled",
}

func (channel *Channel) Fields() *ChannelFields {
	table := channel.T()

	return &ChannelFields{
		ID:           table.F(ChannelField.ID),
		ChannelID:    table.F(ChannelField.ChannelID),
		Name:         table.F(ChannelField.Name),
		ClientID:     table.F(ChannelField.ClientID),
		ClientSecret: table.F(ChannelField.ClientSecret),
		AuthURL:      table.F(ChannelField.AuthURL),
		TokenURL:     table.F(ChannelField.TokenURL),
		CreateTime:   table.F(ChannelField.CreateTime),
		UpdateTime:   table.F(ChannelField.UpdateTime),
		Enabled:      table.F(ChannelField.Enabled),
	}
}

func (channel *Channel) IndexFieldNames() []string {
	return []string{"ChannelID", "ID"}
}

func (channel *Channel) ConditionByStruct() *github_com_johnnyeven_libtools_sqlx_builder.Condition {
	table := channel.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(channel)

	conditions := []*github_com_johnnyeven_libtools_sqlx_builder.Condition{}

	for _, fieldName := range channel.IndexFieldNames() {
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

func (channel *Channel) PrimaryKey() github_com_johnnyeven_libtools_sqlx.FieldNames {
	return github_com_johnnyeven_libtools_sqlx.FieldNames{"ID"}
}
func (channel *Channel) UniqueIndexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{"U_channel_id": github_com_johnnyeven_libtools_sqlx.FieldNames{"ChannelID", "Enabled"}}
}
func (channel *Channel) Comments() map[string]string {
	return map[string]string{
		"AuthURL":      "认证URL",
		"ChannelID":    "业务ID",
		"ClientID":     "ClientID",
		"ClientSecret": "ClientSecret",
		"CreateTime":   "",
		"Enabled":      "",
		"ID":           "",
		"Name":         "名称",
		"TokenURL":     "交换tokenURL",
		"UpdateTime":   "",
	}
}

func (channel *Channel) Create(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if channel.CreateTime.IsZero() {
		channel.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	channel.UpdateTime = channel.CreateTime

	stmt := channel.D().
		Insert(channel).
		Comment("Channel.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		channel.ID = uint64(lastInsertID)
	}

	return err
}

func (channel *Channel) DeleteByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (err error) {
	table := channel.T()

	stmt := table.Delete().
		Comment("Channel.DeleteByStruct").
		Where(channel.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (channel *Channel) CreateOnDuplicateWithUpdateFields(db *github_com_johnnyeven_libtools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if channel.CreateTime.IsZero() {
		channel.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	channel.UpdateTime = channel.CreateTime

	table := channel.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(channel, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range channel.UniqueIndexes() {
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
		Comment("Channel.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (channel *Channel) FetchByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()
	stmt := table.Select().
		Comment("Channel.FetchByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(channel.ID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	return db.Do(stmt).Scan(channel).Err()
}

func (channel *Channel) FetchByIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()
	stmt := table.Select().
		Comment("Channel.FetchByIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(channel.ID),
			table.F("Enabled").Eq(channel.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(channel).Err()
}

func (channel *Channel) DeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()
	stmt := table.Delete().
		Comment("Channel.DeleteByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(channel.ID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	return db.Do(stmt).Scan(channel).Err()
}

func (channel *Channel) UpdateByIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Channel.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(channel.ID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	dbRet := db.Do(stmt).Scan(channel)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return channel.FetchByID(db)
	}
	return nil
}

func (channel *Channel) UpdateByIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(channel, zeroFields...)
	return channel.UpdateByIDWithMap(db, fieldValues)
}

func (channel *Channel) SoftDeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Channel.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(channel.ID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	dbRet := db.Do(stmt).Scan(channel)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return channel.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (channel *Channel) FetchByChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()
	stmt := table.Select().
		Comment("Channel.FetchByChannelID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ChannelID").Eq(channel.ChannelID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	return db.Do(stmt).Scan(channel).Err()
}

func (channel *Channel) FetchByChannelIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()
	stmt := table.Select().
		Comment("Channel.FetchByChannelIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ChannelID").Eq(channel.ChannelID),
			table.F("Enabled").Eq(channel.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(channel).Err()
}

func (channel *Channel) DeleteByChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()
	stmt := table.Delete().
		Comment("Channel.DeleteByChannelID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ChannelID").Eq(channel.ChannelID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	return db.Do(stmt).Scan(channel).Err()
}

func (channel *Channel) UpdateByChannelIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("Channel.UpdateByChannelIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ChannelID").Eq(channel.ChannelID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	dbRet := db.Do(stmt).Scan(channel)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return channel.FetchByChannelID(db)
	}
	return nil
}

func (channel *Channel) UpdateByChannelIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(channel, zeroFields...)
	return channel.UpdateByChannelIDWithMap(db, fieldValues)
}

func (channel *Channel) SoftDeleteByChannelID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	channel.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := channel.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("Channel.SoftDeleteByChannelID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ChannelID").Eq(channel.ChannelID),
			table.F("Enabled").Eq(channel.Enabled),
		))

	dbRet := db.Do(stmt).Scan(channel)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return channel.DeleteByChannelID(db)
		}
		return err
	}
	return nil
}

type ChannelList []Channel

// deprecated
func (channelList *ChannelList) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (count int32, err error) {
	*channelList, count, err = (&Channel{}).FetchList(db, size, offset, conditions...)
	return
}

func (channel *Channel) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (channelList ChannelList, count int32, err error) {
	channelList = ChannelList{}

	table := channel.T()

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Channel.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(github_com_johnnyeven_libtools_sqlx_builder.Count(github_com_johnnyeven_libtools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&channelList).Err()

	return
}

func (channel *Channel) List(db *github_com_johnnyeven_libtools_sqlx.DB, condition *github_com_johnnyeven_libtools_sqlx_builder.Condition) (channelList ChannelList, err error) {
	channelList = ChannelList{}

	table := channel.T()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Channel.List").
		Where(condition)

	err = db.Do(stmt).Scan(&channelList).Err()

	return
}

func (channel *Channel) ListByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (channelList ChannelList, err error) {
	channelList = ChannelList{}

	table := channel.T()

	condition := channel.ConditionByStruct()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Channel.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&channelList).Err()

	return
}

// deprecated
func (channelList *ChannelList) BatchFetchByChannelIDList(db *github_com_johnnyeven_libtools_sqlx.DB, channelIDList []uint64) (err error) {
	*channelList, err = (&Channel{}).BatchFetchByChannelIDList(db, channelIDList)
	return
}

func (channel *Channel) BatchFetchByChannelIDList(db *github_com_johnnyeven_libtools_sqlx.DB, channelIDList []uint64) (channelList ChannelList, err error) {
	if len(channelIDList) == 0 {
		return ChannelList{}, nil
	}

	table := channel.T()

	condition := table.F("ChannelID").In(channelIDList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Channel.BatchFetchByChannelIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&channelList).Err()

	return
}

// deprecated
func (channelList *ChannelList) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (err error) {
	*channelList, err = (&Channel{}).BatchFetchByIDList(db, idList)
	return
}

func (channel *Channel) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (channelList ChannelList, err error) {
	if len(idList) == 0 {
		return ChannelList{}, nil
	}

	table := channel.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("Channel.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&channelList).Err()

	return
}
