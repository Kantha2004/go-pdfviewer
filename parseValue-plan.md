# ROADMAP: Implementing ParseValue() for a PDF Parser

This roadmap describes, in exact implementation order, how to build the `ParseValue()` layer that sits on top of your PDF lexer.

`ParseValue` is the bridge between:
- Tokens (lexer output)
- Structured PDF values (parser output)

This is the MOST IMPORTANT parsing phase after the lexer.

---

## 0. DEFINE THE SCOPE (FREEZE IT)

`ParseValue()` MUST handle ONLY the following:

- Numbers            -> `PDFNumber`
- Names              -> `PDFName`
- Literal strings    -> `PDFString`
- Hex strings        -> `PDFHexString` (or `PDFString` later)
- Arrays             -> `PDFArray`
- Dictionaries       -> `PDFDict`
- Booleans           -> `true` / `false`
- Null               -> `null`
- Indirect references -> `<obj> <gen> R`

`ParseValue()` MUST NOT:

- Parse "obj ... endobj"
- Parse streams
- Resolve references
- Interpret operators
- Parse xref tables

If this scope is violated, the parser will collapse later.

---

## 1. DEFINE PARSER STATE (ABSOLUTELY REQUIRED)

You MUST implement single-token lookahead.

**Reason:**
- Indirect references ("12 0 R") require speculative reads.

**Structure:**

```go
type Parser struct {
    lexer  *Lexer
    buffer *Token // nil or holding one unread token
}
```

**Required helpers:**

- `next()`   -> returns next token (buffer-aware)
- `unread()` -> pushes a token back into buffer

Without this, `ParseValue` cannot be correct.

---

## 2. IMPLEMENT ParseValue() AS A DISPATCHER

`ParseValue()` must:
- Read exactly ONE value
- Delegate based on token type
- NEVER loop endlessly

**High-level logic:**

```go
tok = next()

switch tok.Type {
case TokNumber:
    return parseNumberOrReference()
case TokName:
    return PDFName
case TokString:
    return PDFString
case TokHexString:
    return PDFHexString
case TokArrayStart:
    return parseArray()
case TokDictStart:
    return parseDict()
case TokKeyword:
    return parseKeyword()
default:
    return error
}
```

**Important:**
- `ParseValue` must NEVER consume tokens belonging to the next value.

---

## 3. IMPLEMENT INDIRECT REFERENCE DETECTION (CRITICAL)

**Pattern:**
```
<number> <number> R
```

**Example:**
```
12 0 R
```

**Implementation strategy:**
- First number already read
- Peek second token
- Peek third token

If pattern matches:
-> return `PDFIndirectRef`

If pattern fails:
-> unread tokens in reverse order
-> return `PDFNumber`

**Rules:**
- Only `TokNumber` `TokNumber` `TokKeyword("R")` qualifies
- Anything else must be rolled back

This is the most common source of bugs in PDF parsers.

---

## 4. IMPLEMENT ARRAY PARSING

**Syntax:**
```
[ value value value ]
```

**Implementation:**
- Create empty `PDFArray`
- Loop:
    - Read next token
    - If `TokArrayEnd` -> break
    - unread token
    - `ParseValue()`
    - append to array

**Rules:**
- Arrays are heterogeneous
- Arrays can be nested
- Whitespace is irrelevant

**Exit condition:**
- `TokArrayEnd` ONLY

---

## 5. IMPLEMENT DICTIONARY PARSING

**Syntax:**
```
<< /Key value /Key value >>
```

**Rules:**
- Keys MUST be `TokName`
- Values are parsed via `ParseValue()`
- Dictionary ends ONLY on `TokDictEnd`

**Implementation:**
- Create empty `map[string]PDFValue`
- Loop:
    - Read key token
    - If `TokDictEnd` -> break
    - If not `TokName` -> error
    - `ParseValue()` for value
    - Store in map

Dictionaries can be nested.

---

## 6. IMPLEMENT KEYWORD HANDLING

**Keywords to handle here:**

- `true`  -> `PDFBoolean(true)`
- `false` -> `PDFBoolean(false)`
- `null`  -> `PDFNull`

All other keywords are NOT values.

**Examples of NON-value keywords:**
- `obj`
- `endobj`
- `stream`
- `BT`
- `Tf`
- `Tj`

These must error or be handled at higher layers.

---

## 7. ERROR HANDLING STRATEGY

`ParseValue()` must:
- Fail fast on impossible syntax
- Preserve token stream integrity
- Never swallow tokens silently

**Typical errors:**
- Dictionary key not a name
- Unexpected EOF
- Unterminated array or dictionary

**Error messages should include:**
- Token type
- Token value

---

## 8. VALIDATION CHECKLIST (MUST PASS)

Before moving forward, `ParseValue` must correctly parse:

```
/Type /Page
[0 0 300 300]
<< /Parent 2 0 R >>
<< /Kids [3 0 R] >>
true
false
null
```

**Nested example:**
```
<< /A [1 << /B 2 >>] >>
```

If ANY of these fail, do NOT proceed.

---

## 9. WHAT ParseValue ENABLES

Once implemented, you can:
- Parse page dictionaries
- Parse resource dictionaries
- Parse content stream dictionaries

You CANNOT yet:
- Parse "obj ... endobj"
- Resolve references
- Render anything

That comes next.

---

## 10. NEXT STEP (STRICT ORDER)

After `ParseValue` is COMPLETE and VERIFIED:

-> Implement object parsing:
   `"<obj> <gen> obj ... endobj"`

-> Then xref parsing

**DO NOT SKIP ORDER.**

---

> [!IMPORTANT]
> **PROFESSIONAL NOTE:**
> If `ParseValue` is correct, everything else becomes mechanical.
> If `ParseValue` is wrong, NOTHING else will work.
