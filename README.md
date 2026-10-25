# hexlet-path-size

[![Actions Status](https://github.com/RomanVoronovskiy/go-project-242/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/RomanVoronovskiy/go-project-242/actions/workflows/hexlet-check.yml)
[![CI](https://github.com/RomanVoronovskiy/go-project-242/actions/workflows/ci.yml/badge.svg)](https://github.com/RomanVoronovskiy/go-project-242/actions/workflows/ci.yml)

Консольная утилита на Go для определения размера файла или директории.
Интерфейс построен на [urfave/cli v3](https://cli.urfave.org/v3/).
Модуль проекта — `code`.

## Установка и сборка

Нужны Go 1.23.3 или новее и GNU Make. Для линтера используется
[golangci-lint v2.15.0](https://golangci-lint.run/docs/welcome/install/local/).

```sh
git clone https://github.com/RomanVoronovskiy/go-project-242.git
cd go-project-242
go mod download
make build
./bin/hexlet-path-size --help
```

Без Make сборка выполняется так:

```sh
go build -o bin/hexlet-path-size ./cmd/hexlet-path-size
```

Go создаёт каталог `bin` при сборке. Бинарники исключены из Git:
после клонирования перед запуском нужно выполнить `make build`.
Собранный в Linux исполняемый файл получает необходимые права автоматически;
`git update-index --chmod=+x` для него не нужен.

В Windows `make build` создаёт `bin/hexlet-path-size.exe`.
В PowerShell можно собрать и запустить напрямую:

```powershell
go build -o bin/hexlet-path-size.exe ./cmd/hexlet-path-size
.\bin\hexlet-path-size.exe --help
```

## Использование

```sh
./bin/hexlet-path-size testdata/demo/hello.txt
# 6B    testdata/demo/hello.txt

./bin/hexlet-path-size testdata/demo
# 6B    testdata/demo

./bin/hexlet-path-size testdata/demo -a
# 13B   testdata/demo

./bin/hexlet-path-size testdata/demo -r
# 13B   testdata/demo

./bin/hexlet-path-size testdata/demo -r -a -H
# 27B   testdata/demo
```

Между размером и исходным путём выводится один символ табуляции.
Программа принимает ровно один путь; путь с пробелами нужно заключить в кавычки.
Флаги можно указывать как до, так и после пути.

| Флаг | Краткая форма | Значение |
| --- | --- | --- |
| `--recursive` | `-r` | Суммировать файлы на всех уровнях вложенности |
| `--human` | `-H` | Автоматически выбирать единицы размера |
| `--all` | `-a` | Учитывать скрытые файлы и директории |
| `--help` | `-h` | Показать справку |

Все три флага по умолчанию выключены. Можно явно передать
`--human=false`, `--recursive=false` или `--all=false`.

Без `-r` учитываются только файлы первого уровня, без размеров самих каталогов.
Без `-a` пропускаются элементы с именем, начинающимся с точки, включая всё
содержимое скрытых каталогов. Сам явно указанный путь обрабатывается всегда,
даже если его имя начинается с точки.

Символические ссылки учитываются по собственному размеру через `Lstat`;
программа не обходит их цели, поэтому ссылки на родительский каталог
не создают бесконечную рекурсию.

При ошибке доступа программа выводит сообщение в stderr и завершает работу
с кодом 1. Частичный размер при ошибке не возвращается.

## Формат размера

Без `-H` выводятся байты: `25165824B`.
С `-H` используются степени 1024 и единицы `B`, `KB`, `MB`, `GB`,
`TB`, `PB`, `EB`. Для байтов дробная часть отсутствует, для остальных
единиц выводится один знак после точки:

```text
123       → 123B
1024      → 1.0KB
1234567   → 1.2MB
25165824  → 24.0MB
```

## Использование как библиотеки

```go
import "code"

size, err := code.GetPathSize(path, recursive, human, all)
```

Функция `GetPathSize(path string, recursive, human, all bool) (string, error)`
возвращает только строку с размером. Формирование вывода и обработка
аргументов находятся в `cmd/hexlet-path-size`.

## Проверки

```sh
make build
make test
make coverage
make lint
make lint-fix
```

`make lint-fix` форматирует код и запускает автоматическое исправление
предупреждений. Ограничение цикломатической сложности — 10.
CI собирает программу из исходников, проверяет CLI, тесты, покрытие от 80%
и линтер. Тесты используют временные каталоги и фикстуры в `testdata`.
Проверки Unix-прав доступа и недоступных на машине символических ссылок
могут пропускаться в Windows.

Для запуска из Make: `make run ARGS="testdata/demo -r -a -H"`.

## Демонстрация

[Посмотреть запись на asciinema](https://asciinema.org/a/uNa8luaGlOXj9UlH)

[![asciicast](https://asciinema.org/a/uNa8luaGlOXj9UlH.svg)](https://asciinema.org/a/uNa8luaGlOXj9UlH)

Запись установки зависимостей, сборки, справки и запуска с разными флагами
находится в [docs/demo.cast](docs/demo.cast) (формат asciicast v2).

Локальную копию можно воспроизвести командой `asciinema play docs/demo.cast`.
