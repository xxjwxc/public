# consul-kv
A consul key/value tool for golang

## Usage

### New Config
```golang
conf := NewConfig()
```

### With Options
```golang
conf := NewConfig(
    WithPrefix(prefix),             // consul kv prefix
    WithAddress(address),           // consul address
    WithAuth(username, password),   // cosul auth
    WithToken(token),               // cousl token
    WithLoger(loger),               // loger
)

```

### Init
```golang
if err := conf.Init();err !=nil {
    return err
}
```

### Put
```golang
if err := conf.Put(key, value);err !=nil {
    return err
}
```

### Delete
```golang
if err := conf.Delete(key);err !=nil {
    return err
}
```

### Get
```golang
// scan
if err := conf.Get(key).Scan(x);err !=nil {
    return err
}

// get float
float := conf.Get(key).Float()

// get float with default
float := conf.Get(key).Float(defaultFloat)

// get int
i := conf.Get(key).Int()

// get int with default
i := conf.Get(key).Int(defaultInt)

// get uint
uInt := conf.Get(key).Uint()

// get uint with default
uInt := conf.Get(key).Uint(defaultUint)

// get bool
b := conf.Get(key).Bool()

// get bool with default
b := conf.Get(key).Bool(defaultBool)

// get []byte
bytes := conf.Get(key).Bytes()

// get uint with default
bytes := conf.Get(key).bytes(defaultBytes)

// get string
str := conf.Get(key).String()

// get string with default
str := conf.Get(key).String(defaultStr)

// get time
t := conf.Get(key).Time()

// get time with default
t := conf.Get(key).Time(defaultTime)

// get nested key values
conf.Get(key).Get(nextKey1).Get(nextKey2).String()
```

### Watch
```golang
conf.Watch(path, func(r *Result){
    r.Scan(x)
})

```

### Stop Watch
```golang
// stop single watcher
conf.StopWatch(path)

// stop multiple watcher
conf.StopWatch(path1, path2)

// stop all watcher
conf.StopWatch()
```

## Next Features:
- Distributed lock
- Transactions
