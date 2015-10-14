# Go Parsec2

Go Parsec2 改写了 [goParsec](http://github.com/Dwarfartisan/goParsec) 。在性能上作出一定让步后，对结构和形式做出了改良。

Go P 2中，bind 不再是一个 P 算子，而是所有 P 算子的 Monad 特征。所有的 P 算子都组合为带 Bind/Then/Over 的结构。 Go P 2 提供了一些方法简化这些封装操作。

提供了 Do 形式。
