package sqlite

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

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
			field := types.TableField{
				ID:           coll.FieldID.Int64,
				CollectionID: coll.CollectionID,
				FieldName:    coll.FieldName.String,
				FieldType:    coll.FieldType.String,
				FieldOptions: types.NullString(coll.FieldOptions),
			}
			collection.Fields = append(collection.Fields, field)
		}
	}
	return collection, err
}

func (this *SqliteDB) GetCollectionById(ctx context.Context, id int64) (types.Collection, error) {
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
			field := types.TableField{
				ID:           coll.FieldID.Int64,
				CollectionID: coll.CollectionID,
				FieldName:    coll.FieldName.String,
				FieldType:    coll.FieldType.String,
				FieldOptions: types.NullString(coll.FieldOptions),
			}
			collection.Fields = append(collection.Fields, field)
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

			curr_collection.Fields = append(curr_collection.Fields, types.TableField{
				ID:           row.FieldID.Int64,
				FieldName:    row.FieldName.String,
				FieldType:    row.FieldType.String,
				FieldOptions: types.NullString(row.FieldOptions),
			})
		}

        collections = append(collections, curr_collection)
	}
	return collections, err
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
    res, err := this.objects_queries.CreateTable(ctx, name)

    if err != nil {
        this.logger.Error("Cant add collection: " + err.Error())
        return collection, err
    }

    collection.ID = res.ID
    collection.Name = res.TableName
    collection.QueryRulesDirectoryPath = res.QueryRulesDirectoryPath
    return collection, err
}

func (this *SqliteDB) DeleteCollection(ctx context.Context, name string) error {
    // deleting table
    _, err := this.db.Exec("DROP TABLE " + name + ";")
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
