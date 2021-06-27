# Modellbasierte Softwareentwicklung

Das Ziel dieses Projekts ist es, eine kleine, aber vollständige, Programmiersprache zu entwickeln. Dabei geht es um das Auslesen des Codes als ein String in eine Datenstruktur (parsing), die Validierung der Typen (type-checking) und die Möglichkeit diesen Code auszuführen. Die Sprache ist angelehnt an C, jedoch wurden nur die wichtigsten Konstrukte umgesetzt. 

## Beschreibung der Sprache

Die Programmiersprache ist aus Sicht der Syntax an C angelehnt. Es folgt ein Beispielcode, der die Features der Sprache zeigt:


```c=
a = 123;
b = "abc";
c = true;
d = 4.2;

if (c) {
    println("c is true");
}

if (a == 123) {
    println("a is 123");
}

if ((c && true) == true) {
    println("c && true");
}

if (b == "abc") {
    println("b is abc");
}

println(b + "123");

for (;false;) {
}

for (e = 1; e < 4; e = e + 1) {
    println("e");
}

input = readln();
println(input);
```

### Variablen und Werte

Ein Wert in der Sprache hat einen der Typen `String`, `Int`, `Float` oder `Boolean`. Die Definition eigener Typen wird nicht unterstützt. Es gibt kein spezielles Schlüsselwort um eine Variable erstmals zu erstellen denn dies geschieht implizit bei der ersten Zuweisung. Es gibt jedoch „Scopes“, was bedeutet, dass Variablen nicht mehr gelten, nachdem man den Block, in dem sie definiert wurden, wieder verlässt.

Der Typ einer Variable wird ihr bei der ersten Zuweisung verliehen. Die Typen müssen also nicht explizit angegeben werden. Jedoch kann sich der Typ einer Variable nach der Zuweisung nicht mehr ändern. Es handelt sich also nicht um „duck typing“ sondern um statische Typisierung.

### Bedingungen

Einzelne Codeabschnitte können bedingt ausgeführt werden, indem sie mit einer „if“-Bedingung abgesichert werden. Hier gibt es in Klammern eine Bedingung und danach einen Block Code in geschweiften Klammern. Wird diese Bedingung während der Ausführung erfüllt, so wird der in geschweiften Klammern stehende Code ausgeführt. Ansonsten wird er übersprungen und die nachfolgenden Operationen ausgeführt.

### Schleifen

In unserer Programmiersprache gibt es eine „for“-Schleife. Da es keine anderen Schleifentypen und auch keine Rekursion gibt, ist diese „for“-Schleife auch der einzige Weg, um Aktionen eine beliebige Anzahl mal zu wiederholen.

Die „for“-Schleife ist syntaktisch ähnlich wie in C. Es gibt 3 verschiedene Ausdrücke innerhalb der Klammern, die mit einem Semikolon getrennt sind. Mit dem ersten kann man eine Variable initialisieren. Hier ist es erzwungen, dass es ein Ausdruck der Form „x = Wert“ ist. Danach folgt ein Ausdruck, der festlegt, wann die Schleife abbrechen soll. Dieser Ausdruck wird nach jedem Schleifendurchlauf ausgeführt. Ist dieser Wert „false“, dann wird die Schleife abgebrochen. Der dritte Ausdruck erfordert, wie auch schon der erste Ausdruck, eine Beschreibung einer Variable. Dieser wird am Ende jedes Schleifendurchlaufs ausgeführt und kann beispielsweise dazu verwendet werden, um eine Iterationsvariable um 1 zu erhöhen. Danach folgt in geschweiften Klammern ein Block Code.

Die 3 Ausdrücke zwischen den Klammern können jeweils weggelassen werden. Bei dem ersten und dritten Ausdruck bedeutet dies das einfach keine variable initialisiert bzw. erhöht wird. Wenn der mittlere Ausdruck weggelassen wird, dann bedeutet das, dass es sich um eine Endlosschleife handelt.

### Funktionsaufrufe

Es gibt in der Sprache 2 „hartcodierte“ Funktionen. Unterstützt werden „readln“ zum Auslesen einer Zeile aus „stdin“ und „println“ zum Ausgeben einer Zeile auf „stdout“. Auf diesem Weg kann man mit dem Programm auf der Konsole kommunizieren und Eingaben tätigen sowie Ausgaben auslesen. „println“ nimmt hierbei einen String an, der dann ausgegeben wird. „readln“ hat dementsprechend einen Rückgabewert von String und nimmt keine Parameter an.

### Operatoren

Einzelne Ausdrücke können mit einem Operator verbunden werden. Dies ähnelt theoretisch einem Funktionsaufruf der zwei Parameter hat. Jedoch wird das nicht mit einem Namen aufgerufen, sondern mit einem Symbol, welches zwischen den beiden Ausdrücken steht. Die Parameter sind auch hier angelehnt an C. Sie teilen sich in 4 verschiedene Kategorien auf:
-	Die Gleichheitsoperatoren (==, !=). Beim testen auf Gleichheit muss der linke und der rechte Ausdruck den gleichen Typ haben. Beispiel: „123 == 123“. Dieser Ausdruck gibt einen Boolean zurück.
-	Der boolesche „und“ und „oder“ Operator (&&, ||). Mit diesen kann man einzelne boolesche Werte miteinander verketten.
-	Die Vergleichsoperatoren für Zahlen (>, <, >=, <=). Mit ihnen kann man Zahlenwerte vergleichen.
-	Die arithmetischen Operationen (+, -, *, /). Sie können verwendet werden um mit den Zahlenwerten zu rechnen.
-	Der Operator für String-Konkatenation (+). Mit ihm können mehrere Strings verbunden werden. Hier handelt es sich um das gleiche Symbol wie bei der Addition von Zahlen. Es hängt also von den Typen ab, was gemacht wird.

## Ziele der einzelnen Phasen

Das Programm setzt sich aus 3 Teilen zusammen. Der Parser, der Type-Checker und die Code-Ausführung.

### Parsen

Das Ziel des Parsens ist es, die Eingabe (also ein großer String) in einen „Abstract Syntax Tree“ umzuwandeln. Dabei sollen alle korrekten Programme richtig erkannt werden und alle fehlerhaften abgelehnt werden.

### Type-Checking

Der Type-Checker ist dazu da, den ausgelesenen AST auf semantische Probleme zu testen. Hierbei soll herausgefunden werden, ob eine Ausführung aus Sicht des Typsystems Sinn ergibt.

### Code-Ausführung
Wenn der Type-Checker keine Probleme festgestellt hat, muss der geschriebene Code nur noch ausgeführt werden. Hierbei kommt wieder der AST, den der Parser generiert hat, zum Einsatz. Dieser wird Schritt für Schritt evaluiert und somit die ursprünglich im Code angegebenen Operationen ausgeführt.

## Angewandte Methoden

### Parserkombinatoren

Einige der Syntaxkonstrukte bestehen aus mehrere nacheinanderfolgenden Teilen. Ein Beispiel dafür ist die „if“-Bedingungen, die aus dem Wort „if“, 4 verschiedenen klammern („(){}“) und zwischen in Klammern ein Ausdruck bzw. ein Block von Code.

Der „Parser“ Typ stellt einen einzelnen Teil der Syntax dar, die ausgelesen werden soll:
```go
type Parser func(string) (string, error)
```
Dieser Typ ist eine Funktion, die den Code als String-Parameter annimmt. Ein Parser kann jedoch nicht immer den ganzen String auswerten, da ja nach dem relevanten Teil noch weitere folgen kann. Beispielsweise möchte ein String-Parser ja nur einen String auslesen, den Rest des Codes jedoch ignorieren. Aus diesem Grund wird der restliche Code als String zusammen mit einem Fehler zurückgegeben. Der Aufruf von einer Parser-Funktion nimmt also am Anfang des Codes einen Teil des String weg.

Da Go keine Generics hat kann ein Parser-Funktion keinen speziellen Typ zurückgeben, da die Signatur der Funktion sonst dem "Parser"-Typ nicht mehr entspricht. Deswegen werden Resultate von einzelnen Parser-Funktionen hier als Out-Parameter übergeben. Beispielsweise schreibt die Funktion `name` den gelesenen Namen an die Adresse, die mittels dem Parameter `out` übergeben wurde. Da immer nur einzelne Codeteile wie z.B. eine "for"-Schleife mit Parserkombinatoren ausgewertet werden, funktioniert dies gut.

#### Lesen von einzelnen Begriffen

```go
func token(t string) Parser
func name(out *string) Parser
```

Die Funktion `token` akzeptiert nur eine spezielles Wort am Anfang des Strings. Dieses wird der Funktion übergeben.

`name` liest einen einzelnen Namen (bsp. Variablenname) aus und schreib ihn an die Adresse in `out`.

#### Lesen von Ausdrücken

```go
func expr(out *Expr) Parser
func block(out *Block) Parser
```

Mit `expr` kann man einen ganzen Ausdruck auslesen. Dies umfasst alle Begriffe bei denen ein Wert zurückgegeben wird (z.B. "asd" oder "1 + (4 * 2)"). Bei `block` geschieht dies analog miteinem Code-Block (einer folge von einzelnen Befehlen mit `;` getrennt).

#### Kombination mit anderen Parsern

```go
func pfunc(out *Expr, fn func(string) (string, Expr, error)) Parser
func sequence(parsers ...Parser) Parser
func alternative(parsers ...Parser) Parser
func opt(p Parser) Parser
```

`pfunc` erlaubt es eine Funktion, die zusätzlich zu den restlichen Code und einem Fehler noch eine `Expr` zurückgibt in einen Parser umzuwandeln. Somit kann man diese Funktionen mit anderen kombinieren.

Mehrere aufeinanderfolgende Parser können mit `sequence` verbunden werden. Sobald einer der Parser fehlschlägt gibt dieser Parser einen Fehler zurück.

`alternative` wählt den ersten erfolgreichen Parser und gibt dessen Ergebnis zurück. Hat keiner der Parser Erfolg, dann wird ein Fehler zurückgegeben

`opt` stellt einen optionalen Parser dar. Schlägt der übergebene Parser fehl, dann wird kein Code konsumiert, es gibt jedoch keinen Fehler.

```go=
func ParseIf(code string) (string, Expr, error) {
	if_ := If{}
	code, err := sequence(token("if"), token("("), expr(&if_.Condition), token(")"), token("{"), block(if_.Body), token("}"))(code)

	if err != nil {
		return code, nil, err
	}

	return code, if_, nil
}
```

### Regexp
Um beim parsen bestimmte Muster zu erkennen, wird Gebrauch von regulären Ausdrücken gemacht. Beispielsweise werden sie dazu verwendet, um Variablennamen, Ganzzahlen oder Gleitkommazahlen aus dem Code zu extrahieren.
```go=
var nameRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*`)

func ParseName(code string) (string, string, error) {
	codeWithoutWhitespace := stripWhitespaceLeft(code)
	name := nameRegex.FindString(codeWithoutWhitespace) //extracting the variable name

	if name == "" {
		return "", "", &ParseError{Message: "Couldn't parse the name"}
	}

	return codeWithoutWhitespace[len(name):], name, nil
}
```
### Polymorphie
Im Type-Checker und bei der Code-Ausführung wird stark auf das Konzept der Polymorphie zurückgegriffen. Der AST ist ein Konstrukt, welches aus vielen verschiedenen Ausdruckstypen besteht. Jeder dieser Ausdruckstypen implementiert das „Expr“-Interface („Expr“ steht hier für „Expression“), welches Funktionen enthält, die für die Weiterverarbeitung des ASTs wichtig sind. Diese Funktionen werden von allen Ausdruckstypen unterschiedlich implementiert. So können zur Laufzeit immer genau die Operationen ausgeführt werden, die zu dem jeweiligen Ausdruck passen. Das folgende Beispiel zeigt zwei unterschiedliche Implementierungen und die Verwendung der „eval“-Funktion, die für die Code-Ausführung zuständig ist.

```go
type Expr interface {
	Print() string
	Eval() interface{}
	Type() Type
}

func (b Block) Eval() interface{} {
	outerscopeVars := map[string]interface{}{}
	for k, v := range variables {
		outerscopeVars[k] = v
	}
	for _, expr := range b.Statements {
		expr.Eval()
	}
	variables = outerscopeVars
	return nil
}

func (f For) Eval() interface{} {
	for f.Init.Eval(); f.Condition.Eval().(bool); f.Advancement.Eval() {
		f.Body.Eval()
	}
	return nil
}
```
