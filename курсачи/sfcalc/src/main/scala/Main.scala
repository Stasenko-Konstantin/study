//
// repl # условие
//

import scala.io.StdIn.readLine

object Main extends App:

  repl()

  def repl(): Unit =

    print("< ")
    val input = Option(" " + readLine().toLowerCase)
    val tokens = Scanner.scan(input)
    tokens.get.foreach(println)
    val expr = Parser.parse(tokens.get)
    println(AstPrinter.print(expr))
    if input.get == ":q" then repl() // == -> !=

  def error(start: Int, msg: String, content: String): Unit =
    val out = msg + content + ", pos = " + start.toString
    System.err.println(out)