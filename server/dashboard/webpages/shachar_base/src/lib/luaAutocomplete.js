import { CompletionContext } from "@codemirror/autocomplete";

const luaKeywords = [
    "and",
    "break",
    "do",
    "else",
    "elseif",
    "end",
    "false",
    "for",
    "function",
    "if",
    "in",
    "local",
    "nil",
    "not",
    "or",
    "repeat",
    "return",
    "then",
    "true",
    "until",
    "while",
    "print",
    "io",
    "math",
    "string",
    "table",
];

export function generateLuaAutocomplete(custom_keywords) {
    return function luaAutocomplete(context) {
        const word = context.matchBefore(/\w*/);
        if (word.from === word.to && !context.explicit) return null;

        let docText = context.state.doc.toString();

        docText = docText.slice(0, docText.length - 1);
        const variables = Array.from(docText.matchAll(/\b\w+\b/g)).map(
            (match) => match[0]
        );

        const options = [
            ...luaKeywords,
            ...Object.keys(custom_keywords),
            ...variables,
        ]
            .filter((item, index, array) => array.indexOf(item) === index) // Remove duplicates
            .filter((keyword) => keyword.startsWith(word.text))
            .map((keyword) => ({ label: keyword, type: "keyword" }));

        const word_parts = word.text.split(".");
        console.log("word_parts: ", word_parts);
        if (word_parts.length > 0) {
            // if custom keywords contains word_parts[0]
            if (word_parts[0] in custom_keywords) {
                let last_complete = word_parts[0];
                word_parts.slice(1).forEach((part) => {
                    if (last_complete.hasOwnProperty(part)) {
                        last_complete = last_complete[part];
                    }
                });

                console.log("last_complete: ", last_complete);

                for (const key in last_complete) {
                    options.push({
                        label: key,
                        type: "keyword",
                    });
                }
            }
        }

        return {
            from: word.from,
            options,
        };
    };
}
