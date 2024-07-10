package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/shachar1236/Baasa/database/sqlite/objects"
	"github.com/shachar1236/Baasa/database/types"
)

func (this *SqliteDB) GetCollectionByName(ctx context.Context, name string) (types.Collection, error) {
	res, err := this.objects_queries.GetTableAndFieldsByTableName(ctx, name)

	var collection types.Collection
	if err != nil {
		this.logger.Error("Cant get collection: ", err)
		return collection, err
	}

	if len(res) > 0 {
		collection.ID = res[0].CollectionID
		collection.Name = res[0].TableName
		collection.QueryRulesDirectoryPath = res[0].Queryrulesdirectorypath
		for _, coll := range res {
			if coll.FieldID.Int64 != 0 {
				field := types.TableField{
					ID:           coll.FieldID.Int64,
					CollectionID: coll.CollectionID,
					FieldName:    coll.FieldName.String,
					FieldType:    coll.FieldType.String,
					FieldOptions: types.NullString(coll.FieldOptions),
                    IsForeignKey: coll.IsForeignKey.Bool,
                    FkTableName: types.NullString(coll.FkTableName),
                    FkFieldName: types.NullString(coll.FkFieldName),
				}
				collection.Fields = append(collection.Fields, field)
			}
		}
	}
	return collection, err
}

func (this *SqliteDB) GetCollectionById(ctx context.Context, id int64) (types.Collection, error) {
	logger := this.logger.With("collection_id", id)

	logger.Info("Get collection by id")
	res, err := this.objects_queries.GetTableAndFieldsByTableId(ctx, id)

	var collection types.Collection
	if err != nil {
		this.logger.Error("Cant get collection: " + err.Error())
		return collection, err
	}

	if len(res) > 0 {
		collection.ID = res[0].CollectionID
		collection.Name = res[0].TableName
		collection.QueryRulesDirectoryPath = res[0].Queryrulesdirectorypath
		for _, coll := range res {
			if coll.FieldID.Int64 != 0 {
				field := types.TableField{
					ID:           coll.FieldID.Int64,
					CollectionID: coll.CollectionID,
					FieldName:    coll.FieldName.String,
					FieldType:    coll.FieldType.String,
					FieldOptions: types.NullString(coll.FieldOptions),
                    IsForeignKey: coll.IsForeignKey.Bool,
                    FkTableName: types.NullString(coll.FkTableName),
                    FkFieldName: types.NullString(coll.FkFieldName),
				}
				collection.Fields = append(collection.Fields, field)
			}
		}
	}
	return collection, err
}

func (this *SqliteDB) GetCollections(ctx context.Context) ([]types.Collection, error) {
	res, err := this.objects_queries.GetAllTablesAndFields(ctx)

	collections := []types.Collection{}
	this.logger.Info(fmt.Sprintf("res: %v", res))
	if len(res) > 0 {
		curr_collection := types.Collection{
			ID:                      res[0].CollectionID,
			Name:                    res[0].TableName,
			QueryRulesDirectoryPath: res[0].Queryrulesdirectorypath,
		}

		for _, row := range res {
			if row.CollectionID != curr_collection.ID {
				collections = append(collections, curr_collection)

				curr_collection = types.Collection{
					ID:                      row.CollectionID,
					Name:                    row.TableName,
					QueryRulesDirectoryPath: row.Queryrulesdirectorypath,
				}
			}

			if row.FieldID.Int64 != 0 {
				curr_collection.Fields = append(curr_collection.Fields, types.TableField{
					ID:           row.FieldID.Int64,
					FieldName:    row.FieldName.String,
					FieldType:    row.FieldType.String,
					FieldOptions: types.NullString(row.FieldOptions),
				})
			}
		}

		collections = append(collections, curr_collection)
	}
	return collections, err
}

func (this *SqliteDB) GetBaseCollections() (collections_names []string, err error) {
	collections_names = []string{"users", "collections", "table_fields", "queries"}
	return collections_names, nil
}

func (this *SqliteDB) AddCollection(ctx context.Context) (types.Collection, error) {
	var collection types.Collection
	name := "NewCollection_" + strconv.Itoa(rand.Int())

	// creating table
	create_statment := "CREATE TABLE IF NOT EXISTS " + name + " (\n    id INTEGER PRIMARY KEY\n);"
	fmt.Println(create_statment)
	_, err := this.db.Exec(create_statment)
	if err != nil {
		this.logger.Error("Cant add collection: " + err.Error())
		return collection, err
	}

	// creating object
	res, err := this.objects_queries.CreateCollection(ctx, name)

	if err != nil {
		this.logger.Error("Cant add collection: " + err.Error())
		return collection, err
	}

	collection.ID = res.ID
	collection.Name = res.TableName
	collection.QueryRulesDirectoryPath = res.QueryRulesDirectoryPath
	return collection, err
}

func (this *SqliteDB) dropTable(ctx context.Context, name string) error {
	_, err := this.db.Exec("DROP TABLE " + name + ";")
	return err
}

func (this *SqliteDB) DeleteCollection(ctx context.Context, name string) error {
	// deleting table
	err := this.dropTable(ctx, name)
	if err != nil {
		this.logger.Error("Cant delete collection: " + err.Error())
		return err
	}

	// deleting object
	err = this.objects_queries.DeleteCollection(ctx, name)

	if err != nil {
		this.logger.Error("Cant delete collection: " + err.Error())
		return err
	}

	return nil
}

// fk = key(table.field, contsrains)
func getTableFieldAndConstrainsFromFK(fk string) (table_relation_name string, table_field_relation_name string, constraints string) {
	t := fk[4 : len(fk)-1]
	params := strings.Split(t, ",")
	if len(params) > 0 {
		relation_params := strings.Split(params[0], ".")
		table_relation_name = relation_params[0]
		table_field_relation_name = relation_params[1]
		if len(params) == 2 {
			constraints = params[1]
		} else {
			constraints = ""
		}
	}
	return
}

func getAddTableQueryForCollection(ctx context.Context, collection types.Collection) string {
	var sb strings.Builder
	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString(collection.Name)
	sb.WriteString("\n(id INTEGER PRIMARY KEY,\n")

	var fks []types.TableField

	for _, field := range collection.Fields {
		sb.WriteString(field.FieldName)
		sb.WriteString(" ")

		t := field.FieldType
		if t[:3] == "key" {
			fks = append(fks, field)
			t = "INTEGER"
		}
		sb.WriteString(t)
		sb.WriteString(" ")

		if field.FieldOptions.Valid {
			sb.WriteString(field.FieldOptions.String)
			sb.WriteString(",\n")
		}
	}

	for _, fk := range fks {
		table_relation_name, table_field_relation_name, constraints := getTableFieldAndConstrainsFromFK(fk.FieldType)

		sb.WriteString("FOREIGN KEY(") // collections(id) ON DELETE CASCADE")
		sb.WriteString(fk.FieldName)
		sb.WriteString(") REFERENCES ")
		sb.WriteString(table_relation_name)
		sb.WriteString("(")
		sb.WriteString(table_field_relation_name)
		sb.WriteString(") ")

		// writing contsrains
		sb.WriteString(constraints)
		sb.WriteString("\n")
	}

	// sb.WriteString(");")
	q := sb.String()
	q = q[:len(q)-2]

	return q + ");"
}

func analyze_if_fk(field_type string) (is_fk bool, parts [2]string) {
	if len(field_type) > 6 {
		if field_type[:3] == "key" {
			// inside := filed_type[3:len(filed_type) - 1]
			// p := strings.Split(inside, ".")
			t, f, _ := getTableFieldAndConstrainsFromFK(field_type)
			parts[0] = t
			parts[1] = f
			is_fk = true
			return
		}
	}

	is_fk = false
	return
}
func (this *SqliteDB) SaveCollectionChanges(ctx context.Context, new_collection types.Collection) error {
	logger := this.logger.With("collection_id", new_collection.ID)

	old_collection, err := this.GetCollectionById(ctx, new_collection.ID)
	if err != nil {
		logger.Error("Cant get collection: " + err.Error())
		return err
	}

	if new_collection.Name != old_collection.Name {
		err = this.ChangeCollectionName(ctx, old_collection.Name, new_collection.Name)
		if err != nil {
			err = errors.New("Cant change collection name: " + err.Error())
			logger.Error(err.Error())
			return err
		}
	}

	needs_to_create_new := false
	needs_to_copy := false
	var copy_query strings.Builder
	copy_query.WriteString("INSERT INTO ")
	copy_query.WriteString(new_collection.Name)
	copy_query.WriteString(" (")
	var new_fields_strings []string

	var if_succeds []func()

	for _, new_field := range new_collection.Fields {
		field_exists := false
		for _, old_field := range old_collection.Fields {

			if new_field.ID == old_field.ID {
				field_exists = true

				if new_field.FieldName != old_collection.Name {
					err = this.ChangeFieldName(ctx, new_collection.Name, old_field.FieldName, new_field.FieldName, old_field.ID)
					if err != nil {
						err = errors.New("Cant change field name: " + err.Error())
						logger.Error(err.Error())
						return err
					}
				}

				new_fields_strings = append(new_fields_strings, new_field.FieldName)

				if new_field.FieldType != old_field.FieldType {
					needs_to_create_new = true
					needs_to_copy = true

					is_new_field_fk, new_field_fk_parts := analyze_if_fk(new_field.FieldType)
					is_old_field_fk, _ := analyze_if_fk(old_field.FieldType)

					if_succeds = append(if_succeds, func() {
						var err error
						if !is_new_field_fk {
							err = this.objects_queries.ChangeFieldType(ctx, objects.ChangeFieldTypeParams{
								FieldType: new_field.FieldType,
								ID:        old_field.ID,
							})

						}

						if is_new_field_fk {
                            logger.Info("Changing field " + new_field.FieldName + " to foreign key: " + new_field_fk_parts[0] + "." + new_field_fk_parts[1])
							// change to fk
							err = this.objects_queries.ChangeFieldToForeignKey(ctx, objects.ChangeFieldToForeignKeyParams{
								FieldType: new_field.FieldType,
								FkTableName: sql.NullString{
									String: new_field_fk_parts[0],
								},
								FkFieldName: sql.NullString{
									String: new_field_fk_parts[1],
								},
								ID: old_field.ID,
							})
						}

						if is_old_field_fk && !is_new_field_fk {
							// change from fk
							err = this.objects_queries.ChangeFieldToNotBeForeignKey(ctx, old_field.ID)
						}

						if err != nil {
							err = errors.New("Cant change field type: " + err.Error())
							logger.Error(err.Error())
							return
						}
					})
				}

				if new_field.FieldOptions.String != old_field.FieldOptions.String && new_field.FieldOptions.Valid {
					// checking if added not null and if so if also add deafult
					has_not_null := strings.Contains(new_field.FieldOptions.String, "NOT NULL") || strings.Contains(new_field.FieldOptions.String, "NOT NULL")
					if has_not_null {
						has_default := strings.Contains(new_field.FieldOptions.String, "DEFAULT") || strings.Contains(new_field.FieldOptions.String, "default")
						if !has_default {
							err = errors.New("Field options changed but not added default value")
							logger.Error(err.Error())
							return err
						}
					}
					// everything is good
					needs_to_create_new = true
					needs_to_copy = true

					if_succeds = append(if_succeds, func() {
                            err := this.objects_queries.ChangeFieldOptions(ctx, objects.ChangeFieldOptionsParams{
                                FieldOptions: sql.NullString(new_field.FieldOptions),
                                ID:           old_field.ID,
                            })

                            if err != nil {
                                err = errors.New("Cant change field options: " + err.Error())
                                logger.Error(err.Error())
                                return
                            }
						})
				}
				break
			}
		}

		if !field_exists {
			// there is a new field
			needs_to_create_new = true
			needs_to_copy = true

            is_new_field_fk, new_field_fk_parts := analyze_if_fk(new_field.FieldType)

			if_succeds = append(if_succeds, func() {
                params := objects.CreateFieldParams{
					FieldName:    new_field.FieldName,
					FieldType:    new_field.FieldType,
					FieldOptions: sql.NullString(new_field.FieldOptions),
					CollectionID: new_collection.ID,
                    IsForeignKey: is_new_field_fk,
				}
                if is_new_field_fk {
                    params.FkTableName = sql.NullString{
                        String: new_field_fk_parts[0],
                        Valid: true,
                    }

                    params.FkFieldName = sql.NullString{
                        String: new_field_fk_parts[1],
                        Valid: true,
                    }
                }
				err := this.objects_queries.CreateField(ctx, params)

				if err != nil {
					err = errors.New("Cant create field: " + err.Error())
					logger.Error(err.Error())
					return
				}
			})
		}
	}

	// fields were removed
	for _, old_field := range old_collection.Fields {
		fields_exists := false
		for _, new_field := range new_collection.Fields {
			if old_field.ID == new_field.ID {
				fields_exists = true
			}
		}

		if !fields_exists {
			needs_to_create_new = true
			needs_to_copy = true

			if_succeds = append(if_succeds, func() {
				err := this.objects_queries.DeleteField(ctx, old_field.ID)

				if err != nil {
					err = errors.New("Cant delete field: " + err.Error())
					logger.Error(err.Error())
					return
				}
			})
		}
	}

	if needs_to_create_new {
		// changing table name
		changed_collection_name := "old_" + new_collection.Name
		err = this.changeCollectionTableName(ctx, new_collection.Name, changed_collection_name)
		if err != nil {
			err = errors.New("Cant change table name: " + err.Error())
			logger.Error(err.Error())
			return err
		}

		// creating new table
		create_query := getAddTableQueryForCollection(ctx, new_collection)
		logger.Info("new table create query: " + create_query)
		_, err = this.db.Exec(create_query)
		if err != nil {
			err = errors.New("Cant create new table: " + err.Error())
			logger.Info("new table create query: " + create_query)
			logger.Error(err.Error())
			// reverting table name change
			this.changeCollectionTableName(ctx, changed_collection_name, new_collection.Name)
			// returning error
			return err
		}

		// copying table
		if needs_to_copy && len(old_collection.Fields) > 0 {
			fields_string := strings.Join(new_fields_strings, ", ")
			copy_query.WriteString(fields_string)
			copy_query.WriteString(") ")
			copy_query.WriteString("SELECT ")
			copy_query.WriteString(fields_string)
			copy_query.WriteString(" FROM ")
			copy_query.WriteString(changed_collection_name)
			copy_query.WriteString(";")

			logger.Info("copy query: " + copy_query.String())

			_, err = this.db.Exec(copy_query.String())
			if err != nil {
				err = errors.New("Cant copy collection: " + err.Error())
				logger.Info("copy query: " + copy_query.String())
				logger.Error(err.Error())
				// droping new table
				this.dropTable(ctx, new_collection.Name)
				// reverting table name change
				this.changeCollectionTableName(ctx, changed_collection_name, new_collection.Name)
				// returning error
				return err
			}

			// changing Collection and table_fields
		}

		// droping old table
		err = this.dropTable(ctx, changed_collection_name)
		if err != nil {
			err = errors.New("Cant drop old table: " + err.Error())
			logger.Error(err.Error())
			return err
		}

		for _, f := range if_succeds {
			f()
		}
	}

	logger.Info("Change table succeds")

	return nil
}

func (this *SqliteDB) changeCollectionTableName(ctx context.Context, old_name string, new_name string) error {
	_, err := this.db.Exec("ALTER TABLE " + old_name + " RENAME TO " + new_name + ";")
	return err
}

func (this *SqliteDB) ChangeCollectionName(ctx context.Context, old_name string, new_name string) error {
	// removing table
	logger := this.logger.With("old_name", old_name, "new_name", new_name)
	err := this.changeCollectionTableName(ctx, old_name, new_name)

	if err != nil {
		logger.Error("Cant change collection name: " + err.Error())
		return err
	}

	err = this.objects_queries.ChangeTableName(ctx, objects.ChangeTableNameParams{
		NewName: new_name,
		OldName: old_name,
	})
	if err != nil {
		logger.Error("Cant change collection name: " + err.Error())
		return err
	}

	return err
}

func (this *SqliteDB) ChangeFieldName(ctx context.Context, collection_name string, old_name string, new_name string, field_id int64) error {
	logger := this.logger.With("collection_name", collection_name, "old_name", old_name, "new_name", new_name)
	// removing table
	_, err := this.db.Exec("ALTER TABLE " + collection_name + " RENAME COLUMN " + old_name + " TO " + new_name + ";")
	if err != nil {
		logger.Error("Cant change field name: " + err.Error())
		return err
	}

	err = this.objects_queries.ChangeFieldName(ctx, objects.ChangeFieldNameParams{
		FieldName: new_name,
		ID:        field_id,
	})
	if err != nil {
		logger.Error("Cant change field name: " + err.Error())
		return err
	}

	return err
}
