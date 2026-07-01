# set — довідник

Повний довідник пакета `set`: ментальна модель, ідентичність елементів,
конструювання, повна поверхня методів, пакетні генерики, ітерація, JSON і
практичні рецепти.

Англійська версія: **[DOC.md](DOC.md)**.

## Зміст

- [Ментальна модель](#ментальна-модель)
- [Ідентичність елементів](#ідентичність-елементів)
- [Конструювання](#конструювання)
- [Базові операції](#базові-операції)
- [Алгебра множин](#алгебра-множин)
- [Відношення](#відношення)
- [Ітерація й впорядкування](#ітерація-й-впорядкування)
- [Функціональні помічники](#функціональні-помічники)
- [JSON](#json)
- [Конкурентність](#конкурентність)
- [Рецепти й поради](#рецепти-й-поради)

## Ментальна модель

`set` — це маленький, швидкий, узагальнений `Set` для Go: невпорядкована колекція
унікальних елементів `comparable`-типу, побудована прямо на вбудованій мапі.

Визначальний вибір — що **ідентичність елемента є власною рівністю мови (`==`)**.
Два елементи однакові тоді й лише тоді, коли вони рівні за порівнянням, а
унікальність вирішує рантаймова мапа — немає гешування, рефлексії й кастомного
контракту рівності. Як наслідок:

- `Set` ніколи не може тихо втратити елемент через колізію гешу.
- `Len` завжди відображає справжню кількість різних елементів.
- Поведінка відповідає вашій інтуїції щодо `==` для типу елемента.

```go
import "github.com/goloop/set/v2"
```

## Ідентичність елементів

Обмеження `comparable` допускає числові види, `string`, `bool`, вказівники,
канали, інтерфейси й будь-яку структуру чи масив, поля яких самі порівнювані.

Ідентичність — це `==`, тож структури порівнюються поле за полем, за значенням:

```go
type Address struct{ City string }
type User struct{ Name string; Age int; Address Address }

users := set.New(
    User{"John", 21, Address{"Kyiv"}},
    User{"Bob", 22, Address{"Chernihiv"}},
    User{"John", 21, Address{"Kyiv"}}, // дублікат -> згортається
)
users.Len() // 2
```

Структура, що містить **вказівник**, порівнюється за цим вказівником, а не за
значенням, на яке він указує: дві структури з різними вказівниками — різні
елементи, навіть якщо значення, на які вони вказують, рівні. Зберігайте значення,
а не вказівники, коли хочете ідентичність за значенням.

Зрізи, мапи й функції не порівнювані й не можуть бути елементами напряму. Щоб
дедуплікувати такі значення, виведіть порівнюваний ключ (`string` чи структуру з
порівнюваних полів) і збудуйте `Set` цього ключа.

## Конструювання

```go
func New[T comparable](items ...T) *Set[T]
func NewWithCapacity[T comparable](capacity int, items ...T) *Set[T]
func Collect[T comparable](seq iter.Seq[T]) *Set[T]
```

`New` виводить тип елемента з аргументів або бере його явно, коли порожньо
(`set.New[int]()`). `NewWithCapacity` попередньо виділяє мапу, коли розмір
відомий. `Collect` будує множину з будь-якої `iter.Seq[T]`.

Нульове значення придатне напряму — `var s set.Set[int]` є порожньою, готовою до
використання множиною (перша вставка виділяє мапу). `New` усе одно кращий, коли
розмір відомий.

```go
ints  := set.New[int]()
words := set.New("one", "two", "three")
keys  := set.Collect(maps.Keys(m))
var s set.Set[int] // теж валідно
```

## Базові операції

```go
func (s *Set[T]) Add(items ...T)
func (s *Set[T]) AddSeq(seq iter.Seq[T])
func (s *Set[T]) Delete(items ...T)
func (s *Set[T]) Contains(item T) bool
func (s *Set[T]) ContainsAll(items ...T) bool
func (s *Set[T]) ContainsAny(items ...T) bool
func (s *Set[T]) Len() int
func (s *Set[T]) IsEmpty() bool
func (s *Set[T]) Clear()
func (s *Set[T]) Copy() *Set[T]
func (s *Set[T]) Pop() (T, bool)
func (s *Set[T]) Elements() []T
func (s *Set[T]) Append(others ...*Set[T])
func (s *Set[T]) Overwrite(items ...T)
```

`Add`/`Delete` варіативні. `AddSeq` додає всі значення з `iter.Seq[T]`. `Pop`
видаляє й повертає довільний елемент (`ok=false`, коли порожньо). `Elements`
повертає членів невпорядкованим зрізом. `Append` зливає інші множини в цю на
місці; `Overwrite` замінює вміст заданими елементами.

```go
ints.Add(1, 2, 3, 4)
ints.Delete(1, 2)
ints.Contains(3)          // true
ints.ContainsAll(3, 4)    // true
s.AddSeq(slices.Values(items))
```

## Алгебра множин

```go
func (s *Set[T]) Union(others ...*Set[T]) *Set[T]
func (s *Set[T]) Intersection(others ...*Set[T]) *Set[T] // псевдонім: Inter
func (s *Set[T]) Difference(others ...*Set[T]) *Set[T]   // псевдонім: Diff
func (s *Set[T]) SymmetricDifference(others ...*Set[T]) *Set[T] // псевдонім: Sdiff
```

Кожна повертає нову множину й приймає кілька операндів одразу
(`a.Union(b, c, d)`):

```go
a := set.New(1, 3, 5, 7)
b := set.New(0, 2, 4, 7)

set.Sorted(a.Union(b))               // [0 1 2 3 4 5 7]
set.Sorted(a.Intersection(b))        // [7]
set.Sorted(a.Difference(b))          // [1 3 5]
set.Sorted(a.SymmetricDifference(b)) // [0 1 2 3 4 5]
```

## Відношення

```go
func (s *Set[T]) Equal(other *Set[T]) bool
func (s *Set[T]) IsSubset(other *Set[T]) bool         // псевдонім: IsSub
func (s *Set[T]) IsSuperset(other *Set[T]) bool       // псевдонім: IsSup
func (s *Set[T]) IsProperSubset(other *Set[T]) bool
func (s *Set[T]) IsProperSuperset(other *Set[T]) bool
func (s *Set[T]) IsDisjoint(other *Set[T]) bool
```

`IsSubset`/`IsSuperset` — це **нестрогі** відношення: множина є підмножиною й
надмножиною самої себе, відповідно до стандартних математичних означень.
Використовуйте варіанти `Proper` для строгих відношень.

```go
a := set.New(1, 2, 3)
b := set.New(1, 2, 3, 4, 5)

a.IsSubset(b)               // true  (a ⊆ b)
a.IsProperSubset(b)         // true  (a ⊊ b)
b.IsSuperset(a)             // true  (b ⊇ a)
a.Equal(set.New(3, 2, 1))   // true  (порядок не має значення)
a.IsDisjoint(set.New(8, 9)) // true
```

## Ітерація й впорядкування

Порядок ітерації множини **невизначений**.

```go
func (s *Set[T]) Iter() iter.Seq[T]
func (s *Set[T]) Elements() []T
func (s *Set[T]) Sorted(cmp func(a, b T) int) []T
func Sorted[T cmp.Ordered](s *Set[T]) []T
```

Використовуйте `Iter` для циклу `range`, `Elements` для невпорядкованого зрізу
чи `Sorted`, коли потрібен стабільний порядок:

```go
for v := range s.Iter() {
    _ = v // невизначений порядок
}

set.Sorted(s)                                 // [1 2 3]  (природний порядок)
s.Sorted(func(a, b int) int { return b - a }) // [3 2 1]  (власний порядок)
```

Пакетний `set.Sorted` працює для елементів типу `cmp.Ordered` без аргументу;
метод `Sorted` бере функцію порівняння (той самий контракт, що й `cmp.Compare`)
для будь-якого іншого порядку.

## Функціональні помічники

Методи, що зберігають тип елемента, плюс пакетні генерики, що можуть його
змінити:

```go
func (s *Set[T]) Map(fn func(item T) T) *Set[T]
func (s *Set[T]) Filter(fn func(item T) bool) *Set[T]
func (s *Set[T]) Filtered(fn func(item T) bool) []T
func (s *Set[T]) Reduce(fn func(acc, item T) T) T
func (s *Set[T]) Any(fn func(item T) bool) bool
func (s *Set[T]) All(fn func(item T) bool) bool

func Map[T, R comparable](s *Set[T], fn func(item T) R) *Set[R]
func Reduce[T comparable, R any](s *Set[T], fn func(acc R, item T) R) R
func Fold[T comparable, R any](s *Set[T], initial R, fn func(acc R, item T) R) R
```

```go
s := set.New(1, 2, 3, 4, 5)

even := s.Filter(func(v int) bool { return v%2 == 0 }) // {2, 4}
doubled := s.Map(func(v int) int { return v * 2 })     // той самий тип елемента

// Функція Map може змінити тип елемента.
labels := set.Map(s, func(v int) string {
    if v%2 == 0 { return "even" }
    return "odd"
}) // {"odd", "even"}

sum := s.Reduce(func(acc, v int) int { return acc + v })           // 15
product := set.Fold(s, 1, func(acc, v int) int { return acc * v }) // 120
```

`Reduce` (метод) стартує з нульового значення; `Fold` бере явний старт і може
акумулювати в інший тип. `Any`/`All` — прості лінійні проходи.

## JSON

`Set` реалізує стандартні інтерфейси `encoding/json`:

```go
s := set.New(1, 2, 3)

data, _ := json.Marshal(s) // напр. [1,2,3] (порядок невизначений)

var back set.Set[int]
_ = json.Unmarshal(data, &back)
back.Equal(s) // true
```

Множина маршалиться в JSON-масив і демаршалиться з нього, дедуплікуючи на вході.

## Конкурентність

`Set` **не** безпечний для конкурентного використання кількома горутинами, точно
як вбудована мапа, на якій він побудований. Якщо множина спільна й хоч одна
горутина її мутує, захистіть доступ власною синхронізацією:

```go
type SafeSet[T comparable] struct {
    mu sync.RWMutex
    s  *set.Set[T]
}

func (c *SafeSet[T]) Add(items ...T) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.s.Add(items...)
}
```

Тримання ядра несинхронізованим уникає блокування на кожну операцію й дозволяє
викличнику, який знає модель конкурентності застосунку, обрати правильну
стратегію.

## Рецепти й поради

**Дедуплікуйте зріз.** `set.Collect(slices.Values(xs)).Elements()` (чи
`set.Sorted(...)`) повертає різні значення.

**Ключуйте непорівнювані значення.** Для зрізів/мап виведіть порівнюваний ключ
(`fmt.Sprint`, геш-рядок чи структуру з порівнюваних полів) і зберігайте `Set`
цього ключа.

**Ідентичність за значенням для структур.** Зберігайте значення структур, а не
вказівники, щоб рівні записи згорталися в один елемент.

**Стабільний вивід.** Порядок ітерації невизначений — пропускайте результати
крізь `set.Sorted` (чи метод `Sorted`), коли порядок спостережний, напр. у тестах
чи серіалізованому виводі.

**Попередньо виділяйте, коли знаєте розмір.** `NewWithCapacity(n, …)` уникає
повторного росту мапи при додаванні багатьох елементів.
