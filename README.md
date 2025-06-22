<p align="center">
  <a href="" rel="noopener">
 <img src="https://i.ibb.co/bgV3SCTb/XOR-FILTER-6-22-2025.png" alt="Project Logo"></a>
</p>

## 📖 Версия / Version
- [🇷🇺 RU](#RU)
- [🇺🇸 ENG](#ENG)

---

## 🇷🇺 RU <a name="RU"></a>

## 📝 Содержимое

- [О проекте](#aboutRU)
- [Принцип работы](#how_it_worksRU)
- [Структура XOR-фильтра](#structureRU)
- [Установка](#installationRU)
- [Пример использования](#example_usageRU)
- [Заключение](#conclusionRU)
- [Авторы](#authorsRU)

---

## 🧐 О проекте <a name="aboutRU"></a>

Проект **XORFilter** — это реализация [XOR-фильтра](https://en.wikipedia.org/wiki/XOR_filter) на языке Go, созданная в ознакомительных целях для изучения возможностей языка программирования Go. XOR-фильтр представляет собой компактную вероятностную структуру данных, которая используется для проверки принадлежности элементов к множеству с минимальной вероятностью ложных срабатываний (false positives). 

Цель проекта — продемонстрировать базовую реализацию XOR-фильтра на Go, включая этапы инициализации, построения фильтра и проверки принадлежности ключей. Проект не ставит задачу сравнения с другими структурами данных, такими как фильтр Блума, а служит для освоения работы с хэш-функциями, структурами данных и вероятностными алгоритмами в Go.

---

## 🔄 Принцип работы <a name="how_it_worksRU"></a>

XOR-фильтр использует хэш-функции и побитовые операции XOR для создания компактного представления множества ключей. Основные этапы работы:

1. **Инициализация**: 
   - Создается массив фиксированного размера, рассчитываемого как `floor(1.23 * len(Elements)) + 32`, что обеспечивает достаточное пространство для хранения отпечатков ключей.
   - Фильтр делится на три блока для вычисления позиций ключей.
   - Используется случайный `Seed` для воспроизводимости хэш-функций.

2. **Построение фильтра**: 
   - Для каждого ключа вычисляются три хэш-значения (с помощью библиотеки `xxhash`), соответствующие позициям в трех блоках фильтра.
   - Алгоритм построения решает систему уравнений с использованием XOR, чтобы назначить уникальные отпечатки (fingerprints) в массиве фильтра.
   - Если построение не удается, генерируется новый `Seed`, и процесс повторяется (до 1000 попыток).

3. **Проверка принадлежности**: 
   - Для заданного ключа вычисляются его хэш-значения и отпечаток.
   - Выполняется операция XOR над значениями в соответствующих позициях фильтра.
   - Если результат совпадает с отпечатком ключа, ключ считается принадлежащим множеству (с возможностью ложных срабатываний).

### Особенности:
- **Компактность**: Использует примерно 1.23 байта на ключ.
- **Скорость**: Операции построения и проверки имеют среднюю сложность O(1).
- **Статичность**: После построения фильтр не поддерживает добавление или удаление ключей.

---

## 🏗️ Структура XOR-фильтра <a name="structureRU"></a>

### Описание структуры

XOR-фильтр состоит из массива ячеек (`BitSet`), где каждая ячейка хранит 64-битное значение (`Bits`). Размер фильтра определяется по формуле:

```
Size = floor(1.23 * len(Elements)) + 32
```

Массив логически делится на три блока:
- **Блок 0**: Первые `BlockSize` ячеек.
- **Блок 1**: Следующие `BlockSize` ячеек.
- **Блок 2**: Оставшиеся ячейки.

Каждый ключ отображается в три позиции (по одной в каждом блоке) с помощью трёх хэш-функций. Это позволяет использовать операцию XOR для хранения и проверки отпечатков ключей.

### Пример добавления ключа

Добавление ключа в XOR-фильтр включает следующие шаги:

1. **Вычисление хэш-значений**:
   - Для ключа вычисляются три позиции (`H0`, `H1`, `H2`), каждая из которых соответствует ячейке в одном из блоков.
   - Пример для ключа `"key1"` с начальным значением `seed = 12345`:
     - `H0 = hash("key1", 12345, 0) % BlockSize`
     - `H1 = BlockSize + hash("key1", 12345, 1) % BlockSize`
     - `H2 = 2 * BlockSize + hash("key1", 12345, 2) % (Size - 2 * BlockSize)`

2. **Генерация отпечатка**:
   - Отпечаток (`Fingerprint`) — это 64-битное значение, вычисляемое как `hash("key1", 12345, 3)`.

3. **Построение фильтра**:
   - Алгоритм решает систему уравнений с использованием XOR, чтобы значения в ячейках `BitSet[H0]`, `BitSet[H1]` и `BitSet[H2]` удовлетворяли условию:

   $$BitSet[H0] \oplus BitSet[H1] \oplus BitSet[H2] = Fingerprint$$

   - Это достигается итеративным процессом, который назначает значения ячейкам, избегая конфликтов.

### Пример проверки наличия ключа

Для проверки, содержится ли ключ в фильтре:

1. **Вычисление позиций и отпечатка**:
   - Для проверяемого ключа (например, `"key1"`) вычисляются те же `H0`, `H1`, `H2` и `Fingerprint`, как при добавлении.

2. **Выполнение XOR**:
   - Вычисляется:
   
   $$result = BitSet[H0] \oplus BitSet[H1] \oplus BitSet[H2]

3. **Сравнение**:
   - Если `result == Fingerprint`, ключ считается присутствующим. Однако из-за вероятностной природы фильтра возможны ложные срабатывания.

### Схема структуры XOR-фильтра

Ниже представлена схема XOR-фильтра:

<div style="text-align: center;">
   <img src="https://i.ibb.co/zhw7c2qD/XOR-Filter-2-1.png" alt="Структура XOR-фильтра" width="150%" height="150%">
</div>

- Массив `BitSet` разделён на три блока.
- Позиции `H0`, `H1`, `H2` показывают, куда отображается ключ.
- Стрелки иллюстрируют связь между ячейками и отпечатком через операцию XOR.

---

## ⚙️ Установка <a name="installationRU"></a>

1. Склонируйте репозиторий:
   ```bash
   git clone git@github.com:lvbealr/XORFilter.git
   cd XORFilter
   ```

2. Убедитесь, что у вас установлен Go версии 1.21 или выше:
   ```bash
   go version
   ```

3. Установите зависимости:
   ```bash
   go mod tidy
   ```

4. Скомпилируйте и запустите программу:
   ```bash
   go build
   ./xorfilter
   ```
   Программа принимает ключи либо из стандартного ввода (stdin), либо из файла, указанного в аргументах командной строки.

**Зависимости**: 
- Go 1.21+
- `github.com/cespare/xxhash/v2`

---

## 🕹️ Пример использования <a name="example_usageRU"></a>

Пример кода для создания XOR-фильтра и проверки принадлежности ключей:

```go
package main

import (
    "fmt"
    "XORFilter/xorfilter"
)

func main() {
    // Инициализация фильтра с набором ключей
    keys := []string{"key1", "key2", "key3"}
    filter := xorfilter.InitializeXORfilter(keys)

    // Построение фильтра
    if err := filter.Build(); err != nil {
        fmt.Printf("Ошибка при построении фильтра: %v\n", err)
        return
    }
    fmt.Println("Фильтр успешно построен")

    // Проверка существующих ключей
    for _, key := range keys {
        if !filter.Contains([]byte(key)) {
            fmt.Printf("Ошибка: ключ %q должен быть в фильтре\n", key)
        } else {
            fmt.Printf("Ключ %q вероятно находится в фильтре\n", key)
        }
    }

    // Проверка несуществующего ключа
    testKey := []byte("nonexistent_key")
    if filter.Contains(testKey) {
        fmt.Printf("Ключ %q вероятно находится в фильтре (ложное срабатывание)\n", testKey)
    } else {
        fmt.Printf("Ключ %q отсутствует в фильтре\n", testKey)
    }
}
```

Программа также поддерживает чтение ключей из файла или stdin (см. функцию `ProcessKeys` в `xorfilter.go`).

---

## 🛠️ Заключение <a name="conclusionRU"></a>

Проект XORFilter представляет собой простую реализацию XOR-фильтра на Go и демонстрирует работу с хэш-функциями, побитовыми операциями и вероятностными структурами данных. Это отличный пример для изучения языка Go, включая такие аспекты, как работа с пакетами, управление зависимостями и обработка ошибок.

---

## ✍ Авторы <a name="authorsRU"></a>

- [@lvbealr](https://github.com/lvbealr) — Разработка и реализация.

---

## 🇺🇸 ENG <a name="ENG"></a>

## 📝 Table of Contents

- [About](#about)
- [How It Works](#how_it_works)
- [XOR Filter Structure](#structure)
- [Installation](#installation)
- [Example Usage](#example_usage)
- [Conclusion](#conclusion)
- [Authors](#authors)

---

## 🧐 About <a name="about"></a>

The **XORFilter** project is an implementation of an [XOR filter](https://en.wikipedia.org/wiki/XOR_filter) in Go, developed as an exploratory project to learn the Go programming language. The XOR filter is a compact probabilistic data structure designed for set membership queries with a low false positive rate.

The purpose of this project is to showcase a basic implementation of an XOR filter in Go, covering initialization, filter construction, and membership testing. It is not intended for performance comparisons with other data structures but rather as a tool to explore Go's features, such as hash functions, data structures, and probabilistic algorithms.

---

## 🔄 How It Works <a name="how_it_works"></a>

The XOR filter leverages hash functions and bitwise XOR operations to create a compact representation of a set of keys. The key steps are:

1. **Initialization**: 
   - An array is created with a size calculated as `floor(1.23 * len(Elements)) + 32`, providing sufficient space for key fingerprints.
   - The filter is divided into three blocks for hash position calculations.
   - A random `Seed` is used for reproducible hash functions.

2. **Filter Construction**: 
   - Three hash values are computed for each key (using the `xxhash` library), mapping to positions in the three blocks.
   - The construction algorithm solves an XOR-based system to assign unique fingerprints to the filter array.
   - If construction fails, a new `Seed` is generated, and the process retries (up to 1000 attempts).

3. **Membership Query**: 
   - For a given key, its hash values and fingerprint are computed.
   - An XOR operation is performed on the values at the corresponding filter positions.
   - If the result matches the key’s fingerprint, the key is considered present (with a chance of false positives).

### Features:
- **Compactness**: Uses approximately 1.23 bytes per key.
- **Speed**: Construction and query operations run in O(1) average time.
- **Static Nature**: Does not support adding or removing keys after construction.

---

## 🏗️ XOR Filter Structure <a name="structure"></a>

### Structure Description

The XOR filter consists of an array of cells (`BitSet`), where each cell holds a 64-bit value (`Bits`). The filter size is determined by the formula:

```
Size = floor(1.23 * len(Elements)) + 32
```

The array is logically divided into three blocks:
- **Block 0**: First `BlockSize` cells.
- **Block 1**: Next `BlockSize` cells.
- **Block 2**: Remaining cells.

Each key is mapped to three positions (one in each block) using three hash functions. This allows the use of XOR operations to store and verify key fingerprints.

### Example of Adding a Key

Adding a key to the XOR filter involves the following steps:

1. **Hash Value Calculation**:
   - Three positions (`H0`, `H1`, `H2`) are computed for the key, each corresponding to a cell in one of the blocks.
   - Example for key `"key1"` with seed `12345`:
     - `H0 = hash("key1", 12345, 0) % BlockSize`
     - `H1 = BlockSize + hash("key1", 12345, 1) % BlockSize`
     - `H2 = 2 * BlockSize + hash("key1", 12345, 2) % (Size - 2 * BlockSize)`

2. **Fingerprint Generation**:
   - The fingerprint is a 64-bit value computed as `hash("key1", 12345, 3)`.

3. **Filter Construction**:
   - The algorithm solves an XOR-based system to ensure that the values in `BitSet[H0]`, `BitSet[H1]`, and `BitSet[H2]` satisfy:

     $$BitSet[H0] \oplus BitSet[H1] \oplus BitSet[H2] = Fingerprint$$

   - This is achieved through an iterative process that assigns values to cells while avoiding conflicts.

### Example of Checking Key Presence

To check if a key is in the filter:

1. **Position and Fingerprint Calculation**:
   - For the key (e.g., `"key1"`), compute the same `H0`, `H1`, `H2`, and `Fingerprint` as during insertion.

2. **XOR Operation**:

   $$result = BitSet[H0] \oplus BitSet[H1] \oplus BitSet[H2]$$

3. **Comparison**:
   - If `result == Fingerprint`, the key is considered present. However, due to the probabilistic nature, false positives are possible.

### XOR Filter Structure Diagram

Below is a diagram of the XOR filter structure:

<div style="text-align: center;">
   <img src="https://i.ibb.co/zhw7c2qD/XOR-Filter-2-1.png" alt="XOR filter structure" width="150%" height="150%">
</div>

- The `BitSet` array is divided into three blocks.
- Positions `H0`, `H1`, `H2` indicate where the key is mapped.
- Arrows illustrate the relationship between cells and the fingerprint via XOR.

---

## ⚙️ Installation <a name="installation"></a>

1. Clone the repository:
   ```bash
   git clone git@github.com:lvbealr/XORFilter.git
   cd XORFilter
   ```

2. Ensure Go version 1.21 or higher is installed:
   ```bash
   go version
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Build and run the program:
   ```bash
   go build -o xorfilter main.go
   ./xorfilter
   ```
   The program accepts keys from either standard input (stdin) or a file specified via command-line arguments.

**Dependencies**: 
- Go 1.21+
- `github.com/cespare/xxhash/v2`

---

## 🕹️ Example Usage <a name="example_usage"></a>

Example code for creating an XOR filter and testing key membership:

```go
package main

import (
    "fmt"
    "XORFilter/xorfilter"
)

func main() {
    // Initialize filter with a set of keys
    keys := []string{"key1", "key2", "key3"}
    filter := xorfilter.InitializeXORfilter(keys)

    // Build the filter
    if err := filter.Build(); err != nil {
        fmt.Printf("Error building filter: %v\n", err)
        return
    }
    fmt.Println("Filter built successfully")

    // Check existing keys
    for _, key := range keys {
        if !filter.Contains([]byte(key)) {
            fmt.Printf("Error: key %q should be in the filter\n", key)
        } else {
            fmt.Printf("Key %q is likely in the filter\n", key)
        }
    }

    // Check nonexistent key
    testKey := []byte("nonexistent_key")
    if filter.Contains(testKey) {
        fmt.Printf("Key %q is likely in the filter (false positive)\n", testKey)
    } else {
        fmt.Printf("Key %q is not in the filter\n", testKey)
    }
}
```

The program also supports reading keys from a file or stdin (see the `ProcessKeys` function in `xorfilter.go`).

---

## 🛠️ Conclusion <a name="conclusion"></a>

The XORFilter project provides a straightforward implementation of an XOR filter in Go, showcasing hash functions, bitwise operations, and probabilistic data structures. It serves as an excellent example for learning Go, including package management, dependency handling, and error processing.

---

## ✍ Authors <a name="authors"></a>

- [@lvbealr](https://github.com/lvbealr) — Development and Implementation.

---