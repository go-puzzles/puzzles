// File:		dice.go
// Created by:	Hoven
// Created on:	2025-03-31
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package dice

import (
	"math/rand"
)

type Dice struct {
	total    int
	weights  []int
	original []int
}

func NewDice(weights []int) *Dice {
	total := 0

	for _, w := range weights {
		total += w
	}
	return &Dice{
		total:    total,
		weights:  append([]int(nil), weights...),
		original: append([]int(nil), weights...),
	}
}

// Next 根据权重进行随机抽取，支持不放回抽样模式
//
// 工作原理：
// 1. 在权重总和范围内生成随机数
// 2. 按权重区间依次尝试命中
// 3. 命中后将对应位置权重置0（实现不放回）
//
// 使用场景：
// 1. 不放回抽样（默认模式）：
//
//   - 适用于抽奖、抽卡等不希望重复的场景
//
//   - 每个元素在一轮中只会被抽取一次
//
//   - 当所有元素都被抽取后返回-1
//     示例：抽取不重复的奖品
//     for {
//     n := dice.Next()
//     if n == -1 {
//     break  // 所有奖品已抽完
//     }
//     // 处理第n号奖品
//     }
//     此时每一轮第一次抽取的概率是权重最大的元素，第二次抽取的概率是权重第二大的元素，以此类推。
//     并且每一轮元素出现的次数是一样的，
//     在大量抽取时，每个元素出现次数基本相同。
//
// 2. 放回抽样：
//   - 在每次抽取后调用Reset()
//   - 适用于需要保持权重比例的重复抽样
//     示例：模拟多次投掷骰子
//     n := dice.Next()
//     dice.Reset()  // 立即重置，保持权重比例
//     此时在大量抽取时，权重大的元素更容易被抽取，权重小的元素更难被抽取。
//
// 3. 分组抽样：
//   - 抽完一轮后调用Reset()
//   - 适用于需要分组或周期性的抽样
//     示例：每4个一组的抽取
//     for {
//     // 一轮抽取
//     for i := 0; i < 4; i++ {
//     if n := dice.Next(); n != -1 {
//     // 处理第n号元素
//     }
//     }
//     dice.Reset()  // 重置开始新一轮
//     }
//
// 返回值：
//   - 返回被抽中元素的索引（从0开始）
//   - 当所有元素都被抽取后返回-1
func (d *Dice) Next() int {
	if d.total == 0 {
		return -1
	}

	v := rand.Intn(d.total)
	for i, w := range d.weights {
		if v < w {
			d.total -= w
			d.weights[i] = 0
			return i
		}
		v -= w
	}
	return -1
}

func (d *Dice) Reset() {
	d.total = 0
	copy(d.weights, d.original)
	for _, w := range d.weights {
		d.total += w
	}
}
