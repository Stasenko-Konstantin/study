# sfcalc (deprecated)
Небольшой функциональный язык программирования для простых математических расчетов


### Примеры
< - приглашение REPL\`а                                                                 
\> - вывод интерпретатора

```haskell
< 1 + 1 
> 2

< a = 5!
> a = 25

< A = 26 -- язык регистронезависим, а все "переменные", на самом деле константы
> Override error: a = 26, pos = 1

< if not a == 5! then true else false -- ветка else обязательная
> false

< 5 == 5.0
> true

< 5 === 5.0
> false

< f = @(a b c)(if a == b then sqrt(c) else с^a) -- лямбда
> f = Func => (Num Num Num) -> Num

< f2 = @()()
> f2 = Func => Nil

< f2()
> nil

< log(-1)
> nil

< f3 = @()\
| ()
> f3 = Func => Nil
```
