# spike-interpreter-go
Learn by doing project - interpreter written in Go for my own language "Spike"

## Language features

Immutable variables

```
let Name = "kenny"
Name = "def" // Error!
```

Tuples
```
let Person = (123, "Lukasz")
```

Atoms
```
let Person = (person, 123, "Lukasz")
let Place = (place, 222, "Krakow")
```

Pattern matching
```
let Person = (person, 123, "Lukasz")

case Person of
    (person, Id, Name) -> {
        printf("Person: %s", Name)
    }
    (place, Id, Name) -> {
        printf("Place: %s", Name)
    }
end

```

First class functions
```
let f = fn(x) -> {
    return x + 5
}

let result = f(10)
```

Lightweight processes 
```
process User {
    let name = "lukasz"
    
    cast changeName(newName) {
        return User{
            name: newName
        }
    }
    
    call getName() {
        return this.name 
    }
}

let u = spawn User()

u.changeName("newName")
```

## ToDo

- [x] Lexing of all basic mathematical operators
- [x] Parsing booleans
- [x] Parsing grouped expressions
- [x] Parsing if-else expressions
- [x] Parsing function literals
- [x] Parsing call expressions