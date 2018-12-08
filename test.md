# Godown - The Markdown Parser in Go -

Sorry, this parser is English only

## Markdown Spec
- Heading
- Emphasis
- *em*
- **strong**
- ***ems trong***
- ~~Strike through~~
- ~~**strong**~~
- List(DISC)
- List(Decimal)
- Quote
- Horizontal Line
- Code Block
- Inline Code
- `print()`
- Link
- Imange
- Table

## GoDown Spec
Sorry, This parser is English Only

---

*Do* *Not* *Use* *Japanese*
**Do** ***Not*** *Use* ***Japanese***
This `Parser` makes `AST`.

---

**Block**
- Heading
- Paragraph
- Horizontal
- DiscList
- DecimalList
- CheckList
- Code
- Blockquote
- Table

**Inline**
- Em
- Strong
- EmStrong
- Strikethrough
- InlineCode
- Text
- Html
- HtmlComment
- Image
- Link

---

```go
func main() {
    fmt.Printf("Hello, World!")
    fmt.Printf("## ignore")
    fmt.Printf("-ignore")
    fmt.Printf("`ignore`")
    fmt.Printf("`ignore``")
    fmt.Printf("*ignore*")
}
```

```rust
fn main() {
    println!("Hello World!");
}
```
