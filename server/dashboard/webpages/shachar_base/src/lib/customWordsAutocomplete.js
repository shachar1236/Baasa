import { CompletionContext } from "@codemirror/autocomplete";

export function createCustomWordsAutocomplete() {
    console.log("here1");
    let ck = ["SELECT", "FROM", "WHERE", "AND", "OR", "ORDER BY", "GROUP BY"];
    return function customWordsAutocomplete(context) {
        console.log("here2");
        const word = context.matchBefore(/\w*/);
        if (word.from === word.to && !context.explicit) return null;
        console.log("here3");

        let options = [];

        ck.forEach((word) =>
            options.push({
                label: word,
                type: "keyword",
            })
        );
        console.log("here4");

        return {
            from: word.from,
            options: options,
        };
    };
}
