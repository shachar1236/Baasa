package sqlite

import (
	"context"

	"github.com/shachar1236/Baasa/database/types"
)

func (this *SqliteDB) GetCollectionByName(ctx context.Context, name string) (types.Collection, error) {
    res, err := this.objects_queries.GetAllTablesAndFields(ctx)

    var collection types.Collection
    if err != nil {
        this.logger.Error("Cant get collection: ", err)
        return collection, err
    }

    collection.ID = res[0].CollectionID
    collection.Name = res[0].TableName
    collection.QueryRulesDirectoryPath = res[0].Queryrulesdirectorypath
    for _, coll := range res {
        field := types.TableField{
            ID: coll.FieldID,
            CollectionID: coll.CollectionID,
            FieldName: coll.FieldName,
            FieldType: coll.FieldType,
            FieldOptions: types.NullString(coll.FieldOptions),
        }
        collection.Fields = append(collection.Fields, field)
    }
    return collection, err
}


func (this *SqliteDB) GetCollections(ctx context.Context) ([]types.Collection, error) {
    res, err := this.objects_queries.GetAllTablesAndFields(ctx)

    collections := []types.Collection{}
    if len(res) > 0 {

        curr_collection := types.Collection{
            ID: res[0].CollectionID,
            Name: res[0].FieldName, 
            QueryRulesDirectoryPath: res[0].Queryrulesdirectorypath,
        }

        for _, row := range res {
            if row.CollectionID != curr_collection.ID {
                collections = append(collections, curr_collection)

                curr_collection = types.Collection{
                    ID: row.CollectionID,
                    Name: row.FieldName, 
                    QueryRulesDirectoryPath: row.Queryrulesdirectorypath,
                }
            } 

            curr_collection.Fields = append(curr_collection.Fields, types.TableField{
                ID: row.FieldID,
                FieldName: row.FieldName,
                FieldType: row.FieldType,
                FieldOptions: types.NullString(row.FieldOptions),
            })
        }
    }
    return collections, err
}
