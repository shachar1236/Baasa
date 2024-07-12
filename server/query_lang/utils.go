package querylang

import "github.com/shachar1236/Baasa/database/types"

func is_white_space(c byte) bool {
	return c == ' ' || c == '\n'
}

func hasCollection(collections []types.Collection, collection_name string) bool {
	for _, curr_collection := range collections {
		if curr_collection.Name == collection_name {
			return true
		}
	}
	return false
}

func DoesCollectionHasField(collection types.Collection, field_name string) bool {
	for _, field := range collection.Fields {
		if field.FieldName == field_name {
			return true
		}
	}
	return false
}
