# PDF PARSER – COMPLETE IMPLEMENTATION PLAN (END-TO-END)

This document describes the FULL conceptual and implementation roadmap for building a PDF parser from scratch (like a real PDF engine), using clean grammar layering and correct architecture.

This is NOT code. This is the PLAN you follow step-by-step.

## PHASE 0 — CORE PRINCIPLES (READ ONCE, NEVER FORGET)

1. A PDF is NOT one grammar — it is MULTIPLE grammars stacked together.
2. Each grammar layer must be STRICTLY isolated.
3. Each layer consumes tokens from the SAME stream.
4. No layer "passes data down" — layers READ from the stream.

**Grammar layers (bottom → top):**
1. Lexer
2. Value Grammar
3. Object Grammar
4. Document Grammar
5. Cross-Reference Grammar

**Violating this separation guarantees bugs.**

## PHASE 1 — LEXER (DONE)

**Responsibility:**
Convert raw bytes → token stream

**Tokens include:**
- Numbers
- Names
- Strings
- Hex strings
- Keywords
- Array start/end
- Dictionary start/end
- EOF

**Lexer rules:**
- Skip whitespace and comments
- Do NOT interpret meaning
- Do NOT group tokens

**Output:**
Stream of tokens, nothing more.

**Status:**
- [x] COMPLETED

## PHASE 2 — VALUE PARSER (DONE)

**Responsibility:**
Parse ONE PDF VALUE from token stream.

**Supported values:**
- Number (int / float)
- Boolean (true / false)
- Null
- Name
- String
- Hex string
- Array
- Dictionary
- Indirect Reference (n n R)

**Critical rule:**
Indirect references are VALUE grammar, not object grammar.

**Stop condition:**
The value parser MUST stop exactly at the end of ONE value.

**Status:**
- [x] COMPLETED

## PHASE 3 — OBJECT PARSER (DONE)

**Responsibility:**
Parse ONE indirect object.

**Grammar:**
```
<object-number> <generation> obj
    <value>
endobj
```

**Rules:**
- object-number and generation MUST be integers
- "obj" and "endobj" are STRUCTURAL keywords
- Exactly ONE value per object

**Critical rule:**
Object parser reads the value itself; it is NOT passed in.

**Status:**
- [x] COMPLETED

## PHASE 4 — OBJECT TABLE (DONE)

**Responsibility:**
Store all parsed objects for later lookup.

**Data structure:**
`map[int]map[int]*PDFObject`

**Meaning:**
`objectTable[objNum][gen] → object`

**Rules:**
- Duplicate objects overwrite (last wins)
- Missing object is a normal condition

**Status:**
- [x] COMPLETED

## PHASE 5 — DOCUMENT-LEVEL PARSER (CURRENT PHASE)

**Responsibility:**
Control the HIGH-LEVEL parsing flow.

**This layer:**
- Repeatedly parses objects
- Detects transition to XRef section
- Switches grammar modes

**This is where the parser becomes STATEFUL.**

**States:**
- ReadingObjects
- ReadingXRef
- ReadingTrailer

**Status:**
- [ ] **IN PROGRESS**

## PHASE 6 — XREF PARSING (NEXT IMPLEMENTATION STEP)

**Responsibility:**
Parse the Cross-Reference Table.

**Concept:**
XRef is the AUTHORITATIVE index of object locations.

**Grammar:**
```
xref
<firstObj> <count>
<offset> <generation> <n|f>
...
```

**Data structure:**
`map[int]XRefEntry`

where:
```go
XRefEntry {
    Offset     int
    Generation int
    InUse      bool
}
```

**Rules:**
- Object 0 is always free
- Offsets are byte positions
- XRef does NOT use value grammar

**Do NOT:**
- Seek to offsets yet
- Resolve references yet

**Status:**
- [ ] NOT STARTED

## PHASE 7 — TRAILER PARSING

**Responsibility:**
Parse the trailer dictionary.

**Grammar:**
```
trailer
<< /Size N /Root n n R ... >>
```

**Trailer tells:**
- Total object count
- Root catalog object
- Encryption info
- Info dictionary

**Trailer is a VALUE (dictionary),**
but is introduced by STRUCTURAL keyword "trailer".

**Status:**
- [ ] NOT STARTED

## PHASE 8 — STARTXREF + EOF HANDLING

**Responsibility:**
Locate the correct xref offset.

**Grammar:**
```
startxref
<offset>
%%EOF
```

**Rules:**
- startxref is authoritative
- Last xref wins (incremental updates)

**Status:**
- [ ] NOT STARTED

## PHASE 9 — REFERENCE RESOLUTION

**Responsibility:**
Resolve indirect references (n n R).

**Strategy:**
- Use object table + xref
- Lazy resolution (on demand)

**Rule:**
Never resolve references during parsing.

**Status:**
- [ ] NOT STARTED

## PHASE 10 — ROOT CATALOG & PAGE TREE

**Responsibility:**
Build the logical document structure.

**Steps:**
- Resolve /Root
- Resolve /Pages
- Traverse page tree
- Extract page contents

**Status:**
- [ ] NOT STARTED

## PHASE 11 — CONTENT STREAM INTERPRETER

**Responsibility:**
Interpret page content streams.

**This is a STACK-BASED language:**
- Operands pushed
- Operators consume operands

**Examples:**
- Text drawing
- Graphics operators
- Path construction

**Status:**
- [ ] NOT STARTED

## PHASE 12 — RENDERING (VERY LAST)

**Responsibility:**
Convert PDF instructions → pixels / vectors.

**This is OUTSIDE the parser.**

**Parser must already produce:**
- Page structure
- Content streams
- Resources

**Status:**
- [ ] NOT STARTED

## FINAL RULES (PRINT THIS)

1. **NEVER** mix grammar layers
2. **NEVER** parse ahead without unread support
3. **NEVER** resolve references early
4. **ALWAYS** trust XRef over sequential parsing

If you follow this plan in order, you will end up with a REAL PDF engine.
