package access_rules

import (
	"errors"
	"log"
	"strings"

	"github.com/shachar1236/Baasa/database"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)


func databaseQuery(L *lua.LState) int {
    query := L.ToString(1)
    args := L.ToTable(2)
    my_args := make(map[string]any)
    args.ForEach(func(l1, l2 lua.LValue) {
        if str, ok := l1.(lua.LString); ok {
            my_args[string(str)] = l2
        }
    })
    res, err := database.RunQuery(query, my_args)
    if err != nil {
        log.Println("Cant run query from lua error: ", err)
    }

    L.Push(luar.New(L, res))
    return 1                     /* number of results */
}

func databaseCount(L *lua.LState) int {
    query_select := "SELECT COUNT(*) FROM "
    query_where := " WHERE "

    table := L.ToString(1)
    count_if := L.ToString(2)
    args := L.ToTable(3)
    
    var sb strings.Builder
    sb.Grow(len(query_select) + len(table) + len(query_where) + len(count_if) + 1)
    sb.WriteString(query_select)
    sb.WriteString(table)
    sb.WriteString(query_where)
    sb.WriteString(count_if)
    sb.WriteString(";")

    var my_args []any
    args.ForEach(func(l1, l2 lua.LValue) {
        my_args = append(my_args, l2)
    })

    count, err := database.RunCountQuery(sb.String(), my_args)
    if err != nil {
        log.Println("Cant run count query from lua: ", err)
    }

    L.Push(lua.LNumber(count))
    return 1
}

func CheckRules(rules_file_path string, filters *string, request Request) (bool, error) {
    L := lua.NewState()
    defer L.Close()

    L.SetGlobal("Request", luar.New(L, request))
    L.SetGlobal("Filters", lua.LString(""))
    L.SetGlobal("Accept", lua.LFalse)

    L.SetGlobal("Query", L.NewFunction(databaseQuery))
    L.SetGlobal("Count", L.NewFunction(databaseCount))

    if err := L.DoFile(rules_file_path); err != nil {
        return false, err
    }

    accept := lua.LVAsBool(L.GetGlobal("Accept"))
    new_filters := L.GetGlobal("Filters")
    if new_filters.Type() == lua.LTString {
        if str, ok := new_filters.(lua.LString); ok {
            *filters = string(str)
        } else {
            return false, errors.New("Filters is not a string")
        }
    } else {
        return false, errors.New("Filters is not a string")
    }


    return accept, nil
}
