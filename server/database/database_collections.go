package database

import (
	"context"

	"github.com/shachar1236/Baasa/database/objects"
)

type Collection struct {
    ID int64
    Name string
    QueryRulesDirectoryPath string
    Fields []objects.TableField
}

func GetTableByName(ctx context.Context, name string) (objects.GetTableAndFielsByTableNameRow, error) {
    res, err := objects_queries.GetTableAndFielsByTableName(ctx, name)
    return res, err
}


func GetCollections(ctx context.Context) ([]Collection, error) {
    res, err := objects_queries.GetAllTablesAndFields(ctx)

    collections := []Collection{}
    if len(res) > 0 {

        curr_collection := Collection{
            ID: res[0].CollectionID,
            Name: res[0].FieldName, 
            QueryRulesDirectoryPath: res[0].Queryrulesdirectorypath,
        }

        for _, row := range res {
            if row.CollectionID != curr_collection.ID {
                collections = append(collections, curr_collection)

                curr_collection = Collection{
                    ID: row.CollectionID,
                    Name: row.FieldName, 
                    QueryRulesDirectoryPath: row.Queryrulesdirectorypath,
                }
            } 

            curr_collection.Fields = append(curr_collection.Fields, objects.TableField{
                ID: row.FieldID,
                FieldName: row.FieldName,
                FieldType: row.FieldType,
                FieldOptions: row.FieldOptions,
            })
        }
    }
    return collections, err
}
