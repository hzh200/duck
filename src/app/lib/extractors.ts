interface Extractor {
    name: string
}

const extractors: Array<Extractor> = [{
    name: 'http'
}, {
    name: 'bilibili'
}, {
    name: 'youtube'
}];

export { extractors, Extractor };
