object Scanner {

  def scan(code: Option[String]): Option[List[Token]] = {

    def help(tokens: List[Token], start: Int): Option[List[Token]] = {

      def minhelp(t: List[Token], s: Int = start): Option[List[Token]] = help(t, s)

      def addNum(s: Int, r: String = "", isFloat: Boolean = false): Option[List[Token]] = {
        def result(): Option[List[Token]] = minhelp(Token(if isFloat then Token.FLOAT
                                            else Token.INT, r.reverse, s + 1) :: tokens, s = s)
        if s < 0 then return result()
        val c = code.get(s)
        c match {
          case '.' => addNum(s - 1, r + c, true)
          case _ => if !c.isDigit then result() else addNum(s - 1, r + c, isFloat)
        }
      }

      def addName(s: Int, r: String = ""): Option[List[Token]] = {
        def result(): Option[List[Token]] = {
          var t: TokenType = Token.NIL
          r.reverse match {
            case "if" => t = Token.IF
            case "then" => t = Token.THEN
            case "else" => t = Token.ELSE
            case "true" => t = Token.TRUE
            case "false" => t = Token.FALSE
            case "and" => t = Token.AND
            case "or" => t = Token.OR
            case "not" => t = Token.NOT
            case "sqrt" => t = Token.SQRT
            case _ => t = Token.NAME
          }
          minhelp(Token(t, r.reverse, s + 1) :: tokens, s = s)
        }
        if s < 0 then return result()
        val c = code.get(s)
        if c.isLetter then addName(s - 1, r + c) else result()
      }

      def addEqual(s: Int): Option[List[Token]] = {
        def result(t: TokenType, r: String, offset: Int = 1): Option[List[Token]] =
          minhelp(Token(t, r, s - 1) :: tokens, s = s - offset)
        val c = code.get(s - 1)
        c match {
          case '<' => result(Token.LEQUAL, "<=", 2)
          case '>' => result(Token.GEQUAL, ">=", 2)
          case '=' => if code.get(s - 2) == '=' then result(Token.TEQUAL, "===", 3) else result(Token.DEQUAL, "==", 2)
          case _ => code.get(s) match {
            case '=' => result(Token.EQUAL, "=")
            case '<' => result(Token.LESS, "<")
            case '>' => result(Token.GREAT, ">")
          }
        }
      }

      if start < 0 then return Option(tokens)
      val c = code.get(start)
      code match {
        case Some(code) => c match {
          case '+' => minhelp(Token(Token.PLUS, "+", start) :: tokens, s = start - 1)
          case '-' => minhelp(Token(Token.MINUS, "-", start) :: tokens, s = start - 1)
          case '*' => minhelp(Token(Token.STAR, "*", start) :: tokens, s = start - 1)
          case '/' => minhelp(Token(Token.SLASH, "/", start) :: tokens, s = start - 1)
          case '\\' => minhelp(Token(Token.RSLASH, "\\", start) :: tokens, s = start - 1)
          case '!' => minhelp(Token(Token.FACTORIAL, "!", start) :: tokens, s = start - 1)
          case '^' => minhelp(Token(Token.POWER, "^", start) :: tokens, s = start - 1)
          case '(' => minhelp(Token(Token.LPAREN, "(", start) :: tokens, s = start - 1)
          case ')' => minhelp(Token(Token.RPAREN, ")", start) :: tokens, s = start - 1)
          case '@' => minhelp(Token(Token.LAMBDA, "@", start) :: tokens, s = start - 1)
          case '=' | '<' | '>' => addEqual(start)
          case '\n' => minhelp(Token(Token.SEMICOLON, "\\n", start) :: tokens, s = start - 1)
          case ' ' | '\t' | '\r' => minhelp(tokens, s = start - 1)
          case _ if c.isDigit => addNum(start)
          case _ if c.isLetter => addName(start)
          case _ => Main.error(start, "Syntax error: ", c.toString)
                    Option(tokens)
        }
        case None => None
      }
    }

    val len = code.get.length - 1
    Option(help(Nil, len).get ::: List(Token(Token.EOF, "\u0000", -1)))
  }
}