# NEXT STEP ROADMAP — XREF PARSING (WITH CONCEPTUAL EXPLANATION)

You have successfully implemented:
- [x] Lexer
- [x] Value parser
- [x] Object parser
- [x] Object table

You can now parse and store ALL indirect objects.

The IMMEDIATE next step is XREF parsing.
BUT: XRef is conceptually different from objects, so understanding comes first.

## 1. WHAT IS XREF? (CONCEPTUAL UNDERSTANDING)

**XREF = Cross-Reference Table**

**Purpose:**
- Maps object numbers → byte offsets in the file
- Allows random access to objects
- Enables incremental updates

**IMPORTANT:**
You do NOT need xref to *parse* objects sequentially, but you DO need it to:
- Validate correctness
- Resolve references reliably
- Parse real-world PDFs

### Example

```
xref
0 5
0000000000 65535 f
0000000010 00000 n
0000000079 00000 n
0000000178 00000 n
0000000314 00000 n
```

**Meaning:**
- object 0 → free
- object 1 → offset 10
- object 2 → offset 79
...

## 2. IMPORTANT RULE (FREEZE THIS)

**XREF IS NOT VALUE GRAMMAR.**
**XREF IS NOT OBJECT GRAMMAR.**

**XREF is a DOCUMENT-LEVEL STRUCTURE.**

Therefore:
- Do NOT use `Parse()`
- Do NOT use `ParseObject()`
- Use raw token / line-based parsing

**This is CRITICAL.**

## 3. WHY XREF IS PARSED AFTER OBJECTS (EVEN THOUGH IT APPEARS LATER)

In a real PDF:

```
[objects...]
xref
trailer
```

BUT:
- Objects can appear in any order
- PDFs can be incrementally updated
- XRef is authoritative

**Strategy:**
1. Parse objects sequentially (DONE)
2. Parse xref to validate / locate them
3. Use xref for reference resolution

## 4. XREF GRAMMAR (STRICT)

**Grammar:**

```
xref
<first-object-number> <count>
<offset> <generation> <n|f>
...
```

**Notes:**
- Offsets are 10-digit, zero-padded
- Generation is 5-digit
- 'n' = in-use
- 'f' = free

**Example line:**
`0000000010 00000 n`

## 5. DATA STRUCTURE FOR XREF (DESIGN NOW)

**Minimal correct structure:**

```go
xrefTable = map[int]XRefEntry
```

Where:

```go
type XRefEntry struct {
    Offset     int
    Generation int
    InUse      bool
}
```

Object number is the key.

## 6. DOCUMENT-LEVEL PARSER RESPONSIBILITY (NEW LAYER)

You now need a **DOCUMENT** parser that:
- Calls `ParseObject()` in a loop
- Detects keyword "xref"
- Switches to XRef parsing mode
- Then parses trailer

**IMPORTANT:**
The parser now becomes STATEFUL:

**State:**
- `ReadingObjects`
- `ReadingXRef`
- `ReadingTrailer`

## 7. CONTROL FLOW FOR XREF PARSING

**Pseudocode:**

```python
while true:
    tok = next token

    if tok == "xref":
        parseXRef()
        break

    else:
        unread tok
        ParseObject()
```

## 8. HOW parseXRef WORKS (STEP BY STEP)

**Step 1:**
Expect keyword "xref"

**Step 2:**
Read two numbers:
- `startObject`
- `count`

**Step 3:**
For i in range(count):
- read offset
- read generation
- read status (n/f)

**Step 4:**
Store in `xrefTable`:
`xrefTable[startObject + i] = entry`

**Step 5:**
Repeat if another subsection appears

## 9. WHAT YOU MUST NOT DO YET

**DO NOT:**
- Seek to offsets
- Re-parse objects via xref
- Resolve references

**First goal:**
- Correctly READ xref
- Store it in memory

## 10. VALIDATION CHECKLIST (MUST PASS)

After xref parsing:
- `xrefTable` size == number of entries
- object 0 must be free
- offsets must be non-negative

Optional (later):
- Validate offsets against parsed object positions

## 11. WHY THIS STEP MATTERS

**Without xref:**
- `/Root` cannot be trusted
- Incremental updates break
- Reference resolution is unreliable

**With xref:**
- Document structure becomes stable
- Trailer becomes meaningful
- Page tree can be resolved

## 12. NEXT STEPS AFTER XREF (DO NOT SKIP)

**Correct order:**
1. XRef parsing (YOU ARE HERE)
2. Trailer parsing
3. Root catalog resolution
4. Page tree traversal

---
**END OF ROADMAP**

# STRUCTURE FOR XREF PARSING (ARCHITECTURE + FILE LAYOUT)

This describes the **STRUCTURE**, not the code. It answers:
- What components you need
- Where xref parsing lives
- How it connects to existing parser layers

You should be able to implement xref cleanly after this.

## 1. GRAMMAR POSITION OF XREF

XRef lives at the **DOCUMENT GRAMMAR** level.

Current layers you have:
- **Lexer** → tokens
- **Value Parser** → `Parse()`
- **Object Parser** → `ParseObject()`
- **Object Table** → storage

**XRef sits ABOVE ALL OF THESE.**

It must:
- **NOT** call `Parse()`
- **NOT** call `ParseObject()`
- Read tokens / raw values directly

## 2. FILE / PACKAGE STRUCTURE (RECOMMENDED)

```text
internal/
├── parser/
│   ├── lexer.go              # DONE
│   ├── value_parser.go       # DONE (Parse)
│   ├── object_parser.go      # DONE (ParseObject)
│   ├── document_parser.go    # NEW (controls flow)
│   ├── xref_parser.go        # NEW (xref logic)
│   └── trailer_parser.go     # LATER
│
└── model/
    ├── object.go             # PDFObject
    ├── values.go             # PDFValue types
    ├── xref.go               # XRefEntry, XRefTable
    └── trailer.go            # Trailer model
```

## 3. CORE DATA STRUCTURES FOR XREF

### XRefEntry
Represents ONE object reference entry.

**Fields:**
- **ObjectNumber**: (implicit from table key)
- **Offset**: (byte offset in file)
- **Generation**: (generation number)
- **InUse**: (true if 'n', false if 'f')

### XRefTable
Stores ALL xref entries.

**Structure:**
```go
xrefTable = map[int]XRefEntry
```

**Key:** object number

## 4. DOCUMENT PARSER STATE MACHINE (VERY IMPORTANT)

The document parser controls WHAT grammar is active.

**States:**
- `STATE_OBJECTS`: parsing indirect objects
- `STATE_XREF`: parsing xref table
- `STATE_TRAILER`: parsing trailer
- `STATE_DONE`: EOF reached

Only ONE state is active at a time.

## 5. CONTROL FLOW FOR DOCUMENT PARSING

**Pseudocode structure (conceptual):**

```python
state = STATE_OBJECTS

while state != STATE_DONE:
    if state == STATE_OBJECTS:
        tok = peek token
        if tok == "xref":
            state = STATE_XREF
        else:
            obj = ParseObject()
            ObjectTable.Add(obj)

    elif state == STATE_XREF:
        ParseXRef()
        state = STATE_TRAILER

    elif state == STATE_TRAILER:
        ParseTrailer()
        state = STATE_DONE
```

## 6. STRUCTURE OF XRef PARSER (RESPONSIBILITIES)

XRef parser must:
1. Expect keyword `xref`
2. Read one or more **SUBSECTIONS**
3. For each subsection:
    - Read `startObjectNumber`
    - Read `count`
    - Read exactly `<count>` entries
4. Populate `XRefTable`

It must **STOP** reading at keyword `trailer`.

## 7. XREF SUBSECTION STRUCTURE

**Subsection grammar:**
```text
<startObj> <count>
<offset> <generation> <n|f>
<offset> <generation> <n|f>
...
```

**Mapping rule:**
`objectNumber = startObj + index`

## 8. TOKEN HANDLING STRATEGY FOR XREF

XRef parsing should:
- Read numbers directly from tokens
- Read status as keyword or raw byte
- **NOT** use recursive descent

XRef lines are **FLAT, FIXED FORMAT**.

## 9. ERROR HANDLING STRATEGY

XRef parser must detect:
- Missing `xref` keyword
- Invalid subsection header
- Wrong number of entries
- Invalid `n/f` marker

Error early. XRef corruption is fatal.

## 10. WHAT XREF PARSER MUST NOT DO

**DO NOT:**
- Seek in file
- Re-parse objects
- Resolve references
- Validate offsets yet

XRef parsing is **PURELY declarative** at this stage.

## 11. OUTPUT OF XREF PARSING

After `ParseXRef()`:
- `XRefTable` is fully populated
- `ObjectTable` remains unchanged
- Trailer parsing begins next

## 12. INTEGRATION POINTS

`XRefTable` will later be used by:
- Reference resolver
- Incremental update handler
- Object reloading (optional)

---
**END OF XREF STRUCTURE**
