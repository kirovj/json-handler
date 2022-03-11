enum Token {
    Space,
    Colon,        // :
    Comma,        // ,
    Quote,        // "
    OpenBracket,  // [
    CloseBracket, // ]
    OpenCurly,    // {
    CloseCurly,   // }
    Backslash,    // \
    NewlineN,     // \n
    NewlineR,     // \r
    UnicodeFlag,  // u
    String(String),
    Number(f64),
    Boolean(bool),
    Null,
}

fn main() {
    println!("Hello, world!");
}
