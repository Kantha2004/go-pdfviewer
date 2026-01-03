# PDF Lexer — Implementation Roadmap

## Phase 1 — Define the Contract (Do This First)

### Objective

Freeze **what the lexer is responsible for** and **what it must never do**.

### Tasks

* Define token types (`TokNumber`, `TokName`, etc.)
* Define delimiter set
* Define whitespace set
* Decide error-tolerance level

### Deliverables

* `tokens.go`
* Lexer interface: `NextToken() (Token, error)`

### Exit Criteria

* No semantic logic in lexer
* Token set covers PDF syntax (not meaning)

---

## Phase 2 — Byte Reader & Lookahead

### Objective

Make byte reading **predictable and reversible**.

### Tasks

* Wrap input in `bufio.Reader`
* Implement `readByte()` / `unreadByte()`
* Handle EOF cleanly

### Deliverables

* Reader helpers
* Lookahead capability

### Exit Criteria

* Lexer can read and unread bytes safely
* EOF never causes panic or infinite loop

---

## Phase 3 — Whitespace & Comment Skipping

### Objective

Ensure lexer always starts tokenization at meaningful input.

### Tasks

* Implement `isWhitespace()`
* Implement `%` comment skipping
* Loop until non-whitespace, non-comment byte found

### Deliverables

* `skipWhitespaceAndComments()`

### Exit Criteria

* Comments never appear in token stream
* Arbitrary whitespace does not affect tokens

### Common Pitfalls

* Forgetting `0x00` as whitespace
* Mishandling CR/LF combinations

---

## Phase 4 — Single-Character Tokens

### Objective

Handle trivial tokens first.

### Tasks

* Emit tokens for:

  * `[`
  * `]`
  * `(`
  * `/`

### Deliverables

* Basic `NextToken()` dispatch

### Exit Criteria

* These tokens always emit correctly
* No accidental fall-through

---

## Phase 5 — Compound Delimiters (`<<`, `>>`)

### Objective

Correctly distinguish ambiguous delimiters.

### Tasks

* Implement one-byte lookahead
* Differentiate:

  * `<<` vs `<hex>`
  * `>>` vs invalid `>`

### Deliverables

* Dictionary start/end logic

### Exit Criteria

* Dictionaries never misread as hex strings
* No byte loss during lookahead

### Common Pitfalls

* Forgetting to unread lookahead byte
* Treating `<` as always hex string

---

## Phase 6 — Number Tokenization

### Objective

Capture all valid PDF number formats.

### Tasks

* Detect numeric start (`+ - . digit`)
* Read until delimiter or whitespace
* Do **not** parse to float yet

### Deliverables

* `readNumber()`

### Exit Criteria

* All numeric forms captured verbatim
* Invalid numbers tolerated

---

## Phase 7 — Name Tokenization

### Objective

Correctly parse PDF names.

### Tasks

* Trigger on `/`
* Read until delimiter or whitespace
* Exclude `/` from value

### Deliverables

* `readName()`

### Exit Criteria

* Dictionary keys parse correctly
* No delimiter leakage

---

## Phase 8 — Literal String Tokenization (Critical)

### Objective

Correctly parse nested strings.

### Tasks

* Track parentheses depth
* Allow arbitrary content
* Terminate only when depth returns to zero

### Deliverables

* `readLiteralString()`

### Exit Criteria

* Nested strings work
* Strings can contain line breaks

### Common Pitfalls

* Ignoring nested parentheses
* Assuming strings end at first `)`

---

## Phase 9 — Hex String Tokenization

### Objective

Handle binary strings safely.

### Tasks

* Trigger on `<` (when not dictionary)
* Read until `>`
* Ignore internal whitespace

### Deliverables

* `readHexString()`

### Exit Criteria

* Hex strings captured verbatim
* No premature termination

---

## Phase 10 — Keyword Tokenization

### Objective

Capture all remaining symbols.

### Tasks

* Read until whitespace or delimiter
* Preserve case
* Do not interpret meaning

### Deliverables

* `readKeyword()`

### Exit Criteria

* Keywords like `obj`, `stream`, `R` tokenize correctly

---

## Phase 11 — Error Handling Strategy

### Objective

Make lexer resilient to malformed PDFs.

### Tasks

* Decide recoverable vs fatal errors
* Avoid panics
* Return meaningful errors when needed

### Deliverables

* Stable error behavior

### Exit Criteria

* Broken PDFs do not crash lexer
* Lexer stops cleanly on fatal I/O errors

---

## Phase 12 — Integration & Validation

### Objective

Prove correctness.

### Tasks

* Tokenize `minimal.pdf`
* Tokenize real PDFs
* Print token stream for inspection
* Fuzz input with random bytes

### Deliverables

* Debug CLI
* Test PDFs

### Exit Criteria

* No crashes
* Deterministic output
* Correct token boundaries

---

## Final Completion Criteria

The lexer is considered **complete** when:

* It can tokenize real-world PDFs without crashing
* Token boundaries are correct
* No semantic logic exists in lexer
* Parser can consume tokens without hacks

---

## Professional Advice (Important)

Do **not** optimize early.
Do **not** add semantic checks.
Do **not** skip malformed-PDF tolerance.

A correct lexer is **boring, strict, and forgiving**.

---

### Next Logical Step
> Add extensive tests and then **freeze the lexer**.

