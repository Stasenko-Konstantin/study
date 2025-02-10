sealed abstract class TokenType

class Token(val token: TokenType, val content: String, pos: Int):
  override def toString: String =
    def concat(t: String) = t + ", " + content + ", pos = " + pos.toString
    token match {
      case Token.PLUS => concat("Plus")
      case Token.MINUS => concat("Minus")
      case Token.STAR => concat("Star")
      case Token.SLASH => concat("Slash")
      case Token.RSLASH => concat("RSlash")
      case Token.FACTORIAL => concat("Factorial")
      case Token.POWER => concat("Power") //^
      
      case Token.EQUAL => concat("Equal")
      case Token.DEQUAL => concat("DEqual")
      case Token.TEQUAL => concat("TEqual")
      case Token.LESS => concat("Less")
      case Token.GREAT => concat("Great")
      case Token.LEQUAL => concat("LEqual")
      case Token.GEQUAL => concat("GEqual")
      
      case Token.TRUE => concat("True")
      case Token.FALSE => concat("False")
      case Token.IF => concat("If")
      case Token.THEN => concat("Then")
      case Token.ELSE => concat("Else")
      case Token.AND => concat("And")
      case Token.OR => concat("Or")
      case Token.NOT => concat("Not")
      
      case Token.NAME => concat("Name")
      case Token.LAMBDA => concat("Lambda")
      case Token.SQRT => concat("Sqrt")
      
      case Token.INT => concat("Int")
      case Token.FLOAT => concat("Float")
      
      case Token.LPAREN => concat("LParen")
      case Token.RPAREN => concat("RParent")

      case Token.SEMICOLON => concat("Semicolon")
      case Token.EOF => concat("EOF")

      case _ => concat("Nil")
    }

object Token:
  case object PLUS extends TokenType
  case object MINUS extends TokenType
  case object STAR extends TokenType
  case object SLASH extends TokenType
  case object RSLASH extends TokenType
  case object FACTORIAL extends TokenType
  case object POWER extends TokenType //^
  
  case object EQUAL extends TokenType
  case object DEQUAL extends TokenType
  case object TEQUAL extends TokenType
  case object LESS extends TokenType
  case object GREAT extends TokenType
  case object LEQUAL extends TokenType
  case object GEQUAL extends TokenType
  
  case object FALSE extends TokenType
  case object TRUE extends TokenType
  case object IF extends TokenType
  case object THEN extends TokenType
  case object ELSE extends TokenType
  case object AND extends TokenType
  case object OR extends TokenType
  case object NOT extends TokenType
  
  case object NAME extends TokenType
  case object LAMBDA extends TokenType
  case object SQRT extends TokenType
  
  case object INT extends TokenType
  case object FLOAT extends TokenType
  
  case object LPAREN extends TokenType
  case object RPAREN extends TokenType

  case object SEMICOLON extends TokenType
  case object EOF extends TokenType

  case object NIL extends TokenType