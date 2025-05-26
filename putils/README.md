# putils

一个功能丰富的 Go 工具库，提供常用的数据处理、验证、网络、时间处理等实用功能。

## 安装

```bash
go get github.com/go-puzzles/puzzles/putils
```

## 功能模块

### 📝 字符串处理 (string.go)

提供高效的字符串搜索算法和处理功能：

- **多种字符串搜索算法**：KMP、Boyer-Moore、暴力搜索
- **智能算法选择**：根据文本长度自动选择最适合的搜索算法

```go
// 字符串搜索
index := putils.StringSearch("hello world", "world")

// 使用指定算法
index := putils.StringSearch("text", "pattern", putils.WithKMP())
index := putils.StringSearch("text", "pattern", putils.WithBM())
```

### 🎲 随机数生成 (random.go)

强大的随机数生成工具：

```go
// 随机布尔值
randBool := putils.RandBool()

// 随机整数 [min, max)
randInt := putils.RandInt(1, 100)

// 随机字符串
randStr := putils.RandString(10)                    // 字母
randNum := putils.RandNumeral(6)                   // 数字
randMixed := putils.RandNumeralOrLetter(8)         // 数字+字母

// 随机浮点数
randFloat := putils.RandFloat(1.0, 10.0, 2)       // 精度为2

// 随机字节
randBytes := putils.RandBytes(16)

// 从切片中随机选择
item := putils.RandFromGivenSlice([]string{"a", "b", "c"})

// 生成唯一随机整数切片
uniqueInts := putils.RandUniqueIntSlice(5, 1, 100)
```

### ✅ 数据验证 (validator.go)

全面的数据验证工具：

```go
// 字符串验证
putils.IsAlpha("abc")                  // 纯字母
putils.IsASCII("hello")                // ASCII字符
putils.IsNumberStr("123")              // 数字字符串
putils.IsFloatStr("123.45")            // 浮点数字符串
putils.IsIntStr("-123")                // 整数字符串

// 网络验证
putils.IsIp("192.168.1.1")            // IP地址
putils.IsIpV4("192.168.1.1")          // IPv4
putils.IsIpV6("::1")                  // IPv6
putils.IsUrl("https://example.com")    // URL
putils.IsEmail("test@example.com")     // 邮箱
putils.IsDns("example.com")           // DNS

// 内容验证
putils.IsJSON(`{"key": "value"}`)      // JSON格式
putils.IsBase64("SGVsbG8=")           // Base64编码
putils.ContainChinese("你好世界")       // 包含中文
putils.ContainLetter("abc123")         // 包含字母
putils.ContainNumber("abc123")         // 包含数字

// 密码强度
putils.IsStrongPassword("Abc123!@#", 8)  // 强密码验证
```

### 📅 日期时间处理 (datetime.go)

便捷的时间处理工具：

```go
now := time.Now()

// 时间边界
beginDay := putils.BeginOfDay(now)     // 当日开始
endDay := putils.EndOfDay(now)         // 当日结束
beginWeek := putils.BeginOfWeek(now)   // 本周开始
endWeek := putils.EndOfWeek(now)       // 本周结束
beginMonth := putils.BeginOfMonth(now) // 本月开始
endMonth := putils.EndOfMonth(now)     // 本月结束

// 快速获取今日、本周、本月时间范围
start, end := putils.StartEndDay()     // 今日
start, end = putils.StartEndWeek()     // 本周
start, end = putils.StartEndMonth()    // 本月

// 时间判断
isLeap := putils.IsLeapYear(2024)      // 闰年判断
isWeekend := putils.IsWeekend(now)     // 是否周末
dayOfYear := putils.DayOfYear(now)     // 一年中的第几天

// 时间计算
seconds := putils.BetweenSeconds(time1, time2)  // 时间差(秒)
```

### 🌐 网络工具 (network.go)

网络信息获取：

```go
// 获取网卡信息
ip, err := putils.GetLocalIP("eth0")        // 本地IP
mac, err := putils.GetMacAddr("eth0")       // MAC地址
mask, err := putils.GetSubnetMask("eth0")   // 子网掩码
```

### 🔧 函数式编程 (function.go)

提供常用的函数式编程工具：

```go
numbers := []int{1, 2, 3, 4, 5}

// Map 操作
doubled := putils.Map(numbers, func(x int) int { return x * 2 })

// Filter 过滤
evens := putils.Filter(numbers, func(x int) bool { return x%2 == 0 })

// Reduce 聚合
sum := putils.Reduce(numbers, 0, func(acc, x int) int { return acc + x })

// 查找
value, found := putils.Find(numbers, func(x int) bool { return x > 3 })

// 分组
grouped := putils.GroupBy(numbers, func(x int) string {
    if x%2 == 0 { return "even" } else { return "odd" }
})

// 检查
hasEven := putils.Any(numbers, func(x int) bool { return x%2 == 0 })
allPositive := putils.All(numbers, func(x int) bool { return x > 0 })

// 包含检查
contains := putils.Contains(numbers, 3)

// 分区
evens, odds := putils.Partition(numbers, func(x int) bool { return x%2 == 0 })
```

### 🗂️ 数据去重 (dedup.go)

智能去重功能，根据数据量选择最优算法：

```go
data := []int{1, 2, 2, 3, 3, 4}
unique := putils.Dedup(data)  // [1, 2, 3, 4]

// 自动选择算法：
// - 小于等于50个元素：使用O(n²)简单算法
// - 大于50个元素：使用O(n)哈希表算法
```

### 📦 字节处理 (bytes.go)

字节和类型转换工具：

```go
// MD5 计算
hash := putils.Md5("hello world")      // 完整MD5
shortHash := putils.ShortMd5("hello")  // 短MD5(16字符)

// 类型安全的字节追加
buf := make([]byte, 0)
buf = putils.AppendAny(buf, 123)       // 追加整数
buf = putils.AppendAny(buf, "hello")   // 追加字符串
buf = putils.AppendAny(buf, time.Now()) // 追加时间
```

### 📚 栈数据结构 (stack.go)

泛型栈实现：

```go
stack := &putils.Stack[int]{}

stack.Push(1)
stack.Push(2)
stack.Push(3)

value, ok := stack.Pop()  // 3, true
```

### 🛠️ 工具函数 (utils.go)

基础工具函数：

```go
// 文件操作
exists := putils.FileExists("/path/to/file")
size, err := putils.FileSize("/path/to/file")

// 随机字符串 (已废弃，建议使用 RandString)
str := putils.GenerateRandomString(10)
```

### 🔄 类型转换 (convert.go)

提供类型转换相关的工具函数。

### ⏰ 时间工具 (time.go)

额外的时间处理工具。

### 🔄 迭代器函数 (iter_function.go)

提供迭代器相关的函数式编程工具。

## 测试

```bash
go test ./...
```

部分模块包含测试文件：
- `random_test.go` - 随机数生成测试
- `string_test.go` - 字符串处理测试

## 特性

- ✨ **零依赖**：除了 `golang.org/x/exp` 外无外部依赖
- 🚀 **高性能**：针对不同场景选择最优算法
- 🔒 **类型安全**：广泛使用 Go 泛型
- 📖 **简单易用**：API 设计简洁直观
- 🧪 **测试覆盖**：核心功能都有测试用例

## 许可证

(c) 2024 Example Corp. All rights reserved.