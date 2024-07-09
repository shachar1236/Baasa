package access_rules

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/shachar1236/Baasa/database"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type AccessRules struct {
	logger *slog.Logger
	db     database.Database

    main_lua_state *lua.LState
    main_lua_state_mutex sync.Mutex
}

func New(db database.Database) AccessRules {
	logFile, err := os.OpenFile("logs/access_rules.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	var mw io.Writer
	if err != nil {
		mw = os.Stdout
	} else {
		mw = io.MultiWriter(os.Stdout, logFile)
	}
	logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{AddSource: true}))

    // setting functions
    L := lua.NewState()
	L.SetGlobal("Query", L.NewFunction(getDatabaseQueryFunction(db, logger)))
	L.SetGlobal("Count", L.NewFunction(getDatabaseCountFunction(db, logger)))
	L.SetGlobal("Get", L.NewFunction(getDatabaseGetFunction(db, logger)))

    base_string := `
_name = {
    count = function(filters, args) return Count("_name", filters, args) end,
    get = function(filters, args) return Get("_name", filters, args) end,
}`

    base_collections, err := db.GetBaseCollections()
    for _, collection := range base_collections {
        code := strings.ReplaceAll(base_string, "_name", collection)
        logger.Info("Running base access rules lua code: " + code)
        err = L.DoString(code)
        if err != nil {
            logger.Error("Error in base access rules lua code: " + err.Error())
        }
    }

    collections, err := db.GetCollections(context.Background())
    if err != nil {
        logger.Error("Error in base access rules lua code: " + err.Error())
    } else {
        for _, collection := range collections {
            code := strings.ReplaceAll(base_string, "_name", collection.Name)
            err = L.DoString(code)
            if err != nil {
                logger.Error("Error in base access rules lua code: " + err.Error())
            }
        }
    }

	return AccessRules{
		logger: logger,
		db:     db,
        main_lua_state: L,
	}
}

func getDatabaseQueryFunction(db database.Database, logger *slog.Logger) lua.LGFunction {
	return func(L *lua.LState) int {
		query := L.ToString(1)
		args := L.ToTable(2)
		my_args := make(map[string]any)
		args.ForEach(func(l1, l2 lua.LValue) {
			if str, ok := l1.(lua.LString); ok {
				my_args[string(str)] = l2
			}
		})
		res, err := db.RunQuery(query, my_args)
		if err != nil {
			logger.Info("Cant run query from lua error: ", err)
		}

		L.Push(luar.New(L, res))
		return 1 /* number of results */
	}
}

func getDatabaseCountFunction(db database.Database, logger *slog.Logger) lua.LGFunction {
	return func(L *lua.LState) int {
		collection_name := L.ToString(1)
		filters := L.ToString(2)
		args := L.ToTable(3)

		var my_args []any
		args.ForEach(func(l1, l2 lua.LValue) {
			my_args = append(my_args, l2)
		})

		count, err := db.BasicCount(collection_name, filters, my_args)
		if err != nil {
			logger.Info("Cant run count query from lua: ", err)
		}

		L.Push(lua.LNumber(count))
		return 1
	}
}

func getDatabaseGetFunction(db database.Database, logger *slog.Logger) lua.LGFunction {
	return func(L *lua.LState) int {
		collection_name := L.ToString(1)
		filters := L.ToString(2)
		args := L.ToTable(3)

		var my_args []any
		args.ForEach(func(l1, l2 lua.LValue) {
			my_args = append(my_args, l2)
		})

		res, err := db.Get(collection_name, filters, my_args)
		if err != nil {
			logger.Info("Cant run get query from lua: ", err)
		}

		L.Push(luar.New(L, res))
		return 1
	}
}

func (this *AccessRules) CheckRules(rules_file_path string, filters *string, request Request, query_args map[string]any) (bool, error) {
    this.main_lua_state_mutex.Lock()
    defer this.main_lua_state_mutex.Unlock()

    L := this.main_lua_state

	L.SetGlobal("Request", luar.New(L, request))
	L.SetGlobal("Args", luar.New(L, query_args))
	L.SetGlobal("Filters", lua.LString(""))
	L.SetGlobal("Accept", lua.LFalse)

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
