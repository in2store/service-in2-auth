package database

import (
	fmt "fmt"
	time "time"

	github_com_johnnyeven_libtools_courier_enumeration "github.com/johnnyeven/libtools/courier/enumeration"
	github_com_johnnyeven_libtools_sqlx "github.com/johnnyeven/libtools/sqlx"
	github_com_johnnyeven_libtools_sqlx_builder "github.com/johnnyeven/libtools/sqlx/builder"
	github_com_johnnyeven_libtools_timelib "github.com/johnnyeven/libtools/timelib"
)

var StateMapTable *github_com_johnnyeven_libtools_sqlx_builder.Table

func init() {
	StateMapTable = DBIn2Book.Register(&StateMap{})
}

func (stateMap *StateMap) D() *github_com_johnnyeven_libtools_sqlx.Database {
	return DBIn2Book
}

func (stateMap *StateMap) T() *github_com_johnnyeven_libtools_sqlx_builder.Table {
	return StateMapTable
}

func (stateMap *StateMap) TableName() string {
	return "t_state_map"
}

type StateMapFields struct {
	ID         *github_com_johnnyeven_libtools_sqlx_builder.Column
	State      *github_com_johnnyeven_libtools_sqlx_builder.Column
	ChannelID  *github_com_johnnyeven_libtools_sqlx_builder.Column
	CreateTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	UpdateTime *github_com_johnnyeven_libtools_sqlx_builder.Column
	Enabled    *github_com_johnnyeven_libtools_sqlx_builder.Column
}

var StateMapField = struct {
	ID         string
	State      string
	ChannelID  string
	CreateTime string
	UpdateTime string
	Enabled    string
}{
	ID:         "ID",
	State:      "State",
	ChannelID:  "ChannelID",
	CreateTime: "CreateTime",
	UpdateTime: "UpdateTime",
	Enabled:    "Enabled",
}

func (stateMap *StateMap) Fields() *StateMapFields {
	table := stateMap.T()

	return &StateMapFields{
		ID:         table.F(StateMapField.ID),
		State:      table.F(StateMapField.State),
		ChannelID:  table.F(StateMapField.ChannelID),
		CreateTime: table.F(StateMapField.CreateTime),
		UpdateTime: table.F(StateMapField.UpdateTime),
		Enabled:    table.F(StateMapField.Enabled),
	}
}

func (stateMap *StateMap) IndexFieldNames() []string {
	return []string{"ID", "State"}
}

func (stateMap *StateMap) ConditionByStruct() *github_com_johnnyeven_libtools_sqlx_builder.Condition {
	table := stateMap.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(stateMap)

	conditions := []*github_com_johnnyeven_libtools_sqlx_builder.Condition{}

	for _, fieldName := range stateMap.IndexFieldNames() {
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

func (stateMap *StateMap) PrimaryKey() github_com_johnnyeven_libtools_sqlx.FieldNames {
	return github_com_johnnyeven_libtools_sqlx.FieldNames{"ID"}
}
func (stateMap *StateMap) UniqueIndexes() github_com_johnnyeven_libtools_sqlx.Indexes {
	return github_com_johnnyeven_libtools_sqlx.Indexes{"U_state": github_com_johnnyeven_libtools_sqlx.FieldNames{"State", "Enabled"}}
}
func (stateMap *StateMap) Comments() map[string]string {
	return map[string]string{
		"ChannelID":  "ChannelID",
		"CreateTime": "",
		"Enabled":    "",
		"ID":         "",
		"State":      "业务ID",
		"UpdateTime": "",
	}
}

func (stateMap *StateMap) Create(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if stateMap.CreateTime.IsZero() {
		stateMap.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	stateMap.UpdateTime = stateMap.CreateTime

	stmt := stateMap.D().
		Insert(stateMap).
		Comment("StateMap.Create")

	dbRet := db.Do(stmt)
	err := dbRet.Err()

	if err == nil {
		lastInsertID, _ := dbRet.LastInsertId()
		stateMap.ID = uint64(lastInsertID)
	}

	return err
}

func (stateMap *StateMap) DeleteByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (err error) {
	table := stateMap.T()

	stmt := table.Delete().
		Comment("StateMap.DeleteByStruct").
		Where(stateMap.ConditionByStruct())

	err = db.Do(stmt).Err()
	return
}

func (stateMap *StateMap) CreateOnDuplicateWithUpdateFields(db *github_com_johnnyeven_libtools_sqlx.DB, updateFields []string) error {
	if len(updateFields) == 0 {
		panic(fmt.Errorf("must have update fields"))
	}

	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	if stateMap.CreateTime.IsZero() {
		stateMap.CreateTime = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}
	stateMap.UpdateTime = stateMap.CreateTime

	table := stateMap.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(stateMap, updateFields...)

	delete(fieldValues, "ID")

	cols, vals := table.ColumnsAndValuesByFieldValues(fieldValues)

	m := make(map[string]bool, len(updateFields))
	for _, field := range updateFields {
		m[field] = true
	}

	// fields of unique index can not update
	delete(m, "CreateTime")

	for _, fieldNames := range stateMap.UniqueIndexes() {
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
		Comment("StateMap.CreateOnDuplicateWithUpdateFields")

	return db.Do(stmt).Err()
}

func (stateMap *StateMap) FetchByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()
	stmt := table.Select().
		Comment("StateMap.FetchByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(stateMap.ID),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	return db.Do(stmt).Scan(stateMap).Err()
}

func (stateMap *StateMap) FetchByIDForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()
	stmt := table.Select().
		Comment("StateMap.FetchByIDForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(stateMap.ID),
			table.F("Enabled").Eq(stateMap.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(stateMap).Err()
}

func (stateMap *StateMap) DeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()
	stmt := table.Delete().
		Comment("StateMap.DeleteByID").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(stateMap.ID),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	return db.Do(stmt).Scan(stateMap).Err()
}

func (stateMap *StateMap) UpdateByIDWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("StateMap.UpdateByIDWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(stateMap.ID),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	dbRet := db.Do(stmt).Scan(stateMap)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return stateMap.FetchByID(db)
	}
	return nil
}

func (stateMap *StateMap) UpdateByIDWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(stateMap, zeroFields...)
	return stateMap.UpdateByIDWithMap(db, fieldValues)
}

func (stateMap *StateMap) SoftDeleteByID(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("StateMap.SoftDeleteByID").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("ID").Eq(stateMap.ID),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	dbRet := db.Do(stmt).Scan(stateMap)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return stateMap.DeleteByID(db)
		}
		return err
	}
	return nil
}

func (stateMap *StateMap) FetchByState(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()
	stmt := table.Select().
		Comment("StateMap.FetchByState").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("State").Eq(stateMap.State),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	return db.Do(stmt).Scan(stateMap).Err()
}

func (stateMap *StateMap) FetchByStateForUpdate(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()
	stmt := table.Select().
		Comment("StateMap.FetchByStateForUpdate").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("State").Eq(stateMap.State),
			table.F("Enabled").Eq(stateMap.Enabled),
		)).
		ForUpdate()

	return db.Do(stmt).Scan(stateMap).Err()
}

func (stateMap *StateMap) DeleteByState(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()
	stmt := table.Delete().
		Comment("StateMap.DeleteByState").
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("State").Eq(stateMap.State),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	return db.Do(stmt).Scan(stateMap).Err()
}

func (stateMap *StateMap) UpdateByStateWithMap(db *github_com_johnnyeven_libtools_sqlx.DB, fieldValues github_com_johnnyeven_libtools_sqlx_builder.FieldValues) error {

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()

	delete(fieldValues, "ID")

	stmt := table.Update().
		Comment("StateMap.UpdateByStateWithMap").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("State").Eq(stateMap.State),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	dbRet := db.Do(stmt).Scan(stateMap)
	err := dbRet.Err()
	if err != nil {
		return err
	}

	rowsAffected, _ := dbRet.RowsAffected()
	if rowsAffected == 0 {
		return stateMap.FetchByState(db)
	}
	return nil
}

func (stateMap *StateMap) UpdateByStateWithStruct(db *github_com_johnnyeven_libtools_sqlx.DB, zeroFields ...string) error {
	fieldValues := github_com_johnnyeven_libtools_sqlx.FieldValuesFromStructByNonZero(stateMap, zeroFields...)
	return stateMap.UpdateByStateWithMap(db, fieldValues)
}

func (stateMap *StateMap) SoftDeleteByState(db *github_com_johnnyeven_libtools_sqlx.DB) error {
	stateMap.Enabled = github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE

	table := stateMap.T()

	fieldValues := github_com_johnnyeven_libtools_sqlx_builder.FieldValues{}
	fieldValues["Enabled"] = github_com_johnnyeven_libtools_courier_enumeration.BOOL__FALSE

	if _, ok := fieldValues["UpdateTime"]; !ok {
		fieldValues["UpdateTime"] = github_com_johnnyeven_libtools_timelib.MySQLTimestamp(time.Now())
	}

	stmt := table.Update().
		Comment("StateMap.SoftDeleteByState").
		Set(table.AssignsByFieldValues(fieldValues)...).
		Where(github_com_johnnyeven_libtools_sqlx_builder.And(
			table.F("State").Eq(stateMap.State),
			table.F("Enabled").Eq(stateMap.Enabled),
		))

	dbRet := db.Do(stmt).Scan(stateMap)
	err := dbRet.Err()
	if err != nil {
		dbErr := github_com_johnnyeven_libtools_sqlx.DBErr(err)
		if dbErr.IsConflict() {
			return stateMap.DeleteByState(db)
		}
		return err
	}
	return nil
}

type StateMapList []StateMap

// deprecated
func (stateMapList *StateMapList) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (count int32, err error) {
	*stateMapList, count, err = (&StateMap{}).FetchList(db, size, offset, conditions...)
	return
}

func (stateMap *StateMap) FetchList(db *github_com_johnnyeven_libtools_sqlx.DB, size int32, offset int32, conditions ...*github_com_johnnyeven_libtools_sqlx_builder.Condition) (stateMapList StateMapList, count int32, err error) {
	stateMapList = StateMapList{}

	table := stateMap.T()

	condition := github_com_johnnyeven_libtools_sqlx_builder.And(conditions...)

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("StateMap.FetchList").
		Where(condition)

	errForCount := db.Do(stmt.For(github_com_johnnyeven_libtools_sqlx_builder.Count(github_com_johnnyeven_libtools_sqlx_builder.Star()))).Scan(&count).Err()
	if errForCount != nil {
		err = errForCount
		return
	}

	stmt = stmt.Limit(size).Offset(offset)

	stmt = stmt.OrderDescBy(table.F("CreateTime"))

	err = db.Do(stmt).Scan(&stateMapList).Err()

	return
}

func (stateMap *StateMap) List(db *github_com_johnnyeven_libtools_sqlx.DB, condition *github_com_johnnyeven_libtools_sqlx_builder.Condition) (stateMapList StateMapList, err error) {
	stateMapList = StateMapList{}

	table := stateMap.T()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("StateMap.List").
		Where(condition)

	err = db.Do(stmt).Scan(&stateMapList).Err()

	return
}

func (stateMap *StateMap) ListByStruct(db *github_com_johnnyeven_libtools_sqlx.DB) (stateMapList StateMapList, err error) {
	stateMapList = StateMapList{}

	table := stateMap.T()

	condition := stateMap.ConditionByStruct()

	condition = github_com_johnnyeven_libtools_sqlx_builder.And(condition, table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("StateMap.ListByStruct").
		Where(condition)

	err = db.Do(stmt).Scan(&stateMapList).Err()

	return
}

// deprecated
func (stateMapList *StateMapList) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (err error) {
	*stateMapList, err = (&StateMap{}).BatchFetchByIDList(db, idList)
	return
}

func (stateMap *StateMap) BatchFetchByIDList(db *github_com_johnnyeven_libtools_sqlx.DB, idList []uint64) (stateMapList StateMapList, err error) {
	if len(idList) == 0 {
		return StateMapList{}, nil
	}

	table := stateMap.T()

	condition := table.F("ID").In(idList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("StateMap.BatchFetchByIDList").
		Where(condition)

	err = db.Do(stmt).Scan(&stateMapList).Err()

	return
}

// deprecated
func (stateMapList *StateMapList) BatchFetchByStateList(db *github_com_johnnyeven_libtools_sqlx.DB, stateList []string) (err error) {
	*stateMapList, err = (&StateMap{}).BatchFetchByStateList(db, stateList)
	return
}

func (stateMap *StateMap) BatchFetchByStateList(db *github_com_johnnyeven_libtools_sqlx.DB, stateList []string) (stateMapList StateMapList, err error) {
	if len(stateList) == 0 {
		return StateMapList{}, nil
	}

	table := stateMap.T()

	condition := table.F("State").In(stateList)

	condition = condition.And(table.F("Enabled").Eq(github_com_johnnyeven_libtools_courier_enumeration.BOOL__TRUE))

	stmt := table.Select().
		Comment("StateMap.BatchFetchByStateList").
		Where(condition)

	err = db.Do(stmt).Scan(&stateMapList).Err()

	return
}
