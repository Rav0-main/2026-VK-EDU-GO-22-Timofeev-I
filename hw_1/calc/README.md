# Студент

**WEB-22 Тимофеев Иван**

# Условие задачи calc

[Условие задачи](https://github.com/go-park-mail-ru/lectures/tree/master/1-basics/homework#%D1%87%D0%B0%D1%81%D1%82%D1%8C-2-calc)

# Описание исходных файлов

```
.
├── eval                                   # реализации основной функции для вычисления
│   ├── eval_calculate_rpn_test.go
│   ├── eval_calculate_test.go
│   ├── eval.go                            # реализация вычислительной логики
│   ├── eval_shunting_yard_test.go
│   └── eval_tokenize_test.go
├── go.mod
├── main.go
├── Makefile
└── README.md
```

# Сборка

```bash
# Выполняет go build -o calc.out .
make
```
