interface Parser {
    name: string
}

const parsers: Array<Parser> = [{
    name: 'http'
}, {
    name: 'bilibili'
}, {
    name: 'youtube'
}];

export { parsers, Parser };
