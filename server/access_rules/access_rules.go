package access_rules

import (
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/shachar1236/Baasa/database"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type AccessRules struct {
	logger *slog.Logger
	db     database.Database

	databaseQuery lua.LGFunction
	databaseCount lua.LGFunction 
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

	return AccessRules{
		logger: logger,
		db:     db,
        databaseQuery: getDatabaseQueryFunction(db, logger),
        databaseCount: getDatabaseCountFunction(db, logger),
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

func  getDatabaseCountFunction(db database.Database, logger *slog.Logger) lua.LGFunction {
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

func (this *AccessRules) CheckRules(rules_file_path string, filters *string, request Request) (bool, error) {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("Request", luar.New(L, request))
	L.SetGlobal("Filters", lua.LString(""))
	L.SetGlobal("Accept", lua.LFalse)

	L.SetGlobal("Query", L.NewFunction(this.databaseQuery))
	L.SetGlobal("Count", L.NewFunction(this.databaseCount))

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
