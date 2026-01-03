# PDF Parser & Renderer — Complete Roadmap

## Phase 0 — Foundations (Mandatory)

### Objectives

* Establish project structure
* Lock architectural boundaries
* Prevent scope creep

### Deliverables

* Repository layout finalized
* Core model types defined
* Build pipeline working

### Tasks

* Define `PDFValue` model
* Define token types
* Create CLI entry point
* Add minimal test PDF

### Exit Criteria

* Project builds
* No rendering logic exists yet

---

## Phase 1 — Lexer (Tokenizer)

### Objectives

* Convert byte stream → token stream
* Handle malformed PDFs gracefully

### Tasks

* Byte-level reader
* Whitespace handling
* Comment handling
* Delimiter recognition
* Token generation
* Lookahead support

### Deliverables

* Fully tested lexer
* Token stream debug tool

### Exit Criteria

* Real PDFs tokenize without crash
* Token stream matches spec

### Common Pitfalls

* Line-based parsing
* Missing delimiters
* Incorrect comment termination

---

## Phase 2 — Value Parser (Syntax Layer)

### Objectives

* Convert tokens → structured values

### Tasks

* Recursive value parsing
* Arrays and dictionaries
* Boolean, null handling
* Indirect reference detection
* Single-token lookahead buffer

### Deliverables

* `ParseValue()` implementation
* Value tree printer

### Exit Criteria

* Page dictionaries parse correctly
* Indirect references resolved syntactically

### Common Pitfalls

* Losing tokens during lookahead
* Confusing numbers with references
* Incorrect dictionary termination

---

## Phase 3 — Object & XRef Parsing (Document Backbone)

### Objectives

* Build object table
* Enable random access

### Tasks

* Parse `obj … endobj`
* Parse classic xref table
* Parse `startxref`
* Build object offset map
* Load trailer dictionary

### Deliverables

* Object store
* XRef table
* Catalog reference

### Exit Criteria

* Any object can be resolved by number
* Page tree root identified

### Common Pitfalls

* Ignoring generation numbers
* Incorrect xref offsets
* Not handling incremental updates

---

## Phase 4 — Page Tree Resolution

### Objectives

* Resolve logical pages in order

### Tasks

* Traverse `/Pages` tree
* Inherit resources
* Resolve `/MediaBox`, `/Rotate`
* Collect `/Contents`

### Deliverables

* Ordered page list
* Page metadata

### Exit Criteria

* Correct page count
* Correct media boxes

### Common Pitfalls

* Ignoring resource inheritance
* Mishandling `/Kids`
* Incorrect recursion depth

---

## Phase 5 — Stream Decoding

### Objectives

* Extract raw content streams

### Tasks

* Stream length resolution
* `/FlateDecode` support
* Handle multiple content streams
* Concatenate streams

### Deliverables

* Decoded content stream bytes

### Exit Criteria

* Content streams readable as plain text
* No data corruption

### Common Pitfalls

* Incorrect stream length
* Ignoring indirect `/Length`
* Forgetting multiple streams

---

## Phase 6 — Content Stream Interpreter (Graphics Engine)

### Objectives

* Execute PDF drawing commands

### Tasks

* Operand stack
* Operator dispatch
* Graphics state stack
* Path construction
* Color and line width

### Deliverables

* Interpreter loop
* Operator handlers

### Exit Criteria

* Lines and shapes render correctly
* Graphics state push/pop works

### Common Pitfalls

* Incorrect operand order
* Missing state isolation
* Ignoring CTM

---

## Phase 7 — Rasterizer (Rendering Core)

### Objectives

* Convert vector paths → pixels

### Tasks

* Path flattening
* Scanline fill
* Stroke rendering
* Coordinate transforms
* Anti-aliasing (optional)

### Deliverables

* RGBA image output
* PNG export

### Exit Criteria

* Shapes render identically across runs

### Common Pitfalls

* Floating-point drift
* Incorrect fill rules
* Coordinate inversion errors

---

## Phase 8 — Text Rendering (Hardest Phase)

### Objectives

* Render glyphs accurately

### Tasks

* Text state tracking
* Type1 font support
* Encoding maps
* Glyph positioning
* Text matrix handling

### Deliverables

* Rendered text
* Text bounding boxes

### Exit Criteria

* Simple PDFs show readable text

### Common Pitfalls

* Treating text as strings
* Ignoring glyph metrics
* Incorrect spacing rules

---

## Phase 9 — Image Rendering

### Objectives

* Render embedded images

### Tasks

* Image XObject parsing
* Color space handling
* Pixel unpacking
* CTM application

### Deliverables

* Bitmap images rendered in page

### Exit Criteria

* Images align correctly with text

### Common Pitfalls

* Wrong color space
* Incorrect scaling
* Ignoring masks

---

## Phase 10 — Viewer Integration

### Objectives

* Make it usable

### Tasks

* Page navigation
* Zoom and pan
* Page caching
* DPI scaling

### Deliverables

* CLI or GUI viewer

### Exit Criteria

* Interactive viewing works smoothly

---

## Phase 11 — Robustness & Compliance

### Objectives

* Handle real-world PDFs

### Tasks

* Error recovery
* Broken PDF tolerance
* Incremental updates
* Performance tuning

### Deliverables

* Stable engine

### Exit Criteria

* Common PDFs render without crash

---

## Phase 12 — Advanced Features (Optional)

* Annotations
* Forms
* Transparency
* ICC color profiles
* Encryption

---