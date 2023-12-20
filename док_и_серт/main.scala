import java.awt.Font
import javax.swing.{JFrame, JLabel, JMenu, JMenuBar, JMenuItem, JPopupMenu, SwingConstants}
import scala.sys.exit

@main
def main(): Unit =
  val frame = JFrame("Счастливый год")
  val menuBar = JMenuBar()
  menuBar add {
    val menu = JMenu("Меню")
    val frames = Map[String, JFrame](
      "О программе" -> {
        val p = JFrame("О программе")
        p setSize (300, 200)
        p setResizable false
        p add {
          val text = "<html>Учебная программа для расчета<br>ближайшего счастливого года"
          val label = JLabel(text, SwingConstants.CENTER)
          val font = label.getFont
          label.setFont(Font(font.getFontName(), Font.PLAIN, 16))
          label
        }
        p
      },
      "Об авторе" -> {
        val p = JFrame("Об авторе")
        p setSize (300, 200)
        p setResizable false
        p add {
          val text = "<html>Стасенко К.Ю.<br>github.com/Stasenko-Konstantin<br>stasenko.kost@yandex.ru"
          val label = JLabel(text, SwingConstants.CENTER)
          val font = label.getFont
          label.setFont(Font(font.getFontName(), Font.PLAIN, 15))
          label
        }
        p
      }
    )
    menu adds (
      { val p = JMenuItem("О программе"); p.addActionListener(_ => {
          val frame = frames("О программе")
          frame setVisible !frame.isVisible
        }); p
      },
      { val a = JMenuItem("Об авторе"); a.addActionListener(_ => {
          val frame = frames("Об авторе")
          frame setVisible !frame.isVisible
        }); a
      },
      JPopupMenu.Separator(),
      { val e = JMenuItem("Выход"); e.addActionListener(_ => exit(0)); e }
    )
    menu
  }
  frame setJMenuBar menuBar
  frame setSize (400, 500)
  frame setVisible true

extension (m: JMenu)
  def adds(is: Any*): Unit = for i <- is do i match
      case i: JMenuItem => m add i
      case s: JPopupMenu.Separator => m add s
