# dice

dice 是一个基于权重的随机抽取工具，支持不放回抽样模式，适用于抽奖、抽卡、随机选择等场景。

## 功能特性

- 基于权重的随机抽取
- 支持不放回抽样模式
- 支持重置和重复使用
- 高性能的随机算法
- 简单易用的API

## 基本使用

```go
package main

import (
    "fmt"
    "github.com/go-puzzles/puzzles/dice"
)

func main() {
    // 创建一个权重骰子，权重分别为10, 20, 30, 40
    weights := []int{10, 20, 30, 40}
    d := dice.NewDice(weights)
    
    // 不放回抽样 - 每个元素只会被抽取一次
    fmt.Println("不放回抽样:")
    for {
        n := d.Next()
        if n == -1 {
            break // 所有元素已抽完
        }
        fmt.Printf("抽中第%d号元素\n", n)
    }
    
    // 重置骰子，可以重新开始抽取
    d.Reset()
    
    // 放回抽样 - 每次抽取后立即重置
    fmt.Println("\n放回抽样:")
    for i := 0; i < 5; i++ {
        n := d.Next()
        d.Reset() // 立即重置，保持权重比例
        fmt.Printf("第%d次抽取，抽中第%d号元素\n", i+1, n)
    }
}
```

## 使用场景

### 1. 抽奖系统
```go
// 奖品权重：一等奖1，二等奖5，三等奖10，谢谢参与100
prizes := []int{1, 5, 10, 100}
d := dice.NewDice(prizes)

result := d.Next()
switch result {
case 0:
    fmt.Println("恭喜获得一等奖!")
case 1:
    fmt.Println("恭喜获得二等奖!")
case 2:
    fmt.Println("恭喜获得三等奖!")
case 3:
    fmt.Println("谢谢参与!")
}
```

### 2. 负载均衡
```go
// 服务器权重配置
serverWeights := []int{10, 20, 15, 25} // 4台服务器的权重
d := dice.NewDice(serverWeights)

// 选择服务器
serverIndex := d.Next()
d.Reset() // 重置以保持权重比例

fmt.Printf("选择服务器%d\n", serverIndex)
```

### 3. 分组抽样
```go
// 每4个为一组的抽取
for round := 1; round <= 3; round++ {
    fmt.Printf("第%d轮抽取:\n", round)
    for i := 0; i < 4 && d.Next() != -1; i++ {
        // 处理抽取结果
    }
    d.Reset() // 重置开始新一轮
}
```

## API 文档

### NewDice(weights []int) *Dice
创建一个新的权重骰子。

参数：
- `weights`: 权重数组，每个元素的权重值

返回：
- `*Dice`: 骰子实例

### (*Dice) Next() int
进行一次随机抽取。

返回：
- `int`: 被抽中元素的索引（从0开始），如果所有元素都被抽完则返回-1

### (*Dice) Reset()
重置骰子，恢复所有元素的权重，可以重新开始抽取。

## 注意事项

- 权重值必须为正整数
- 不放回模式下，每个元素在一轮中只会被抽取一次
- 使用Reset()可以重新开始新一轮抽取
- 在大量抽取时，权重大的元素被抽取的概率更高 