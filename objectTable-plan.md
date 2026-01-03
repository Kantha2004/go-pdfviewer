# ROADMAP: IMMEDIATE NEXT STEP — BUILDING THE OBJECT TABLE

You have successfully completed:
- Lexer
- Value parser (Parse)
- Object parser (ParseObject)

Your parser can now correctly read:
```
<objNum> <genNum> obj
    <value>
endobj
```

The IMMEDIATE next step is NOT xref yet.
The immediate next step is to COLLECT objects into an OBJECT TABLE.

---

## 1. WHAT IS THE OBJECT TABLE?

In a PDF, indirect objects are referenced by:

```
<object-number> <generation-number> R
```

To resolve those references later, you need a lookup structure.

Conceptually:

```
(object number, generation) -> PDFObject
```

This is called the OBJECT TABLE.

---

## 2. WHY YOU NEED THE OBJECT TABLE BEFORE XREF

Even though xref exists, you should FIRST:

- Parse objects sequentially
- Store them in memory

Reasons:
- Simpler debugging
- Easier correctness validation
- You already have sequential parsing working

XRef parsing will later OPTIMIZE object lookup, not replace logic.

---

## 3. DATA STRUCTURE DESIGN (KEEP IT SIMPLE)

Start with a SIMPLE structure:

```go
objectTable = map[int]*PDFObject
```

Assumption (valid for now):
- Generation numbers are almost always 0

Later, you can upgrade to:

```go
map[int]map[int]*PDFObject
```

DO NOT over-engineer now.

---

## 4. CONTROL FLOW CHANGE (IMPORTANT)

Your program flow becomes:

1. Initialize parser
2. Loop: ParseObject()
3. Store each object in objectTable
4. Stop on EOF

NO value parsing in main anymore.

---

## 5. MAIN LOOP LOGIC (CONCEPTUAL)

Pseudocode:

```python
objectTable = {}

while True:
    obj, err = ParseObject()

    if err == EOF:
        break

    if err:
        fatal error

    objectTable[obj.Number] = obj
```

This loop MUST:
- Stop cleanly on EOF
- Fail fast on syntax errors

---

## 6. VALIDATION CHECKPOINT (VERY IMPORTANT)

After this step, you MUST be able to:

- Print total number of objects
- Print each object's number and type

Example output:

```
Object 1: Dictionary
Object 2: Dictionary
Object 3: Dictionary
Object 4: Stream (later)
```

If this does not work, DO NOT CONTINUE.

---

## 7. WHAT YOU MUST NOT DO YET

DO NOT:
- Resolve indirect references
- Parse xref
- Parse trailer
- Touch streams

Those depend on a COMPLETE object table.

---

## 8. WHY THIS STEP IS CRITICAL

Once the object table exists:

- Indirect references (n n R) can be resolved
- Page tree traversal becomes possible
- Trailer parsing has a place to attach

Without this step:
- Everything else becomes fragile
- Debugging becomes extremely difficult

---

## 9. COMMON MISTAKES TO AVOID

- Over-handling generation numbers too early
- Trying to parse xref before objects
- Mixing object parsing with value parsing
- Ignoring EOF conditions

---

## 10. EXIT CRITERIA FOR THIS STEP

You are DONE with this step when:

- [ ] All objects in the PDF are parsed
- [ ] All objects are stored in objectTable
- [ ] No syntax errors on valid PDFs
- [ ] EOF stops parsing cleanly

ONLY THEN move forward.

---
# END OF ROADMAP
