# ROADMAP: NEXT IMMEDIATE STEP — INDIRECT OBJECT PARSING

You have completed:
- Lexer
- ParseValue (value grammar)

The error you saw ("unexpected keyword: obj") proves your layering is correct.

The IMMEDIATE next step is to implement the OBJECT GRAMMAR layer.

---

## CONTEXT (VERY IMPORTANT)

PDF has multiple grammars stacked on top of each other:

1. Token grammar      → Lexer
2. Value grammar      → ParseValue (DONE)
3. Object grammar     → ParseObject (NEXT)
4. Cross-reference    → XRef parser
5. Page tree          → Page resolver
6. Content streams    → Interpreter

You are now moving from (2) to (3).

---

## 1. UNDERSTAND THE OBJECT GRAMMAR (FREEZE THIS)

An indirect object in PDF has EXACTLY this form:

```
<object-number> <generation-number> obj
    <value>
endobj
```

Example:

```
1 0 obj
<< /Type /Catalog >>
endobj
```

Rules:
- object-number   : integer
- generation      : integer
- "obj"           : keyword (NOT a value)
- <value>         : exactly ONE PDF value (use ParseValue)
- "endobj"        : keyword

Anything else is a syntax error.

ParseObject must NOT:
- Parse xref
- Parse streams (yet)
- Resolve references

---

## 2. DEFINE THE OBJECT MODEL (IN model PACKAGE)

You need a struct like:

```go
type PDFObject struct {
    Number     int
    Generation int
    Value      PDFValue
}
```

This represents ONE indirect object.

Later, you will store these in a map:

```go
map[int]map[int]*PDFObject
```

(object number → generation → object)

---

## 3. DEFINE ParseObject() API (parser PACKAGE)

Signature:

```go
func (p *Parser) ParseObject() (*model.PDFObject, error)
```

Responsibilities:

1. Read object number token
2. Read generation number token
3. Expect keyword "obj"
4. Call ParseValue() ONCE
5. Expect keyword "endobj"
6. Return PDFObject

ParseObject is a STRICT function.

---

## 4. IMPLEMENT ParseObject() — STEP BY STEP

### Step 4.1 — Read object number

```go
tok := p.next()
// - Must be TokNumber
// - Must parse to integer
```

Error if:
- Not a number
- Not an integer

### Step 4.2 — Read generation number

```go
tok := p.next()
// - Must be TokNumber
// - Must parse to integer
```

### Step 4.3 — Expect "obj"

```go
tok := p.next()
// - Must be TokKeyword
// - Must have Value == "obj"
```

This is NOT optional.

### Step 4.4 — Parse the object VALUE

```go
val, err := p.Parse()
```

This uses your existing ParseValue.

DO NOT loop.
DO NOT parse multiple values.

### Step 4.5 — Expect "endobj"

```go
tok := p.next()
// - Must be TokKeyword
// - Must have Value == "endobj"
```

### Step 4.6 — Construct PDFObject

```go
return &PDFObject{
    Number:     objNum,
    Generation: gen,
    Value:      val,
}
```

---

## 5. HANDLE EOF AND NON-OBJECT TOKENS

PDFs contain non-object sections:
- `xref`
- `trailer`
- `startxref`
- `%%EOF`

ParseObject must:
- Detect EOF cleanly
- Error if syntax is broken

Higher-level document parser will later decide:
- When to stop calling ParseObject
- When to switch to xref parsing

---

## 6. TEMPORARY main.go CHANGE (IMPORTANT)

Right now, `main.go` must NOT call `ParseValue()`.
It must call `ParseObject()` in a loop.

Example (temporary debug):

```go
for {
    obj, err := parser.ParseObject()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Object %d %d: %+v\n", obj.Number, obj.Generation, obj.Value)
}
```

---

## 7. VALIDATION CHECKLIST (MUST PASS)

Using your minimal PDF, you MUST be able to parse:

```
1 0 obj
<< /Type /Catalog >>
endobj

2 0 obj
<< /Type /Pages >>
endobj

3 0 obj
<< /Type /Page >>
endobj
```

Failure on ANY of these means ParseObject is wrong.

---

## 8. WHAT YOU MUST NOT DO YET

DO NOT:
- Parse xref
- Decode streams
- Render pages
- Resolve references

Those come AFTER object parsing is stable.

---

## 9. WHY THIS STEP IS CRITICAL

Once ParseObject works:
- XRef parsing becomes mechanical
- Reference resolution becomes possible
- Page tree traversal becomes trivial

If ParseObject is wrong:
- XRef offsets won't match
- Objects will be corrupted
- Debugging will be extremely hard

---

## 10. NEXT STEPS AFTER THIS (DO NOT JUMP AHEAD)

Correct order AFTER ParseObject:
1. Object table construction
2. XRef table parsing
3. Trailer parsing
4. Root catalog resolution

---
# END OF ROADMAP
