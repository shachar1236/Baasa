import { CompletionContext } from "@codemirror/autocomplete";

const sqlKeywords = [
    "SELECT",
    "FROM",
    "WHERE",
    "AND",
    "OR",
    "ORDER BY",
    "GROUP BY",
    "INSERT INTO",
    "VALUES",
    "UPDATE",
    "SET",
    "DELETE FROM",
    "CREATE TABLE",
    "ALTER TABLE",
    "DROP TABLE",
    "JOIN",
    "INNER JOIN",
    "LEFT JOIN",
    "RIGHT JOIN",
    "ON",
    "AS",
    "DISTINCT",
    "COUNT",
    "SUM",
    "MAX",
    "MIN",
    "AVG",
];

export function createSqlAutocomplete(collections) {
    // Example table and column names (replace with your actual schema)
    // const tables = ["users", "orders", "products"];
    // const columns = {
    //     users: ["id", "name", "email"],
    //     orders: ["id", "user_id", "amount"],
    //     products: ["id", "name", "price"],
    // };

    return function sqlAutocomplete(context) {
        const word = context.matchBefore(/\w*/);
        const text = context.state.doc.toString();
        if (word.from === word.to && !context.explicit) return null;

        let options = [];
        options = sqlKeywords
            .filter((keyword) => keyword.toUpperCase())
            .map((keyword) => ({
                label: keyword,
                type: "keyword",
            }));

        console.log("collections: ", collections);
        console.log("word: ", word);

        const words = text.split(" ");

        let lastWord = null;
        if (words.length > 1) {
            lastWord = words[words.length - 2];
        }

        console.log("lastWord: ", lastWord);

        if (lastWord != null) {
            if (lastWord.toUpperCase() === "SELECT") {
                collections.forEach((collection) => {
                    options.push({
                        label: "id",
                        type: "column",
                    });
                    collection.Fields.forEach((field) => {
                        options.push({
                            label: field.FieldName,
                            type: "column",
                        });
                    });
                });
            } else if (lastWord.toUpperCase() === "FROM") {
                collections.forEach((collection) => {
                    options.push({
                        label: collection.Name,
                        type: "table",
                    });
                });
            }
        }

        return {
            from: word.from,
            options,
        };
    };
}
