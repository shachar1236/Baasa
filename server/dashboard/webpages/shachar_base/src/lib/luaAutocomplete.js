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

        const options = [...luaKeywords, ...custom_keywords, ...variables]
            .filter((item, index, array) => array.indexOf(item) === index) // Remove duplicates
            .filter((keyword) => keyword.startsWith(word.text))
            .map((keyword) => ({ label: keyword, type: "keyword" }));

        return {
            from: word.from,
            options,
        };
    };
}
